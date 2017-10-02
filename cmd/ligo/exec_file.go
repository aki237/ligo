package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/aki237/ligo/pkg/ligo"
)

func runFile(vm *ligo.VM) {
	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			vm.Stop()
		}
	}()
	err = vm.LoadReader(f)
	if err != nil {
		fmt.Println(err)
	}
}
