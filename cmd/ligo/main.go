package main

import (
	"os"

	"github.com/aki237/ligo/pkg/ligo"
)

func main() {
	vm := ligo.NewVM()
	vm.Funcs["require"] = VMRequire
	vm.Funcs["load-plugin"] = VMDlLoad
	if len(os.Args) < 2 {
		vm.Funcs["exit"] = vmExit
		runInteractive(vm)
		return
	}
	if len(os.Args) == 2 && os.Args[1] == "--web" {
		runWeb()
		return
	}
	vm.Funcs["exit"] = vmExit
	runFile(vm)
}
