package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/aki237/ligo/pkg/ligo"
)

var versionString = "0.0.1"

func usage() {
	printVersion()
	fmt.Println("Usage : ligo [filenames]")
	flag.PrintDefaults()
}

func printVersion() {
	fmt.Printf("ligo %s %s %s : ligo language interpreter\n", versionString, runtime.GOOS, runtime.GOARCH)
}

func main() {

	flag.Usage = usage
	version := flag.Bool("version", false, "Print the version information")

	flag.Parse()

	if *version {
		printVersion()
		return
	}

	os.Args = flag.Args()

	vm := ligo.NewVM()
	vm.Funcs["require"] = VMRequire
	vm.Funcs["load-plugin"] = VMDlLoad
	vm.Funcs["exit"] = vmExit
	if len(os.Args) < 1 {
		runInteractive(vm)
		return
	}
	runFile(vm)
}
