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
		panic(err)
	}
	defer rl.Close()

	running := false
	new := true

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			if running {
				fmt.Fprintln(os.Stderr, sig)
				vm.Stop()
			}
		}
	}()

	errorFmt := color.New(color.FgRed).Add(color.Bold).Add(color.BgWhite)

	for {
		if new {
			rl.SetPrompt(getPrompt(vm))
		} else {
			rl.SetPrompt("... ")
		}

		part, err := rl.Readline()
		if err == io.EOF {
			fmt.Println("\rBye...")
			break
		}

		part = ligo.StripComments(part)
		part = strings.TrimSpace(part)

		if part == "" {
			continue
		}
		if new {
			if part[0] != '(' {
				fmt.Printf("Error in the expression passed : %s \n\t %s\n",
					errorFmt.Sprintf("%s", "the expression should start with a '(' got '"+string(part[0])+"'"), part)
				expression = ""
				continue
			}
		}
		if expression != "" {
			expression += "\n"
		}
		expression += part
		if ligo.MatchChars(strings.TrimSpace(expression), 0, '(', ')') > 0 {
			new = true
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
			if v.Type != ligo.TypeNil {
				fmt.Println("Eval :", v.Value)
			}
			expression = ""
			running = false
			continue
		}
		new = false
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
		return defaultPrompt
	}
	psraw, ok := ps1.Value.(string)
	if !ok {
		return defaultPrompt
	}
	return psraw
}
