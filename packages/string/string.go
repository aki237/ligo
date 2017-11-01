package main

import (
	"fmt"
	"strings"

	"github.com/aki237/ligo/pkg/ligo"
)

// PluginInit function is the plugin initializer for the string package
func PluginInit(vm *ligo.VM) {
	vm.Funcs["indexOf"] = vmStringIndexOf
	vm.Funcs["replace"] = vmStringReplace
	vm.Funcs["split"] = vmStringSplit
	vm.Funcs["splitAfter"] = vmStringSplitAfter
	vm.Funcs["splitN"] = vmStringSplitN
	vm.Funcs["splitAfterN"] = vmStringSplitAfterN
	vm.Funcs["trimSpace"] = vmStringTrimSpace
	vm.Funcs["lowerCase"] = vmStringLowerCase
	vm.Funcs["upperCase"] = vmStringUpperCase
	vm.Funcs["fromArray"] = vmStringFromArray
	vm.Funcs["hasPrefix"] = vmStringHasPrefix
	vm.Funcs["hasSuffix"] = vmStringHasSuffix
	vm.Funcs["compare"] = vmStringCompare
	vm.Funcs["repeat"] = vmStringRepeat
	vm.Funcs["count"] = vmStringCount
	vm.Funcs["contains"] = vmStringContains
	vm.Funcs["containsAny"] = vmStringContainsAny
	vm.Funcs["lastIndex"] = vmStringLastIndex
	vm.Funcs["lastIndexAny"] = vmStringLastIndexAny
	vm.Funcs["trim"] = vmStringTrim
	vm.Funcs["trimPrefix"] = vmStringTrimPrefix
	vm.Funcs["trimSuffix"] = vmStringTrimSuffix
	vm.Funcs["trimLeft"] = vmStringTrimLeft
	vm.Funcs["trimRight"] = vmStringTrimRight
	vm.Funcs["join"] = vmStringJoin
}

func vmStringFromArray(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	ret := ""
	if len(a) != 1 {
		return vm.Throw(fmt.Sprintf("string-fromArray : can take only 1 argument, got %d.", len(a)))
	}

	if a[0].Type != ligo.TypeArray {
		return vm.Throw(fmt.Sprintf("string-fromArray : can take only 1 argument of array type, got %s.", a[0].GetTypeString()))
	}

	arr := a[0].Value.([]ligo.Variable)

	for _, val := range arr {
		switch val.Type {
		case ligo.TypeInt:
			if val.Value.(int64) <= 0 {
				return vm.Throw(fmt.Sprintf("string-fromArray : the array can only contain positive integers, got %d", val.Value.(int64)))
			}
			ret = string(val.Value.(int64))
		case ligo.TypeString:
			ret += val.Value.(string)
		default:
			return vm.Throw(fmt.Sprintf("string-fromArray : the array can only contain positive integers or strings, got %s", val.GetTypeString()))
		}
	}

	return ligo.Variable{Type: ligo.TypeString, Value: ret}
}

func vmStringLowerCase(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 1 {
		return vm.Throw(fmt.Sprintf("string-lowerCase : can take only 1 argument, got %d.", len(a)))
	}

	if a[0].Type != ligo.TypeString {
		return vm.Throw(fmt.Sprintf("string-lowerCase : can take only 1 argument of string type, got %s.", a[0].GetTypeString()))
	}

	ret := strings.ToLower(a[0].Value.(string))

	return ligo.Variable{Type: ligo.TypeString, Value: ret}
}

func vmStringUpperCase(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 1 {
		return vm.Throw(fmt.Sprintf("string-lowerCase : can take only 1 argument, got %d.", len(a)))
	}

	if a[0].Type != ligo.TypeString {
		return vm.Throw(fmt.Sprintf("string-lowerCase : can take only 1 argument of string type, got %s.", a[0].GetTypeString()))
	}

	ret := strings.ToUpper(a[0].Value.(string))

	return ligo.Variable{Type: ligo.TypeString, Value: ret}
}

func vmStringTrimSpace(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 1 {
		return vm.Throw(fmt.Sprintf("string-trimSpace : can take only 1 arguments, got %d.", len(a)))
	}

	if a[0].Type != ligo.TypeString {
		return vm.Throw(fmt.Sprintf("string-trimSpace : can take only 1 argument of string type, got %s.", a[0].GetTypeString()))
	}

	return ligo.Variable{Type: ligo.TypeString, Value: strings.TrimSpace(a[0].Value.(string))}
}

