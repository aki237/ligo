package main

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"strings"

	"github.com/aki237/ligo/pkg/ligo"
	"github.com/chzyer/readline"
	"github.com/fatih/color"
)

func runInteractive(vm *ligo.VM) {
	expression := ""

	loadRCFile(vm)

	rl, err := readline.New(getPrompt(vm))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rl.Close()

	running := false

	go handleSignals(vm, &running)

	errorFmt := color.New(color.FgRed).Add(color.Bold).Add(color.BgWhite)

	for {
		part, err := rl.Readline()
		switch err {
		case io.EOF:
			fmt.Println("\rBye...")
			os.Exit(0)
		case readline.ErrInterrupt:
			continue
		}

		part = ligo.StripComments(part)
		part = strings.TrimSpace(part)

		if part == "" {
			continue
		}
		if expression == "" && part[0] != '(' {
			v, err := vm.Eval(part)
			if err != nil {
				fmt.Printf("Error in the expression passed : %s\n\t %s\n", errorFmt.Sprintf("%s", err), expression)
				rl.SetPrompt(getPrompt(vm))
				continue
			}
			printValue(v)
			continue
		}
		if expression != "" {
			expression += "\n"
		}
		expression += part
		if ligo.MatchChars(strings.TrimSpace(expression), 0, '(', ')') > 0 {
			rl.SetPrompt(getPrompt(vm))
			running = true
			v, err := vm.Eval(expression)
			if err == ligo.ErrSignalRecieved {
				fmt.Printf("Caught Signal : %s\n", errorFmt.Sprintf("%s", err))
				expression = ""
				vm.Resume()
				running = false
				continue
			}
			if err != nil {
				fmt.Printf("Error in the expression passed : %s\n\t %s\n", errorFmt.Sprintf("%s", err), expression)
				expression = ""
				running = false
				continue
			}
			printValue(v)
			expression = ""
			running = false
			continue
		}
		rl.SetPrompt("... ")
	}
}

func loadRCFile(vm *ligo.VM) {
	home := filepath.Join(os.Getenv("HOME"), ".ligorc")
	f, err := os.Open(home)
	if err == nil {
		vm.LoadReader(f)
		f.Close()
	}
}

func getPrompt(vm *ligo.VM) string {
	defaultPrompt := ">>> "
	ps1, ok := vm.Vars["PS1"]
	if !ok {
		vm.Vars["PS1"] = ligo.Variable{Type: ligo.TypeString, Value: defaultPrompt}
		return defaultPrompt
	}
	psraw, ok := ps1.Value.(string)
	if !ok {
		return defaultPrompt
	}
	return psraw
}

func printValue(v ligo.Variable) {
	if v.Type != ligo.TypeNil {
		fmt.Println("Eval :", v.Value)
	}
}

func handleSignals(vm *ligo.VM, running *bool) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	for sig := range c {
		if *running {
			fmt.Fprintln(os.Stderr, sig)
			vm.Stop()
		}
	}
}
