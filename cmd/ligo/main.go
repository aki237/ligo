package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"plugin"

	"github.com/aki237/ligo/pkg/ligo"
)

var packages = make([]string, 0)

func run() {
	vm := ligo.NewVM()
	vm.Funcs["require"] = VMRequire
	vm.Funcs["load-plugin"] = VMDlLoad
	if len(os.Args) < 2 {
		err := vm.LoadFile(os.Stdin)
		if err != nil {
			fmt.Println(err)
		}
	}
	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	err = vm.LoadFile(f)
	if err != nil {
		fmt.Println(err)
	}
}

func VMRequire(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 1 {
		panic("require : wrong number of arguments")
	}
	lib := a[0]
	if lib.Type != ligo.TYPE_String {
		panic("require : expected a string, got " + lib.GetTypeString())
	}

	packageName := lib.Value.(string)
	err := LoadPackage(vm, packageName)
	if err != nil {
		panic(err)
	}
	return ligo.Variable{Type: ligo.TYPE_Nil, Value: nil}
}

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

	return ligo.Variable{Type: ligo.TYPE_Nil, Value: nil}
}

func exists(dir string) bool {
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		return true
	}
	return false
}

func LoadPackage(vm *ligo.VM, packageName string) error {
	if slistContains(packages, packageName) {
		return nil
	}
	home := os.Getenv("HOME")
	dir := filepath.Join(home, "lispace", "lib", packageName)
	if !exists(dir) {
		return ligo.LigoError("Package \"" + packageName + "\" not found in the system")
	}
	if info, _ := os.Stat(dir); !info.IsDir() {
		return ligo.LigoError("Package \"" + packageName + "\" is not a valid directory")
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
		err = vm.LoadFile(file)
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
	run()
}
