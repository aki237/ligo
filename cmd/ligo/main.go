package main

import (
	"os"

	"github.com/aki237/ligo/pkg/ligo"
)

func main() {
	vm := ligo.NewVM()
	vm.Funcs["require"] = VMRequire
	vm.Funcs["load-plugin"] = VMDlLoad
	vm.Funcs["exit"] = vmExit
	if len(os.Args) < 2 {
		runInteractive(vm)
		return
	}
	runFile(vm)
}
