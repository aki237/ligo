package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"time"

	"github.com/aki237/ligo/pkg/ligo"
)

// PluginInit function is the plugin initializer for the base package
func PluginInit(vm *ligo.VM) {
	vm.Funcs["println"] = vmPrintln
	vm.Funcs["vmmem"] = vmMem
	vm.Funcs["input"] = vmInput
	vm.Funcs["input-lines"] = vmInputLines
	vm.Funcs["array-index"] = vmArrayIndex
	vm.Funcs["print"] = vmPrint
	vm.Funcs["=="] = vmEquality
	vm.Funcs["car"] = vmCar
	vm.Funcs["cdr"] = vmCdr
	vm.Funcs["len"] = vmLen
	vm.Funcs[">="] = vmInEqualityGTEQ
	vm.Funcs["<="] = vmInEqualityLTEQ
	vm.Funcs[">"] = vmInEqualityGT
	vm.Funcs["+"] = vmAdd
	vm.Funcs["*"] = vmProd
	vm.Funcs["%"] = vmModulus
	vm.Funcs["type"] = vmType
	vm.Funcs["throw"] = vmThrow
	vm.Funcs["sleep"] = vmSleep
	vm.Funcs["reciprocal"] = vmReciprocal
	vm.Funcs["array-set"] = vmArraySet
	vm.Funcs["array-subArray"] = vmArraySubArray
	vm.Funcs["array-append"] = vmArrayAppend
	vm.Funcs["or"] = vmOr
	vm.Funcs["and"] = vmAnd
	vm.Funcs["not"] = vmNot
	vm.Funcs["is-nil"] = vmIsNil
	vm.Funcs["sprintf"] = vmSprintf
	vm.Funcs["map-new"] = vmMapNew
	vm.Funcs["map-store"] = vmMapStore
	vm.Funcs["map-delete"] = vmMapDelete
	vm.Funcs["map-get"] = vmMapGet
}

func vmMapNew(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	return ligo.Variable{Type: ligo.TypeMap, Value: make(ligo.Map, 0)}
}

func vmMapStore(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 3 {
		vm.Throw("map-store : requires 3 arguments")
	}

	if a[0].Type != ligo.TypeMap {
		vm.Throw("map-store : expected a <Map> type as the first argument")
	}

	a[0].Value.(ligo.Map)[a[1]] = a[2]
	return a[0]
}

func vmMapDelete(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 2 {
		vm.Throw("map-delete : requires 2 arguments")
	}

	if a[0].Type != ligo.TypeMap {
		vm.Throw("map-delete : expected a <Map> type as the first argument")
	}

	delete(a[0].Value.(ligo.Map), a[1])
	return a[0]
}

func vmMapGet(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 2 {
		vm.Throw("map-get : requires 2 arguments")
	}

	if a[0].Type != ligo.TypeMap {
		vm.Throw("map-get : expected a <Map> type as the first argument")
	}

	v, ok := a[0].Value.(ligo.Map)[a[1]]
	if !ok {
		return ligo.Variable{Type: ligo.TypeNil, Value: nil}
	}
	return v
}

func vmMem(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	mem := &runtime.MemStats{}
	runtime.ReadMemStats(mem)
	fmt.Println(float64(mem.Alloc) / (1024.0 * 1024.0))
	return ligo.Variable{Type: ligo.TypeNil, Value: nil}
}

func vmArraySubArray(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 3 {
		vm.Throw(fmt.Sprintf("array-subArray: require 3 arguments, got %d arguments", len(a)))
	}

	if (a[0].Type != ligo.TypeArray && a[0].Type != ligo.TypeString) ||
		a[1].Type != ligo.TypeInt ||
		a[2].Type != ligo.TypeInt {
		vm.Throw(fmt.Sprintf("array-subArray: require 3 arguments (array, int, int), got (%s %s %s) arguments",
			a[0].GetTypeString(),
			a[1].GetTypeString(),
			a[2].GetTypeString(),
		))
	}

	if a[0].Type == ligo.TypeString {
		arr := a[0].Value.(string)
		start := a[1].Value.(int64)
		end := a[2].Value.(int64)

		if start >= int64(len(arr)) || start < 0 || end < start || end > int64(len(arr)) {
			vm.Throw(fmt.Sprintf("array-subArray: invalid array index number %d %d %d", start, end, len(arr)))
		}

		return ligo.Variable{Type: ligo.TypeString, Value: string(arr[start:end])}
	}

	arr := a[0].Value.([]ligo.Variable)
	start := a[1].Value.(int64)
	end := a[2].Value.(int64)

	if start >= int64(len(arr)) || start < 0 || end < start || end > int64(len(arr)) {
		vm.Throw(fmt.Sprintf("array-subArray: invalid array index number %d %d %d", start, end, len(arr)))
	}

	return ligo.Variable{Type: ligo.TypeArray, Value: arr[start:end]}
}

