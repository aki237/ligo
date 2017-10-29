package main

import (
	"io/ioutil"
	"net/http"

	"github.com/aki237/ligo/pkg/ligo"
)

// PluginInit function is the plugin initializer for the url package
func PluginInit(vm *ligo.VM) {
	vm.Funcs["url-get"] = vmURLGet
}

func vmURLGet(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 1 {
		return vm.Throw("urlget expects atleast 1 variable")
	}

	url := a[0].Value.(string)

	resp, err := http.Get(url)

	if err != nil {
		return ligo.Variable{Type: ligo.TypeErr, Value: err}
	}
	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ligo.Variable{Type: ligo.TypeErr, Value: err}
	}

	return ligo.Variable{Type: ligo.TypeString, Value: string(bs)}
}

func main() {

}
