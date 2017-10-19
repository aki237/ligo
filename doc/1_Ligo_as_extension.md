# Using ligo as a extention languages to packages

The ligo library can be used to extend golang's power using the scripting functionality.
This tutorial says how to use ligo as a extention language easily. If you have read the [introduction](0_Introduction.md)
this won't be necessary because you'll have the idea how to use ligo as a extention utility.

```go
package main

import (
    "fmt"
    
    "github.com/aki237/ligo/pkg/ligo"
)

func main() {
    // Create a new VM Scope
    vm := ligo.NewVM()
    
    // Similar to this load some functions
    // the function pattern is func(a ...ligo.Variable) ligo.Variable
	vm.Funcs["greet"] = Greet
    
    // you can get inputs from anywhere like from net.Conn or from os.File, or even a string
    // if it is a io.Reader you can directly call vm.LoadReader to just run the script. It won't return any value.
    // on the other hand if it is a string you can call `vm.Eval()`. Which returns a value (ligo.Variable) and an error.
    // Eval can only take in a string that has a single lisp expression. In that case you can pass that string through 
    // vm.BreakChunk method to break the ligo code into a array of ligo expressions, each of them which can be passed
    // through vm.Eval.
    
    // Here let's consider a string with a single expression (ie., BreakChunk step is not necessary).
    exp := `(greet "Alice")`
    
    val, err := vm.Eval(exp)
    if err != nil {
        fmt.Println(err)
        return
    }
    
    fmt.Println(val.Value)
}

func Greet(a ...ligo.Variable) ligo.Variable {
    // cheking the arg count passed for the function
    if len(a) != 1 {
        // do some thing
    }
    
    // check the type of the needed first parameter
    if a[0].Type != ligo.TypeString {
        // do something
    }
    
    return ligo.Variable{Type : ligo.TypeString, Value : "Hello, " + a[0].Value.(string) + "!!"}
}
```

This program will output :

```
Hello, Alice!!
```

Add adapters for your custom functions and go nuts!!

**PS** : None of the basic operations (like +,-,/,*,or,and etc.,) are added in the ligo package. You can copy the adapters from the packages/base/base.go directory
to your project.