func vmArrayIndex(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 2 {
		vm.Throw(fmt.Sprintf("array-index: require 2 arguments, got %d arguments", len(a)))
	}

	if (a[0].Type != ligo.TypeArray && a[0].Type != ligo.TypeString) ||
		a[1].Type != ligo.TypeInt {
		vm.Throw(fmt.Sprintf("array-index: require 2 arguments (array, int), got (%s %s) arguments", a[0].GetTypeString(), a[1].GetTypeString()))
	}

	if a[0].Type == ligo.TypeString {
		arr := a[0].Value.(string)
		nth := a[1].Value.(int64)

		if nth >= int64(len(arr)) {
			vm.Throw(fmt.Sprintf("array-index: index exceeding array-length : index (%d) > array-length (%d)", nth, len(arr)))
		}

		return ligo.Variable{Type: ligo.TypeString, Value: string(arr[nth])}
	}

	arr := a[0].Value.([]ligo.Variable)
	nth := a[1].Value.(int64)

	if nth >= int64(len(arr)) {
		vm.Throw(fmt.Sprintf("array-index: index exceeding array-length : index (%d) > array-length (%d)", nth, len(arr)))
	}

	return arr[nth]
}

func vmInputLines(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) > 1 {
		vm.Throw(fmt.Sprintf("input: require less than 2 arguments, got %d arguments", len(a)))
	}

	if len(a) == 1 {
		vmPrint(vm, a...)
	}

	lines := make([]ligo.Variable, 0)

	rd := bufio.NewReader(os.Stdin)

	for true {
		input, err := rd.ReadString('\n')

		if err == io.EOF {
			return ligo.Variable{Type: ligo.TypeArray, Value: lines}
		}

		if err != nil {
			vm.Throw(fmt.Sprintf("input : panicked when trying to read from stdin : %s", err))
		}

		if len(input) > 0 && input[len(input)-1] == '\n' {
			input = input[:len(input)-1]
		}
		lines = append(lines, ligo.Variable{Type: ligo.TypeString, Value: input})
	}
	return ligo.Variable{Type: ligo.TypeString, Value: lines}
}

func vmInput(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) > 1 {
		vm.Throw(fmt.Sprintf("input: require less than 2 arguments, got %d arguments", len(a)))
	}

	if len(a) == 1 {
		vmPrint(vm, a...)
	}

	rd := bufio.NewReader(os.Stdin)

	input, err := rd.ReadString('\n')
	if err == io.EOF {
		return ligo.Variable{Type: ligo.TypeString, Value: ""}
	}

	if err != nil {
		vm.Throw(fmt.Sprintf("input : panicked when trying to read from stdin : %s", err))
	}

	if len(input) > 0 && input[len(input)-1] == '\n' {
		input = input[:len(input)-1]
	}

	return ligo.Variable{Type: ligo.TypeString, Value: input}
}

func vmSprintf(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) < 1 {
		vm.Throw("sprintf : expected atleast one argument")
	}

	if a[0].Type != ligo.TypeString {
		vm.Throw("sprintf : format expected as a string type")
	}

	values := collectVars(a[1:])

	return ligo.Variable{Type: ligo.TypeString, Value: fmt.Sprintf(a[0].Value.(string), values...)}

}

func collectVars(a []ligo.Variable) []interface{} {
	values := make([]interface{}, 0)
	for _, val := range a {
		if val.Type == ligo.TypeArray {
			values = append(values, collectVars(val.Value.([]ligo.Variable)))
			continue
		}
		values = append(values, val.Value)
	}

	return values
}

func vmOr(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) < 1 {
		vm.Throw("or : expected atleast one argument")
	}

	resultBool := true
	for _, val := range a {
		if a[0].Type != ligo.TypeBool {
			vm.Throw("or : expected only boolean arguments, got " + val.GetTypeString())
		}

		resultBool = resultBool || val.Value.(bool)
	}
	return ligo.Variable{Type: ligo.TypeBool, Value: resultBool}
}

