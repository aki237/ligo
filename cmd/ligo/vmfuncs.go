package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"plugin"

	"github.com/aki237/ligo/pkg/ligo"
)

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

var packages = make([]string, 0)

// LoadPackage is used to load a package of the name "packageName" and load the functions and others in the passed ligo.VM.
func LoadPackage(vm *ligo.VM, packageName string) error {
	if slistContains(packages, packageName) {
		return nil
	}
	home := os.Getenv("HOME")

	// Either /home/$USER or $LIGOPATH can be a path for library searching
	if ligopath := os.Getenv("LIGOPATH"); ligopath != "" {
		home = ligopath
	}

	dir := filepath.Join(home, "ligo", "lib", packageName)
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

func vmExit(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	vm.Stop()
	os.Exit(0)
	return ligo.Variable{Type: ligo.TypeNil, Value: nil}
}
