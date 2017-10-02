package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"plugin"
	"strings"

	"github.com/aki237/ligo/pkg/ligo"
	"github.com/chzyer/readline"
	"github.com/fatih/color"
)

var packages = make([]string, 0)

func runInteractive() {
	vm := ligo.NewVM()
	vm.Funcs["require"] = VMRequire
	vm.Funcs["load-plugin"] = VMDlLoad
	expression := ""

	rl, err := readline.New(">>> ")
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
			rl.SetPrompt(">>> ")
		} else {
			rl.SetPrompt("... ")
		}

		part, err := rl.Readline()
		if err == io.EOF {
			fmt.Println("\rBye...")
			break
		}

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

func runFile() {
	vm := ligo.NewVM()
	vm.Funcs["require"] = VMRequire
	vm.Funcs["load-plugin"] = VMDlLoad
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

// VMRequire function is a ligo.InBuilt that is used to load a package.
// The package system is still not finalized
func VMRequire(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 1 {
		panic("require : wrong number of arguments")
	}
	lib := a[0]
	if lib.Type != ligo.TypeString {
		panic("require : expected a string, got " + lib.GetTypeString())
	}

	packageName := lib.Value.(string)
	err := LoadPackage(vm, packageName)
	if err != nil {
		panic(err)
	}
	return ligo.Variable{Type: ligo.TypeNil, Value: nil}
}

// VMDlLoad function is a ligo.InBuilt function that is used to load a package that is a dynamically loadable
func VMDlLoad(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 1 {
		panic("load-plugin can only take one argument")
	}

	libpath := a[0].Value.(string)

	p, err := plugin.Open(libpath)
	if err != nil {
		panic(err)
	}
	init, err := p.Lookup("PluginInit")
	if err != nil {
		panic(err)
	}

	init.(func(*ligo.VM))(vm)

	return ligo.Variable{Type: ligo.TypeNil, Value: nil}
}

func exists(dir string) bool {
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		return true
	}
	return false
}

// LoadPackage is used to load a package of the name "packageName" and load the functions and others in the passed ligo.VM.
func LoadPackage(vm *ligo.VM, packageName string) error {
	if slistContains(packages, packageName) {
		return nil
	}
	home := os.Getenv("HOME")
	dir := filepath.Join(home, "lispace", "lib", packageName)
	if !exists(dir) {
		return ligo.Error("Package \"" + packageName + "\" not found in the system")
	}
	if info, _ := os.Stat(dir); !info.IsDir() {
		return ligo.Error("Package \"" + packageName + "\" is not a valid directory")
	}

	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		panic("require : " + fmt.Sprint(err))
	}
	for _, val := range fileInfos {
		if val.IsDir() {
			continue
		}
		if filepath.Ext(val.Name()) != ".lg" {
			if filepath.Ext(val.Name()) == ".plg" {
				p, err := plugin.Open(filepath.Join(dir, val.Name()))
				if err != nil {
					panic(err)
				}
				init, err := p.Lookup("PluginInit")
				if err != nil {
					panic(err)
				}

				init.(func(*ligo.VM))(vm)
			}
			continue
		}
		file, err := os.Open(filepath.Join(dir, val.Name()))
		if err != nil {
			return err
		}
		err = vm.LoadReader(file)
		if err != nil {
			panic(err)
		}
	}
	packages = append(packages, packageName)
	return nil
}

func slistContains(sl []string, s string) bool {
	for _, val := range sl {
		if val == s {
			return true
		}
	}
	return false
}

func main() {
	if len(os.Args) < 2 {
		runInteractive()
		return
	}
	runFile()
}