func vmAnd(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) < 1 {
		vm.Throw("and : expected atleast one argument")
	}

	resultBool := true
	for _, val := range a {
		if a[0].Type != ligo.TypeBool {
			vm.Throw("and : expected only boolean arguments, got " + val.GetTypeString())
		}

		resultBool = resultBool && val.Value.(bool)
	}
	return ligo.Variable{Type: ligo.TypeBool, Value: resultBool}
}

func vmNot(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 1 {
		vm.Throw(fmt.Sprintf("not : expected one argument, got %d argument(s)", len(a)))
	}

	if a[0].Type != ligo.TypeBool {
		vm.Throw(fmt.Sprintf("not : expected one argument of boolean type, got type %s", a[0].GetTypeString()))
	}

	return ligo.Variable{Type: ligo.TypeBool, Value: !a[0].Value.(bool)}
}

func vmArrayAppend(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) < 2 {
		vm.Throw("wrong no. of arguments to the append function")
	}

	if a[0].Type < 0x100 && a[0].Type != ligo.TypeString {
		vm.Throw("append function's first argument should be of a array type, Got " + a[0].GetTypeString())
	}

	if a[0].Type == ligo.TypeString {
		str := a[0].Value.(string)
		for _, val := range a[1:] {
			if val.Type == ligo.TypeInt {
				str += string(val.Value.(int64))
				continue
			}
			if val.Type == ligo.TypeString {
				str += val.Value.(string)
				continue
			}
			vm.Throw(fmt.Sprintf("append : unable to append %s type to String", val.GetTypeString()))
		}
		return ligo.Variable{Type: ligo.TypeString, Value: str}
	}
	arrayReturn := make([]ligo.Variable, 0)

	arrayReturn = append(arrayReturn, a[0].Value.([]ligo.Variable)...)
	arrayReturn = append(arrayReturn, a[1:]...)
	return ligo.Variable{Type: a[0].Type, Value: arrayReturn}
}

func loadFile(fileName string, vm *ligo.VM) error {

	ltxtb, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	ltxt := ligo.StripComments(string(ltxtb))
	exps := make([]string, 0)
	line := 0
	inComment := false
	for i := 0; i < len(ltxt); i++ {
		ch := string(ltxt[i])
		switch ch {
		case "(":
			if inComment {
				continue
			}
			off := ligo.MatchChars(ltxt, int64(i), '(', ')') + 1
			if off < int64(i) {
				return fmt.Errorf("Syntax error near %d:%d : %s", i, line, ltxt[i:])
			}
			exps = append(exps, ltxt[i:off])
			i = int(off)
		case " ", "\n", "\r", "\t":
			if ch == "\n" || ch == "\r" {
				line++
				inComment = false
			}
			continue
		case ";":
			inComment = true
		default:
			if inComment {
				continue
			}
			return fmt.Errorf("unexpected Character at line %d : %s", line, ch)
		}
	}

	for _, val := range exps {
		_, err := vm.Eval(val)
		if err != nil {
			return fmt.Errorf("error : %s", err)
		}
	}

	return nil
}

func vmReciprocal(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 1 {
		vm.Throw("reciprocal : wrong number of arguments")
	}
	num := a[0]
	if num.Type != ligo.TypeInt && num.Type != ligo.TypeFloat {
		vm.Throw("reciprocal : expects a number type argument, got " + num.GetTypeString())
	}
	if num.Type == ligo.TypeFloat {
		return ligo.Variable{Type: ligo.TypeFloat, Value: 1 / num.Value.(float64)}
	}

	return ligo.Variable{Type: ligo.TypeFloat, Value: 1 / float64(num.Value.(int64))}
}

func vmCar(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 1 {
		vm.Throw("car can be done for one variable only")
	}
	if a[0].Type < 100 {
		if a[0].Type == ligo.TypeString {
			str := a[0].Value.(string)
			if str == "" {
				return ligo.Variable{Type: ligo.TypeNil, Value: nil}
			}
			return ligo.Variable{Type: ligo.TypeString, Value: string(str[0])}
		}
		vm.Throw("car can be done only for array or string type")
	}

	array := a[0].Value.([]ligo.Variable)
	if len(array) < 1 {
		return ligo.Variable{Type: ligo.TypeNil, Value: nil}
	}
	return array[0]
}

