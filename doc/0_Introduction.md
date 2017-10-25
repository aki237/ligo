# Introduction - developing add-on packages in Go

The add-on packages which are written in Go are compiled to a dynamically loadable
library. In this tutorial lets try writing a small add-on package that contains a
function to greet any name (string type).

## ``(require "package")`` - What happens in the interpreter

When a `(require "somePackage")` is called in ligo, the interpreter searches
for this directory in the ligo package search package (ie., `$HOME/lib/somePackage`).
Say your package contains a go compiled plugin, the compiled file should be stored in
that directory with some name and `.plg` extension. If your package contains simple
ligo source files, the files should be placed in the same directory with `.lg`
extension.

When the interpreter loads files from the package directory, if it encounters a `.plg`
file, it `dlopen`s the file and loads the symbol `PluginInit` (which is a function) and
runs the function. If a `.lg` file is encountered, it simply `Evals` the file through
the interpreter.

## Let's write a `.plg` plugin

The go source of the compiled plugin, should always be a main package and
contain a entry point ie `PluginInit` function with `*ligo.VM` as a parameter.

```go
package main

import (
    "fmt"

    "github.com/aki237/ligo/pkg/ligo"
)

func PluginInit(vm *ligo.VM) {
    // Register all the functions for this package.
}
```
This function is where all the functions of this package are registered into the `*ligo.VM`
Now for the example let's write a function of the type `ligo.InBuilt`
(`func(vm *ligo.VM, a ...ligo.Variable) ligo.Variable`).

```go
func greet(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
    // check the number of parameters passed to the function call
    if len(a) != 1 {
        // either panic or just return a nilValue
        return ligo.Variable{Type : ligo.TypeNil, Value : nil}
    }

    // We need only one string argument. So check the type.
    if a[0].Type != ligo.TypeString {
        // either panic or do something...
    }

    // conversion
    name := a[0].Value.(string)

    // Print the greeting.
    fmt.Printf("Hello, %s!!\n", name)

    // return nil or something else...
    return ligo.Variable{Type : ligo.TypeNil, Value : nil}
}
```

Now that this function is written this has to registered in the VM.
There is no namespace. Every function is in the global namespace. So
it is better to name the function as `packageName-functionName`. In this
case the function name is `mypkg-greet` (assuming you name the package `mypkg`).
To register this, in the `PluginInit` function :

```go
func PluginInit(vm *ligo.VM) {
    // Register all the functions for this package.
    vm.Funcs["mypkg-greet"] = greet
}
```

So your package's source will look like this.

```go
package main

import (
    "fmt"

    "github.com/aki237/ligo/pkg/ligo"
)

func PluginInit(vm *ligo.VM) {
    // Register all the functions for this package.
    vm.Funcs["mypkg-greet"] = greet // <- Registering the function
}

func greet(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
    // check the number of parameters passed to the function call
    if len(a) != 1 {
        // either panic or just return a nilValue
        return ligo.Variable{Type : ligo.TypeNil, Value : nil}
    }

    // We need only one string argument. So check the type.
    if a[0].Type != ligo.TypeString {
        // either panic or do something...
    }

    // conversion
    name := a[0].Value.(string)

    // Print the greeting.
    fmt.Printf("Hello, %s!!\n", name)

    // return nil or something else...
    return ligo.Variable{Type : ligo.TypeNil, Value : nil}
}
```

## Building the `.plg`

Make a new directory in `$HOME/ligo/lib/` named mypkg.

```shell
mkdir $HOME/ligo/lib/mypkg
```

To build the plugin

```shell
go build -buildmode=plugin -o $HOME/ligo/lib/mypkg/mypkg.plg file.go
```

## Test it in ligo.

In a new ligo file, test it by calling the function.

```scheme
(require "mypkg")

(mypkg-greet "Lucas") ;; => Returns nothing, Prints "Hello, Lucas!!"
(mypkg-greet 1 2 3 4) ;; => May fail or panic based on your code.
```
