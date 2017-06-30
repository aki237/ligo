package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"

	"github.com/aki237/ligo"
)

func vmConnect(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 1 {
		panic("split : wrong number of arguments.")
	}
	if a[0].Type != ligo.TYPE_String {
		panic("split : expects string to be splitted")
	}
	ipPort := a[0].Value.(string)
	conn, err := net.Dial("tcp", ipPort)
	if err != nil {
		panic(err)
	}
	return ligo.Variable{Type: ligo.TYPE_ReadWriteCloser, Value: conn}
}

func vmSplit(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 2 {
		panic("split : wrong number of arguments.")
	}
	if a[0].Type != ligo.TYPE_String {
		panic("split : expects string to be splitted")
	}

	if a[1].Type != ligo.TYPE_String {
		panic("split : expects string to be splitted with")
	}
	splitted := strings.Split(a[0].Value.(string), a[1].Value.(string))
	varVals := make([]ligo.Variable, 0)

	for _, value := range splitted {
		tmp := ligo.Variable{Type: ligo.TYPE_String, Value: value}
		varVals = append(varVals, tmp)
	}
	return ligo.Variable{Type: ligo.TYPE_MonoTypeArray, Value: varVals}
}

func vmOpen(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 2 {
		panic("open : wrong number of arguments.")
	}
	if a[0].Type != ligo.TYPE_String {
		panic("open : expects filename as a string")
	}

	if a[1].Type != ligo.TYPE_String {
		panic("open : expects mode as a string")
	}
	filename := a[0].Value.(string)
	mode := a[0].Value.(string)

	flags := os.O_RDONLY
	tp := ligo.TYPE_Reader
	switch mode {
	case "r":
		flags = os.O_RDONLY
		tp = ligo.TYPE_ReadCloser
	case "w":
		flags = os.O_WRONLY | os.O_CREATE
		tp = ligo.TYPE_Writer
	case "rw", "wr":
		flags = os.O_RDWR | os.O_CREATE
		tp = ligo.TYPE_ReadWriteCloser
	case "a":
		flags = os.O_APPEND | os.O_CREATE
		tp = ligo.TYPE_ReadWriteCloser
	}
	file, err := os.OpenFile(filename, flags, 0644)
	if err != nil {
		panic(err)
	}
	return ligo.Variable{Type: tp, Value: file}
}

func vmRead(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 2 {
		panic("read : wrong number of arguments.")
	}
	if a[0].Type != ligo.TYPE_Reader && a[0].Type <= 0x700 {
		panic(fmt.Sprintf("read : expects a reader interface, got : 0x%X , %T", a[0].Type, a[0].Value))
	}
	rd := a[0].Value.(io.Reader)
	if a[1].Type != ligo.TYPE_Int {
		panic("read : expects int for amount to be read got : " + fmt.Sprint(a[1]))
	}
	bytes := a[1].Value.(int64)
	bs := make([]byte, bytes)
	n, err := rd.Read(bs)
	if err != nil {
		panic(err)
	}
	return ligo.Variable{Type: ligo.TYPE_String, Value: string(bs[:n])}
}

func vmWrite(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 2 {
		panic("write : wrong number of arguments.")
	}
	if a[0].Type != ligo.TYPE_Writer && a[0].Type < 0x703 {
		panic(fmt.Sprintf("write : expects a writer interface, got : 0x%X , %T", a[0].Type, a[0].Value))
	}
	if a[1].Type != ligo.TYPE_String {
		panic("write : expects data as a string")
	}
	wr := a[0].Value.(io.Writer)
	cn := a[1].Value.(string)
	n, err := wr.Write([]byte(cn))
	if err != nil {
		panic(err)
	}
	return ligo.Variable{Type: ligo.TYPE_Int, Value: int64(n)}
}

func vmCar(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 1 {
		panic("car can be done for one variable only")
	}
	if a[0].Type < 100 {
		panic("car can be done only for array type")
	}
	array := a[0].Value.([]ligo.Variable)
	if len(array) < 1 {
		return ligo.Variable{Type: ligo.TYPE_Nil, Value: nil}
	}
	return array[0]
}

func vmCdr(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 1 {
		panic("cdr can be done for one variable only")
	}
	if a[0].Type < 0x100 {
		panic(fmt.Sprint("cdr can be done only for array type", a[0]))
	}
	array := a[0].Value.([]ligo.Variable)
	if len(array) <= 1 {
		return ligo.Variable{Type: ligo.TYPE_Nil, Value: nil}
	}
	return ligo.Variable{Type: ligo.TYPE_MonoTypeArray, Value: array[1:]}
}

func vmInEqualityGT(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 2 {
		panic("InEquality can be done for 2 integers only")
	}
	if a[0].Type != ligo.TYPE_Int || a[1].Type != ligo.TYPE_Int {
		panic("InEquality can be done for 2 integers only")
	}
	num1 := a[0].Value.(int64)
	num2 := a[1].Value.(int64)
	return ligo.Variable{Type: ligo.TYPE_Bool, Value: num1 > num2}
}

func vmInEqualityGTEQ(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 2 {
		panic("InEquality can be done for 2 integers only")
	}
	if a[0].Type != ligo.TYPE_Int || a[1].Type != ligo.TYPE_Int {
		panic("InEquality can be done for 2 integers only")
	}
	num1 := a[0].Value.(int64)
	num2 := a[1].Value.(int64)
	return ligo.Variable{Type: ligo.TYPE_Bool, Value: num1 >= num2}
}

func vmInEqualityLTEQ(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 2 {
		panic("InEquality can be done for 2 integers only")
	}
	if a[0].Type != ligo.TYPE_Int || a[1].Type != ligo.TYPE_Int {
		panic("InEquality can be done for 2 integers only")
	}
	num1 := a[0].Value.(int64)
	num2 := a[1].Value.(int64)
	return ligo.Variable{Type: ligo.TYPE_Bool, Value: num1 <= num2}
}