func vmCdr(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 1 {
		vm.Throw("cdr can be done for one variable only")
	}
	if a[0].Type < 0x100 {
		if a[0].Type == ligo.TypeString {
			str := a[0].Value.(string)
			if str == "" {
				return ligo.Variable{Type: ligo.TypeNil, Value: nil}
			}
			return ligo.Variable{Type: ligo.TypeString, Value: string(str[1:])}
		}
		vm.Throw(fmt.Sprint("cdr can be done only for array type", a[0]))
	}
	array := a[0].Value.([]ligo.Variable)
	if len(array) <= 1 {
		return ligo.Variable{Type: ligo.TypeNil, Value: nil}
	}
	return ligo.Variable{Type: ligo.TypeArray, Value: array[1:]}
}

func vmInEqualityGT(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 2 {
		vm.Throw("InEquality can be done for 2 numbers only")
	}
	if (a[0].Type != ligo.TypeInt && a[0].Type != ligo.TypeFloat) ||
		(a[1].Type != ligo.TypeInt && a[1].Type != ligo.TypeFloat) {
		vm.Throw("InEquality can be done for 2 numbers only")
	}
	var num1, num2 float64
	if a[0].Type == ligo.TypeInt {
		num1 = float64(a[0].Value.(int64))
	} else {
		num1 = a[0].Value.(float64)
	}

	if a[1].Type == ligo.TypeInt {
		num2 = float64(a[1].Value.(int64))
	} else {
		num2 = a[1].Value.(float64)
	}
	return ligo.Variable{Type: ligo.TypeBool, Value: num1 > num2}
}

func vmInEqualityGTEQ(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 2 {
		vm.Throw("InEquality can be done for 2 integers only")
	}
	if a[0].Type != ligo.TypeInt || a[1].Type != ligo.TypeInt {
		vm.Throw("InEquality can be done for 2 integers only")
	}
	num1 := a[0].Value.(int64)
	num2 := a[1].Value.(int64)
	return ligo.Variable{Type: ligo.TypeBool, Value: num1 >= num2}
}

func vmInEqualityLTEQ(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 2 {
		vm.Throw("InEquality can be done for 2 integers only")
	}
	if a[0].Type != ligo.TypeInt || a[1].Type != ligo.TypeInt {
		vm.Throw("InEquality can be done for 2 integers only")
	}
	num1 := a[0].Value.(int64)
	num2 := a[1].Value.(int64)
	return ligo.Variable{Type: ligo.TypeBool, Value: num1 <= num2}
}

func vmEquality(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 2 {
		vm.Throw("Equality can be done for 2 integers only")
	}
	if a[0].Type != a[1].Type {
		fmt.Println(a[0], a[1])
		vm.Throw(fmt.Sprintf("Equality can be done for 2 Values of same types only : found %s and %s, %s %s",
			a[0].GetTypeString(), a[1].GetTypeString(), a[0], a[1]))
	}
	return ligo.Variable{Type: ligo.TypeBool, Value: a[0] == a[1]}
}

func vmModulus(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 2 {
		vm.Throw("Modulus can be done for 2 integers only")
	}
	if a[0].Type != ligo.TypeInt || a[1].Type != ligo.TypeInt {
		vm.Throw("Modulus can be done for 2 integers only")
	}
	num1 := a[0].Value.(int64)
	num2 := a[1].Value.(int64)
	return ligo.Variable{Type: ligo.TypeInt, Value: num1 % num2}
}

func vmType(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 1 {
		vm.Throw("Keyword cannot take more than 1 argument")
	}

	return ligo.Variable{Type: ligo.TypeString, Value: a[0].GetTypeString()}
}

func vmPrint(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	for index, val := range a {
		if index != 0 {
			fmt.Print(" ")
		}
		switch true {
		case val.Type < 7:
			fmt.Print(val.Value)
		case val.Type == ligo.TypeArray:
			vmPrint(vm, val.Value.([]ligo.Variable)...)
		case val.Type == ligo.TypeMap:
			mm := val.Value.(ligo.Map)
			fmt.Print("{")
			for key, value := range mm {
				fmt.Print(key.Value, ":", value.Value, ";")
			}
			fmt.Print("}")
		}
	}
	return ligo.Variable{Type: ligo.TypeNil, Value: nil}
}

func vmPrintln(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	vmPrint(vm, a...)
	fmt.Println("")
	return ligo.Variable{Type: ligo.TypeNil, Value: nil}
}

