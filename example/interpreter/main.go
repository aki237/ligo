package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aki237/ligo"
)

func run() {
	vm := ligo.NewVM()
	vm.Funcs["println"] = vmPrintln
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
	vm.Funcs["panic"] = vmPanic
	vm.Funcs["sleep"] = vmSleep
	vm.Funcs["open"] = vmOpen
	vm.Funcs["read"] = vmRead
	vm.Funcs["write"] = vmWrite
	vm.Funcs["split"] = vmSplit
	vm.Funcs["net.connect"] = vmConnect

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage %s [filename.lf]\n", os.Args[0])
		return
	}
	ltxtb, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error : %s\n", err)
		return
	}
	ltxt := string(ltxtb)
	exps := make([]string, 0)
	line := 0
	for i := 0; i < len(ltxt); i++ {
		ch := string(ltxt[i])
		switch ch {
		case "(":
			off := ligo.MatchChars(ltxt, int64(i), '(', ')') + 1
			exps = append(exps, ltxt[i:off])
			i = int(off)
		case " ", "\n", "\r", "\t":
			if ch == "\n" || ch == "\r" {
				line += 1
			}
			continue
		default:
			fmt.Fprintf(os.Stderr, "Unexpected Character at line %d : %s\n", line, ch)
			return
		}
	}

	for _, val := range exps {
		vl, err := vm.Eval(val)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error : %s\n", err)
			break
		}
		if vl.Type != ligo.TYPE_Nil {
			fmt.Println(vl.Value)
		}
	}
}

func main() {
	run()
}
