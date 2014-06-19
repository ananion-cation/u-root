// Copyright 2012 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
sh reads in a line at a time and runs it. 
prompt is '% '
*/

package main

import (
	"os/exec"
	"fmt"
	"os"
	"strings"
	"bufio"
)

func main() {
	if len(os.Args) != 1 {
		fmt.Println("no scripts/args yet")
		os.Exit(1)
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("%% ")
	for scanner.Scan() {
		cmd := scanner.Text()
		argv := strings.Split(cmd, " ")
		run := exec.Command(argv[0], argv[1:]...)
		out, err := run.CombinedOutput()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("%s", out)
		}
		fmt.Printf("%% ")
	}
}