func vmThrow(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 1 {
		vm.Throw("Cannot accept more that 1 variable")
	}
	vm.Throw(fmt.Sprint(a[0].Value))
	return ligo.Variable{Type: ligo.TypeNil, Value: nil}
}

func vmAdd(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	number := false
	tp := ligo.TypeInt
	sums := ""
	sumf := float64(0)
	for i, val := range a {
		if val.Type != ligo.TypeInt && val.Type != ligo.TypeFloat && val.Type != ligo.TypeString {
			vm.Throw("Cannot add a variable of type that is not String, Int or a Float.")
		}
		if i == 0 {
			if val.Type == ligo.TypeInt || val.Type == ligo.TypeFloat {
				number = true
			}
		}
		if number {
			if val.Type == ligo.TypeString {
				vm.Throw("Cannot add a string to a number")
			}
			switch val.Value.(type) {
			case int64:
				sumf += float64(val.Value.(int64))
			case float64:
				tp = ligo.TypeFloat
				sumf += val.Value.(float64)
			}
		} else {
			if val.Type == ligo.TypeInt || val.Type == ligo.TypeFloat {
				vm.Throw("Cannot add a number to a string")
			}
			sums += val.Value.(string)
		}
	}
	if !number {
		return ligo.Variable{Type: ligo.TypeString, Value: sums}
	}
	if tp == ligo.TypeInt {
		return ligo.Variable{Type: ligo.TypeInt, Value: int64(sumf)}
	}
	return ligo.Variable{Type: ligo.TypeFloat, Value: sumf}
}

func vmProd(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	var prodf float64 = 1
	var prodi int64 = 1
	float := false
	for _, val := range a {
		if val.Type != ligo.TypeInt && val.Type != ligo.TypeFloat {
			vm.Throw("Cannot use this type in product")
		}
		switch val.Value.(type) {
		case int64:
			if float {
				prodf *= float64(val.Value.(int64))
			} else {
				prodi *= val.Value.(int64)
			}
		case float64:
			if float {
				prodf *= val.Value.(float64)
			} else {
				prodf = float64(prodi) * val.Value.(float64)
				float = true
			}
		}
	}
	if float {
		return ligo.Variable{Type: ligo.TypeFloat, Value: prodf}
	}
	return ligo.Variable{Type: ligo.TypeInt, Value: prodi}
}

func vmArraySet(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 3 {
		vm.Throw(fmt.Sprintf("wrong number of parameters for the array-set function, required 3, got %d", len(a)))
	}
	array := a[0].Value.([]ligo.Variable)
	index := a[1].Value.(int64)

	if index >= int64(len(array)) || index < 0 {
		vm.Throw(fmt.Sprintf("index value of %d is invalid corresponding to the highest index %d", index, len(array)))
	}

	array[index] = a[2]

	return ligo.Variable{Type: a[0].Type, Value: array}
}

func vmLen(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 1 {
		vm.Throw("len can be done for one variable only")
	}
	if a[0].Type != ligo.TypeArray && a[0].Type != ligo.TypeString && a[0].Type != ligo.TypeMap {
		vm.Throw(fmt.Sprint("len can be done only for array type ", a[0].GetTypeString(), " ", a[0].Value))
	}
	if a[0].Type == ligo.TypeString {
		return ligo.Variable{Type: ligo.TypeInt, Value: int64(len(a[0].Value.(string)))}
	}
	if a[0].Type == ligo.TypeMap {
		return ligo.Variable{Type: ligo.TypeInt, Value: int64(len(a[0].Value.(ligo.Map)))}
	}
	return ligo.Variable{Type: ligo.TypeInt, Value: int64(len(a[0].Value.([]ligo.Variable)))}
}

func vmSleep(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 1 {
		vm.Throw("sleep expects only one argument")
	}
	if a[0].Type != ligo.TypeInt {
		vm.Throw("sleep expects only integers")
	}
	time.Sleep(time.Duration(a[0].Value.(int64)) * time.Millisecond)
	return ligo.Variable{Type: ligo.TypeNil, Value: nil}
}

func vmIsNil(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 1 {
		vm.Throw("is-nil : expects only one argument")
	}

	if a[0].Type == ligo.TypeNil || a[0].Value == nil {
		return ligo.Variable{Type: ligo.TypeBool, Value: true}
	}

	return ligo.Variable{Type: ligo.TypeBool, Value: false}
}
