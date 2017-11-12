package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/aki237/ligo/pkg/ligo"
)

func runFile(vm *ligo.VM) {
	for _, val := range os.Args {
		f, err := os.Open(val)
		if err != nil {
			fmt.Println(err)
			return
		}
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		go func() {
			for range c {
				vm.Stop()
				os.Exit(0)
			}
		}()
		err = vm.LoadReader(f)
		if err != nil {
			fmt.Println(err)
		}
	}
}
