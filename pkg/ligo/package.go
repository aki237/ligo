/*
Package ligo provides implementations for running lisp like scripts

Usage

Using this package is very simple.
  Create a new VM instance
 - Define some functions in the VM
 - Run any ligo code

Sample

In this sample, a new function called printHello is added to VM
This can be called from the ligo code. Running ligo is as simple as
vm.Eval(code). The package itself contains only the basic parsing and
running functionality. No arithmetic or logical functions are added.
For implementing it yourslef, just look at the docs in the github repo : https://github.com/aki237/ligo

    func printHello(vm *VM, a ...ligo.Variable) ligo.Variable {
        fmt.Println ("Hello from ligo!!")
        return ligo.Variable{ligo.TypeNil, nil}
    }


    func main() {
        vm := ligo.NewVM()
        vm.Funcs["hello"] = printHello
        vm.Eval("(hello)")
    }

The above gives an output :

    Hello from ligo!!

*/
package ligo
