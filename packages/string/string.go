package main

import (
	"fmt"
	"strings"

	"github.com/aki237/ligo/pkg/ligo"
)

func PluginInit(vm *ligo.VM) {
	vm.Funcs["string-indexOf"] = vmStringIndexOf
	vm.Funcs["string-replace"] = vmStringReplace
	vm.Funcs["string-split"] = vmStringSplit
	vm.Funcs["string-trimSpace"] = vmStringTrimSpace
	vm.Funcs["string-lowerCase"] = vmStringLowerCase
	vm.Funcs["string-upperCase"] = vmStringUpperCase
	vm.Funcs["string-fromArray"] = vmStringFromArray
	vm.Funcs["string-repeat"] = vmStringRepeat
}

func vmStringFromArray(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	ret := ""
	if len(a) != 1 {
		panic(fmt.Sprintf("string-fromArray : can take only 1 argument, got %d.", len(a)))
	}

	if a[0].Type != ligo.TypeArray {
		panic(fmt.Sprintf("string-fromArray : can take only 1 argument of array type, got %s.", a[0].GetTypeString()))
	}

	arr := a[0].Value.([]ligo.Variable)

	for _, val := range arr {
		switch val.Type {
		case ligo.TypeInt:
			if val.Value.(int64) <= 0 {
				panic(fmt.Sprintf("string-fromArray : the array can only contain positive integers, got %d", val.Value.(int64)))
			}
			ret = string(val.Value.(int64))
		case ligo.TypeString:
			ret += val.Value.(string)
		default:
			panic(fmt.Sprintf("string-fromArray : the array can only contain positive integers or strings, got %s", val.GetTypeString()))
		}
	}

	return ligo.Variable{Type: ligo.TypeString, Value: ret}
}

func vmStringLowerCase(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 1 {
		panic(fmt.Sprintf("string-lowerCase : can take only 1 argument, got %d.", len(a)))
	}

	if a[0].Type != ligo.TypeString {
		panic(fmt.Sprintf("string-lowerCase : can take only 1 argument of string type, got %s.", a[0].GetTypeString()))
	}

	ret := strings.ToLower(a[0].Value.(string))

	return ligo.Variable{Type: ligo.TypeString, Value: ret}
}

func vmStringUpperCase(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 1 {
		panic(fmt.Sprintf("string-lowerCase : can take only 1 argument, got %d.", len(a)))
	}

	if a[0].Type != ligo.TypeString {
		panic(fmt.Sprintf("string-lowerCase : can take only 1 argument of string type, got %s.", a[0].GetTypeString()))
	}

	ret := strings.ToUpper(a[0].Value.(string))

	return ligo.Variable{Type: ligo.TypeString, Value: ret}
}

func vmStringTrimSpace(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 1 {
		panic(fmt.Sprintf("string-trimSpace : can take only 1 arguments, got %d.", len(a)))
	}

	if a[0].Type != ligo.TypeString {
		panic(fmt.Sprintf("string-trimSpace : can take only 1 argument of string type, got %s.", a[0].GetTypeString()))
	}

	return ligo.Variable{Type: ligo.TypeString, Value: strings.TrimSpace(a[0].Value.(string))}
}

func vmStringIndexOf(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 2 {
		panic(fmt.Sprintf("string-indexOf : can take only 2 arguments, got %d.", len(a)))
	}

	if a[0].Type != ligo.TypeString ||
		a[1].Type != ligo.TypeString {
		panic(fmt.Sprintf("string-indexOf : can take only 2 arguments of string type, got %s %s.",
			a[0].GetTypeString(),
			a[1].GetTypeString()))
	}

	str1 := a[0].Value.(string)
	str2 := a[1].Value.(string)

	return ligo.Variable{Type: ligo.TypeInt, Value: int64(strings.Index(str1, str2))}
}

func vmStringReplace(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 4 {
		panic(fmt.Sprintf("string-replace : should take 4 arguments, got %d.", len(a)))
	}

	if a[0].Type != ligo.TypeString ||
		a[1].Type != ligo.TypeString ||
		a[2].Type != ligo.TypeString ||
		a[3].Type != ligo.TypeInt {
		panic(fmt.Sprintf("string-replace : should 4 arguments of (string, string, string, int) types, got (%s, %s, %s, %s).",
			a[0].GetTypeString(),
			a[1].GetTypeString(),
			a[2].GetTypeString(),
			a[3].GetTypeString(),
		))
	}

	str1 := a[0].Value.(string)
	str2 := a[1].Value.(string)
	str3 := a[2].Value.(string)
	times := a[3].Value.(int64)
	return ligo.Variable{Type: ligo.TypeString, Value: strings.Replace(str1, str2, str3, int(times))}
}

func vmStringSplit(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 2 {
		panic(fmt.Sprintf("string-split : should take 2 arguments, got %d.", len(a)))
	}

	if a[0].Type != ligo.TypeString || a[1].Type != ligo.TypeString {
		panic(fmt.Sprintf("string-split : should take 2 arguments of type string, got (%s, %s).", a[0].GetTypeString(), a[1].GetTypeString()))
	}

	str1 := a[0].Value.(string)
	str2 := a[1].Value.(string)

	splitted := strings.Split(str1, str2)
	ret := make([]ligo.Variable, 0)
	for _, val := range splitted {
		ret = append(ret, ligo.Variable{Type: ligo.TypeString, Value: val})
	}

	return ligo.Variable{Type: ligo.TypeArray, Value: ret}
}

func vmStringRepeat(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 2 {
		panic(fmt.Sprintf("string-repeat : should take 2 arguments, got %d.", len(a)))
	}

	if a[0].Type != ligo.TypeString || a[1].Type != ligo.TypeInt {
		panic(fmt.Sprintf("string-repeat : should take 2 arguments of type (string, int), got (%s, %s).", a[0].GetTypeString(), a[1].GetTypeString()))
	}

	str := a[0].Value.(string)
	repetitions := a[1].Value.(int)

	if repetitions < 0 {
		panic(fmt.Sprintf("string-repeat : second argument should be a positive integer, got %d.", repetitions))
	}

	return ligo.Variable{Type: ligo.TypeString, Value: strings.Repeat(str, repetitions)}
}

func main() {

}