func vmEquality(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 2 {
		panic("Equality can be done for 2 integers only")
	}
	if a[0].Type != a[1].Type {
		panic(fmt.Sprintf("Equality can be done for 2 Values of same types only : found %s and %s",
			getTypeString(a[0].Type), getTypeString(a[1].Type)))
	}
	return ligo.Variable{Type: ligo.TYPE_Bool, Value: a[0] == a[1]}
}

func vmModulus(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 2 {
		panic("Modulus can be done for 2 integers only")
	}
	if a[0].Type != ligo.TYPE_Int || a[1].Type != ligo.TYPE_Int {
		panic("Modulus can be done for 2 integers only")
	}
	num1 := a[0].Value.(int64)
	num2 := a[1].Value.(int64)
	return ligo.Variable{Type: ligo.TYPE_Int, Value: num1 % num2}
}

func vmType(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 1 {
		panic("Keyword cannot take more than 1 argument")
	}

	return ligo.Variable{Type: ligo.TYPE_String, Value: getTypeString(a[0].Type)}
}

func getTypeString(t ligo.Type) (tp string) {
	tp = ""
	switch t {
	case ligo.TYPE_Int:
		tp = "int"
	case ligo.TYPE_Float:
		tp = "float"
	case ligo.TYPE_Bool:
		tp = "bool"
	case ligo.TYPE_String:
		tp = "string"
	case ligo.TYPE_Nil:
		tp = "nil"
	case ligo.TYPE_DFunc, ligo.TYPE_IFunc:
		tp = "func"
	case ligo.TYPE_MonoTypeArray:
		tp = "array<mono type>"
	case ligo.TYPE_PolyTypeArray:
		tp = "array<poly type>"
	case ligo.TYPE_Map:
		tp = "map"
	}
	return
}

func vmPrint(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	for index, val := range a {
		if index != 0 {
			fmt.Print(" ")
		}
		fmt.Print(val.Value)
	}
	return ligo.Variable{Type: ligo.TYPE_Nil, Value: nil}
}

func vmPrintln(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	for index, val := range a {
		if index != 0 {
			fmt.Print(" ")
		}
		if val.Type < 7 {
			fmt.Print(val.Value)
		} else {
			vmPrint(vm, val.Value.([]ligo.Variable)...)
		}
	}
	fmt.Println("")
	return ligo.Variable{Type: ligo.TYPE_Nil, Value: nil}
}

func vmPanic(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 1 {
		panic("Cannot accept more that 1 variable")
	}
	panic(a[0].Value)
	return ligo.Variable{Type: ligo.TYPE_Nil, Value: nil}
}

func vmAdd(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	number := false
	tp := ligo.TYPE_Int
	var sums string = ""
	var sumf float64 = 0
	var sumi int64 = 0
	for i, val := range a {
		if val.Type != ligo.TYPE_Int && val.Type != ligo.TYPE_Float && val.Type != ligo.TYPE_String {
			panic("Cannot add a variable of type that is not String, Int or a Float.")
		}
		if i == 0 {
			if val.Type == ligo.TYPE_Int || val.Type == ligo.TYPE_Float {
				number = true
			}
		}
		if number {
			if val.Type == ligo.TYPE_String {
				panic("Cannot add a string to a number")
			}
			if val.Type == ligo.TYPE_Float {
				tp = ligo.TYPE_Float
				sumf = float64(sumi)
			}
			if tp == ligo.TYPE_Int {
				sumi += val.Value.(int64)
			}
			if tp == ligo.TYPE_Float {
				switch val.Value.(type) {
				case int64:
					sumf += float64(val.Value.(int64))
				case float64:
					sumf += val.Value.(float64)
				}
			}
		} else {
			if val.Type == ligo.TYPE_Int || val.Type == ligo.TYPE_Float {
				panic("Cannot add a number to a string")
			}
			sums += val.Value.(string)
		}
	}
	if !number {
		return ligo.Variable{Type: ligo.TYPE_String, Value: sums}
	}
	if tp == ligo.TYPE_Int {
		return ligo.Variable{Type: ligo.TYPE_Int, Value: sumi}
	}
	return ligo.Variable{Type: ligo.TYPE_Float, Value: sumf}
}

func vmProd(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	var prodf float64 = 1
	var prodi int64 = 1
	float := false
	for _, val := range a {
		if val.Type != ligo.TYPE_Int && val.Type != ligo.TYPE_Float {
			panic("Cannot use this type in product")
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
		return ligo.Variable{Type: ligo.TYPE_Float, Value: prodf}
	}
	return ligo.Variable{Type: ligo.TYPE_Int, Value: prodi}
}

func vmLen(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 1 {
		panic("len can be done for one variable only")
	}
	if a[0].Type < 0x100 && a[0].Type != ligo.TYPE_String {
		panic(fmt.Sprint("len can be done only for array type", getTypeString(a[0].Type)))
	}
	if a[0].Type == ligo.TYPE_String {
		return ligo.Variable{Type: ligo.TYPE_Int, Value: int64(len(a[0].Value.(string)))}
	}
	return ligo.Variable{Type: ligo.TYPE_Int, Value: int64(len(a[0].Value.([]ligo.Variable)))}
}

func vmSleep(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 1 {
		panic("sleep expects only one argument")
	}
	if a[0].Type != ligo.TYPE_Int {
		panic("sleep expects only integers")
	}
	time.Sleep(time.Duration(a[0].Value.(int64)) * time.Second)
	return ligo.Variable{Type: ligo.TYPE_Nil, Value: nil}
}