func vmStringIndexOf(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 2 {
		return vm.Throw(fmt.Sprintf("string-indexOf : can take only 2 arguments, got %d.", len(a)))
	}

	if a[0].Type != ligo.TypeString ||
		a[1].Type != ligo.TypeString {
		return vm.Throw(fmt.Sprintf("string-indexOf : can take only 2 arguments of string type, got %s %s.",
			a[0].GetTypeString(),
			a[1].GetTypeString()))
	}

	str1 := a[0].Value.(string)
	str2 := a[1].Value.(string)

	return ligo.Variable{Type: ligo.TypeInt, Value: int64(strings.Index(str1, str2))}
}

func vmStringReplace(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 4 {
		return vm.Throw(fmt.Sprintf("string-replace : should take 4 arguments, got %d.", len(a)))
	}

	if a[0].Type != ligo.TypeString ||
		a[1].Type != ligo.TypeString ||
		a[2].Type != ligo.TypeString ||
		a[3].Type != ligo.TypeInt {
		return vm.Throw(fmt.Sprintf("string-replace : should 4 arguments of (string, string, string, int) types, got (%s, %s, %s, %s).",
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
		return vm.Throw(fmt.Sprintf("string-split : should take 2 arguments, got %d.", len(a)))
	}

	if a[0].Type != ligo.TypeString || a[1].Type != ligo.TypeString {
		return vm.Throw(fmt.Sprintf("string-split : should take 2 arguments of type string, got (%s, %s).", a[0].GetTypeString(), a[1].GetTypeString()))
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

func vmStringHasPrefix(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 2 {
		return vm.Throw(fmt.Sprintf("string-hasPrefix : should take 2 arguments, got %d.", len(a)))
	}

	if a[0].Type != ligo.TypeString || a[1].Type != ligo.TypeString {
		return vm.Throw(fmt.Sprintf("string-hasPrefix : should take 2 arguments of type string, got (%s, %s).", a[0].GetTypeString(), a[1].GetTypeString()))
	}

	str := a[0].Value.(string)
	prefix := a[1].Value.(string)

	return ligo.Variable{Type: ligo.TypeBool, Value: strings.HasPrefix(str, prefix)}
}

func vmStringHasSuffix(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 2 {
		return vm.Throw(fmt.Sprintf("string-hasSuffix : should take 2 arguments, got %d.", len(a)))
	}

	if a[0].Type != ligo.TypeString || a[1].Type != ligo.TypeString {
		return vm.Throw(fmt.Sprintf("string-hasSuffix : should take 2 arguments of type string, got (%s, %s).", a[0].GetTypeString(), a[1].GetTypeString()))
	}

	str := a[0].Value.(string)
	suffix := a[1].Value.(string)

	return ligo.Variable{Type: ligo.TypeBool, Value: strings.HasSuffix(str, suffix)}
}

func vmStringCompare(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 2 {
		return vm.Throw(fmt.Sprintf("string-compare : should take 2 arguments, got %d.", len(a)))
	}

	if a[0].Type != ligo.TypeString || a[1].Type != ligo.TypeString {
		return vm.Throw(fmt.Sprintf("string-compare : should take 2 arguments of type string, got (%s, %s).", a[0].GetTypeString(), a[1].GetTypeString()))
	}

	str1 := a[0].Value.(string)
	str2 := a[1].Value.(string)

	return ligo.Variable{Type: ligo.TypeInt, Value: strings.Compare(str1, str2)}
}

func vmStringRepeat(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 2 {
		return vm.Throw(fmt.Sprintf("string-repeat : should take 2 arguments, got %d.", len(a)))
	}

	if a[0].Type != ligo.TypeString || a[1].Type != ligo.TypeInt {
		return vm.Throw(fmt.Sprintf("string-repeat : should take 2 arguments of type (string, int), got (%s, %s).", a[0].GetTypeString(), a[1].GetTypeString()))
	}

	str := a[0].Value.(string)
	repetitions := a[1].Value.(int64)

	if repetitions < 0 {
		return vm.Throw(fmt.Sprintf("string-repeat : second argument should be a positive integer, got %d.", repetitions))
	}

	return ligo.Variable{Type: ligo.TypeString, Value: strings.Repeat(str, int(repetitions))}
}

func vmStringCount(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 2 {
		return vm.Throw(fmt.Sprintf("string-count : should take 2 arguments, got %d.", len(a)))
	}

	if a[0].Type != ligo.TypeString || a[1].Type != ligo.TypeString {
		return vm.Throw(fmt.Sprintf("string-count : should take 2 arguments of type (string, string), got (%s, %s).",
			a[0].GetTypeString(),
			a[1].GetTypeString(),
		))
	}

	str := a[0].Value.(string)
	substr := a[1].Value.(string)

	return ligo.Variable{Type: ligo.TypeInt, Value: int64(strings.Count(str, substr))}
}

func vmStringContains(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 2 {
		return vm.Throw(fmt.Sprintf("string-contains : should take 2 arguments, got %d.", len(a)))
	}

	if a[0].Type != ligo.TypeString || a[1].Type != ligo.TypeString {
		return vm.Throw(fmt.Sprintf("string-contains : should take 2 arguments of type (string, string), got (%s, %s).",
			a[0].GetTypeString(),
			a[1].GetTypeString(),
		))
	}

	str := a[0].Value.(string)
	substr := a[1].Value.(string)

	return ligo.Variable{Type: ligo.TypeBool, Value: strings.Contains(str, substr)}
}

func vmStringContainsAny(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 2 {
		return vm.Throw(fmt.Sprintf("string-containsAny : should take 2 arguments, got %d.", len(a)))
	}

	if a[0].Type != ligo.TypeString || a[1].Type != ligo.TypeString {
		return vm.Throw(fmt.Sprintf("string-containsAny : should take 2 arguments of type (string, string), got (%s, %s).",
			a[0].GetTypeString(),
			a[1].GetTypeString(),
		))
	}

	str := a[0].Value.(string)
	chars := a[1].Value.(string)

	return ligo.Variable{Type: ligo.TypeBool, Value: strings.ContainsAny(str, chars)}
}

func vmStringLastIndex(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 2 {
		return vm.Throw(fmt.Sprintf("string-lastIndex : should take 2 arguments, got %d.", len(a)))
	}

	if a[0].Type != ligo.TypeString || a[1].Type != ligo.TypeString {
		return vm.Throw(fmt.Sprintf("string-lastIndex : should take 2 arguments of type (string, string), got (%s, %s).",
			a[0].GetTypeString(),
			a[1].GetTypeString(),
		))
	}

	str := a[0].Value.(string)
	substr := a[1].Value.(string)

	return ligo.Variable{Type: ligo.TypeInt, Value: int64(strings.LastIndex(str, substr))}
}

func vmStringLastIndexAny(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 2 {
		return vm.Throw(fmt.Sprintf("string-lastIndexAny : should take 2 arguments, got %d.", len(a)))
	}

	if a[0].Type != ligo.TypeString || a[1].Type != ligo.TypeString {
		return vm.Throw(fmt.Sprintf("string-lastIndexAny : should take 2 arguments of type (string, string), got (%s, %s).",
			a[0].GetTypeString(),
			a[1].GetTypeString(),
		))
	}

	str := a[0].Value.(string)
	substr := a[1].Value.(string)

	return ligo.Variable{Type: ligo.TypeInt, Value: int64(strings.LastIndexAny(str, substr))}
}

func vmStringTrim(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 2 {
		return vm.Throw(fmt.Sprintf("string-trim : should take 2 arguments, got %d.", len(a)))
	}

	if a[0].Type != ligo.TypeString || a[1].Type != ligo.TypeString {
		return vm.Throw(fmt.Sprintf("string-trim : should take 2 arguments of type (string, string), got (%s, %s).",
			a[0].GetTypeString(),
			a[1].GetTypeString(),
		))
	}

	str := a[0].Value.(string)
	cutset := a[1].Value.(string)

	return ligo.Variable{Type: ligo.TypeString, Value: strings.Trim(str, cutset)}
}

func vmStringTrimPrefix(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 2 {
		return vm.Throw(fmt.Sprintf("string-trimPrefix : should take 2 arguments, got %d.", len(a)))
	}

	if a[0].Type != ligo.TypeString || a[1].Type != ligo.TypeString {
		return vm.Throw(fmt.Sprintf("string-trimPrefix : should take 2 arguments of type (string, string), got (%s, %s).",
			a[0].GetTypeString(),
			a[1].GetTypeString(),
		))
	}

	str := a[0].Value.(string)
	prefix := a[1].Value.(string)

	return ligo.Variable{Type: ligo.TypeString, Value: strings.TrimPrefix(str, prefix)}
}

func vmStringTrimSuffix(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 2 {
		return vm.Throw(fmt.Sprintf("string-trimSuffix : should take 2 arguments, got %d.", len(a)))
	}

	if a[0].Type != ligo.TypeString || a[1].Type != ligo.TypeString {
		return vm.Throw(fmt.Sprintf("string-trimSuffix : should take 2 arguments of type (string, string), got (%s, %s).",
			a[0].GetTypeString(),
			a[1].GetTypeString(),
		))
	}

	str := a[0].Value.(string)
	suffix := a[1].Value.(string)

	return ligo.Variable{Type: ligo.TypeString, Value: strings.TrimSuffix(str, suffix)}
}

func vmStringTrimLeft(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 2 {
		return vm.Throw(fmt.Sprintf("string-trimLeft : should take 2 arguments, got %d.", len(a)))
	}

	if a[0].Type != ligo.TypeString || a[1].Type != ligo.TypeString {
		return vm.Throw(fmt.Sprintf("string-trimLeft : should take 2 arguments of type (string, string), got (%s, %s).",
			a[0].GetTypeString(),
			a[1].GetTypeString(),
		))
	}

	str := a[0].Value.(string)
	cutset := a[1].Value.(string)

	return ligo.Variable{Type: ligo.TypeString, Value: strings.TrimLeft(str, cutset)}
}

func vmStringTrimRight(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 2 {
		return vm.Throw(fmt.Sprintf("string-trimRight : should take 2 arguments, got %d.", len(a)))
	}

	if a[0].Type != ligo.TypeString || a[1].Type != ligo.TypeString {
		return vm.Throw(fmt.Sprintf("string-trimRight : should take 2 arguments of type (string, string), got (%s, %s).",
			a[0].GetTypeString(),
			a[1].GetTypeString(),
		))
	}

	str := a[0].Value.(string)
	cutset := a[1].Value.(string)

	return ligo.Variable{Type: ligo.TypeString, Value: strings.TrimRight(str, cutset)}
}

func vmStringSplitAfter(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 2 {
		return vm.Throw(fmt.Sprintf("string-splitAfter : should take 2 arguments, got %d.", len(a)))
	}

	if a[0].Type != ligo.TypeString || a[1].Type != ligo.TypeString {
		return vm.Throw(fmt.Sprintf("string-splitAfter : should take 2 arguments of type string, got (%s, %s).",
			a[0].GetTypeString(),
			a[1].GetTypeString(),
		))
	}

	str := a[0].Value.(string)
	sep := a[1].Value.(string)

	splitted := strings.SplitAfter(str, sep)
	ret := make([]ligo.Variable, 0)
	for _, val := range splitted {
		ret = append(ret, ligo.Variable{Type: ligo.TypeString, Value: val})
	}

	return ligo.Variable{Type: ligo.TypeArray, Value: ret}
}

func vmStringSplitN(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 3 {
		return vm.Throw(fmt.Sprintf("string-splitN : should take 3 arguments, got %d.", len(a)))
	}

	if a[0].Type != ligo.TypeString || a[1].Type != ligo.TypeString || a[2].Type != ligo.TypeInt {
		return vm.Throw(fmt.Sprintf("string-splitN : should take 3 arguments of type (string,string,int), got (%s, %s, %s).",
			a[0].GetTypeString(),
			a[1].GetTypeString(),
			a[2].GetTypeString(),
		))
	}

	str := a[0].Value.(string)
	sep := a[1].Value.(string)
	n := a[2].Value.(int)

	splitted := strings.SplitN(str, sep, n)
	ret := make([]ligo.Variable, 0)
	for _, val := range splitted {
		ret = append(ret, ligo.Variable{Type: ligo.TypeString, Value: val})
	}

	return ligo.Variable{Type: ligo.TypeArray, Value: ret}
}

func vmStringSplitAfterN(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 3 {
		return vm.Throw(fmt.Sprintf("string-splitAfterN : should take 3 arguments, got %d.", len(a)))
	}

	if a[0].Type != ligo.TypeString || a[1].Type != ligo.TypeString || a[2].Type != ligo.TypeInt {
		return vm.Throw(fmt.Sprintf("string-splitAfterN : should take 3 arguments of type (string,string,int), got (%s, %s, %s).",
			a[0].GetTypeString(),
			a[1].GetTypeString(),
			a[2].GetTypeString(),
		))
	}

	str := a[0].Value.(string)
	sep := a[1].Value.(string)
	n := a[2].Value.(int)

	splitted := strings.SplitAfterN(str, sep, n)
	ret := make([]ligo.Variable, 0)
	for _, val := range splitted {
		ret = append(ret, ligo.Variable{Type: ligo.TypeString, Value: val})
	}

	return ligo.Variable{Type: ligo.TypeArray, Value: ret}
}

func vmStringJoin(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 2 {
		return vm.Throw(fmt.Sprintf("string-join : should take 2 arguments, got %d.", len(a)))
	}

	if a[0].Type != ligo.TypeArray || a[1].Type != ligo.TypeString {
		return vm.Throw(fmt.Sprintf("string-join : should take 2 arguments of type (array,string), got (%s, %s).",
			a[0].GetTypeString(),
			a[1].GetTypeString(),
		))
	}

	arrayValues := a[0].Value.([]ligo.Variable)
	if arrayValues[0].Type != ligo.TypeString {
		return vm.Throw(fmt.Sprintf("string-join : 1 argument should be an array of string type, got array of (%s) type.",
			arrayValues[0].GetTypeString(),
		))
	}

	var items []string
	for _, v := range arrayValues {
		items = append(items, v.Value.(string))
	}

	sep := a[1].Value.(string)

	return ligo.Variable{Type: ligo.TypeString, Value: strings.Join(items, sep)}
}

func main() {

}
