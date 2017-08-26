package main
//include a loading bar
import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"	
	"strconv"
	"io/ioutil"
	"os"
	"net/http"
	"io"
)

//TODO generic path 
var PATH = "/usr/local/google/home/ananyajoshi/"


//TODO error handling! 
func goCompatibility(){
	cmd := exec.Command("go", "version")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	// The string is originally in the form: go version go1.9rc2_cl165246139 linux/amd64 where 1.9 is the go version
 	termString, err := strconv.ParseFloat(strings.Split(out.String(), " ")[2][2:5], 64)
	if err != nil {
		log.Fatal(err)
	}
	if(termString > 1.7){		
		fmt.Println("Compatible go version")
	} else {
		fmt.Println("Please install go v1.7 or greater.")
	}
}

func goGet(){
	cmd1 := exec.Command("go", "get", "github.com/u-root/u-root")
	err := cmd1.Run()
	if err != nil {
		log.Fatal(err)
	}
	cmd2 := exec.Command("export ", "GOPATH=\"$HOME/go\"")
	err = cmd2.Run()
	if err != nil {
		log.Fatal(err)
	}
	//syscalls

	cmd3 := exec.Command("cd", "\"$GOPATH/src/github.com/u-root/u-root\"")
	err = cmd3.Run()
	if err != nil {
		log.Fatal(err)
	}
	cmd4 := exec.Command("go", "run", "scripts/ramfs.go")
	err = cmd4.Run()
	if err != nil {
		log.Fatal(err)
	}
	//check if the tmpFile exists
	if _, err :=os.Stat("/tmp/initramgs.linux_amd64.cpio"); err != nil{
		log.Fatal(err)	
	}
	fmt.Printf("Created the initramfs in /tmp/")
}



func kernelGet() error{
	//check about tmp files
	//jsonFile, err := http.Get("https://www.kernel.org/releases.json")
	jsonFile, err := http.Get("http://www.kernel.org/releases.json")
	if err != nil {
		return err
	}
	defer jsonFile.Body.Close()
	d, err := ioutil.ReadAll(jsonFile.Body)
	version := strings.Split(strings.SplitAfter(strings.Split(strings.Split(string(d), "\"moniker\": \"stable\",")[1], "\n")[1], "https:")[1], "\",")[0]
	fmt.Printf("%v", version)
	/* make for user */ 
	out, err := os.Create("/usr/local/google/home/ananyajoshi/Downloads/linux.tar.xz")
	if err != nil  {
   		 return err
  	}
  	defer out.Close()
	resp, err := http.Get(fmt.Sprintf("http:%s",version))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil  {
		return err
	}
	return nil
}
 
//TODO temp directory
//TODO make sure the kernel and the config file syncs up 
//TODO new config file with the defaults set up 
func unpackKernel() error{
//golang tar utility 
	cmd5 := exec.Command("tar", "xf", fmt.Sprintf(PATH, "Downloads/linux.tar.xz")) 
	 err := cmd5.Run() 
	 if err != nil { 
	 	log.Fatal(err)
	 }
	cmd6 := exec.Command("cd", "linux") 
	 err := cmd6.Run() 
	 if err != nil { 
	 	log.Fatal(err)
	 }
	cmd7 := exec.Command("git", "init") 
	 err := cmd7.Run() 
	 if err != nil { 
	 log.Fatal(err)
	 }
	cmd8 := exec.Command("git", "clone", "https://github.com/u-root/NiChrome.git") 
	 err := cmd8.Run() 
	 if err != nil { 
	 log.Fatal(err)
	 }
//use u-root's cp 
	cmd9 := exec.Command("cp", "NiChrome/CONFIG", ".config") 
	 err := cmd9.Run() 
	 if err != nil { 
	 log.Fatal(err)
	 }
	cmd11 := exec.Command("cp", "/tmp/initramfs.linux_amd64.cpio", ".") 
 	err := cmd11.Run() 
	 if err != nil { 
 		log.Fatal(err)
	 }
	cmd11 := exec.Command("sudo", "make", "-j64") 
	 err := cmd11.Run() 
	 if err != nil { 
	 log.Fatal(err)
	 }
	//TODO check if the kernel exists
	//pass the proper parameters to the functions
}


//MERGE CONFLICT khsdlkgjslk
func vbutilIt(bzImageName string, ){
	cmd12 := exec.Command("mkdir", "vbutilDir") 
	 err := cmd12.Run() 
	 if err != nil { 
	 log.Fatal(err)
	 }
	cmd13 := exec.Command("cd", "vbutilDir") 
	 err := cmd13.Run() 
	 if err != nil { 
	 log.Fatal(err)
	 }
	//TODO: do this the way that rm_tests does it 
	if err := ioutil.WriteFile("config.txt", []byte("loglevel=7"), 0777); err != nil {
			return "", err
	}
	cmd14 := exec.Command("vbutil_kernel", "--pack ", NAME_OF_FILE, "--keyblock", "/usr/share/vboot/devkeys/kernel.keyblock", "--signprivate", "/usr/share/vboot/devkeys/kernel_data_key.vbprivk", "--version",  "1",  "--config",  "config.txt",  "--vmlinuz", bzImageName) 
	 err := cmd14.Run() 
	 if err != nil { 
	 log.Fatal(err)
	 }
	cmd15 := exec.Command("sudo dd if=NAME_OF_FILE of=/dev/sda2") 
	 err := cmd15.Run() 
	 if err != nil { 
	 log.Fatal(err)
	 }
}



func main() {
	goCompatibility()
	//goGet()
	kernelGet()	
}




