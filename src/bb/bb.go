// bb converts standalone u-root tools to shell builtins.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path"
	"reflect"
	"text/template"

	"golang.org/x/tools/imports"
)

const cmdFunc = `func {{.CmdName}}(c *Command) (err error) {
save := os.Args
defer func() {
os.Args = save
        if r := recover(); r != nil {
            err = r.(error)
        }
return
    }()
os.Args = append([]string{c.cmd}, c.argv...)
{{.CmdName}}_main()
return
}

func init() {
	addBuiltIn("{{.CmdName}}", {{.CmdName}})
}
`

var config struct {
	Args []string
	CmdName string
	FullPath string
	Src string
}

func one(s string, fset *token.FileSet, f *ast.File) error {
	// Inspect the AST and change all instances of main()
	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			if false {
				fmt.Printf("%v", reflect.TypeOf(x.Type.Params.List[0].Type))
			}
			if x.Name.Name == "main" {
				x.Name.Name = fmt.Sprintf("%v_main", config.CmdName)
			}
			// Append a return.
			x.Body.List = append(x.Body.List, &ast.ReturnStmt{})

			case *ast.CallExpr:
			fmt.Fprintf(os.Stderr, "%v %v\n", reflect.TypeOf(n), n)
			switch z := x.Fun.(type) {
			case *ast.SelectorExpr:
				// somebody tell me how to do this.
				sel := fmt.Sprintf("%v", z.X)
				if  sel == "os" && z.Sel.Name == "Exit" {
					x.Fun = &ast.Ident{Name: "panic"}
				}
			}
		}
		return true
	})
	
	if false {
		ast.Fprint(os.Stderr, fset, f, nil)
	}
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, f); err != nil {
	panic(err)
	}
	fmt.Printf("%s", buf.Bytes())

	t := template.Must(template.New("cmdFunc").Parse(cmdFunc))
	var b bytes.Buffer
	if err := t.Execute(&b, config); err != nil {
		log.Fatalf("spec %v: %v\n", cmdFunc, err)
	}
	out := string(buf.Bytes()) + b.String()

	// fix up any imports. We may have forced the issue
	// with os.Args
	opts := imports.Options{
		Fragment:  true,
		AllErrors: true,
		Comments:  true,
		TabIndent: true,
		TabWidth:  8,
	}
	fullCode, err := imports.Process("commandline", []byte(out), &opts)
	if err != nil {
		log.Fatalf("bad parse: '%v': %v", out, err)
	}

	of := path.Join("bbsh", path.Base(s))
	if err := ioutil.WriteFile(of, []byte(fullCode), 0666); err != nil {
		log.Fatalf("%v\n", err)
	}
	return nil
}

func main() {
	flag.Parse()
	config.Args = flag.Args()
	config.CmdName = config.Args[0]
	if len(config.Args) != 1 {
		log.Fatalf("usage: bb <one directory>\n")
	}
	fset := token.NewFileSet()
	config.FullPath = path.Join(os.Getenv("UROOT"), "src/cmds", config.CmdName)
	p, err := parser.ParseDir(fset, config.FullPath, nil, 0)
	if err != nil {
		panic(err)
	}

	for _, f := range p {
		for n, v := range f.Files {
			one(n, fset, v)
		}
	}
}