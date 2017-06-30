package ligo

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var r_integer *regexp.Regexp = regexp.MustCompile("^[+-]?[0-9]+$")
var r_float *regexp.Regexp = regexp.MustCompile("^[+-]?[0-9]*\\.[0-9]+$")
var r_string *regexp.Regexp = regexp.MustCompile("^\".*\"$")
var r_variable *regexp.Regexp = regexp.MustCompile("^[[:alpha:]]+[[:alnum:]]*$")
var r_expression *regexp.Regexp = regexp.MustCompile("^\\(.*\\)$")
var r_closure *regexp.Regexp = regexp.MustCompile("^\\|([[:alpha:]]+[[:alnum:]]*\\s*)*\\|$")
var r_array *regexp.Regexp = regexp.MustCompile("^\\[(.*\\s*)*\\]$")

type Variable struct {
	Type  Type
	Value interface{}
}

func (v Variable) String() string {
	typeString := "Variable {Type : "
	switch v.Type {
	case TYPE_Int:
		typeString += "Integer<64>"
	case TYPE_Float:
		typeString += "Float<64>"
	case TYPE_String:
		typeString += "String"
	case TYPE_Bool:
		typeString += "Boolean"
	case TYPE_MonoTypeArray:
		typeString += "Array<MonoType>"
	case TYPE_PolyTypeArray:
		typeString += "Array<PolyType>"
	case TYPE_Nil:
		typeString += "Nil"
	}
	return typeString + fmt.Sprint(" ,Value : ", v.Value, "}")
}

type Defined struct {
	scopevars []string
	eval      string
}

type InBuilt func(*VM, ...Variable) Variable

type VM struct {
	global *VM
	Vars   map[string]Variable
	Funcs  map[string]InBuilt
	LFuncs map[string]Defined
}

func NewVM() *VM {
	vm := &VM{}
	vm.Vars = make(map[string]Variable, 0)
	vm.Funcs = make(map[string]InBuilt, 0)
	vm.LFuncs = make(map[string]Defined, 0)
	vm.global = nil
	return vm
}

func (vm *VM) GetVariable(token string) (Variable, error) {
	v := ligoNil
	if len(token) < 1 {
		return ligoNil, LigoError("Invalid Token passed")
	}
	switch true {
	case r_array.MatchString(token):
		ar := token[1:MatchChars(token, 0, '[', ']')]
		tkns, err := ScanTokens("(" + ar + ")")
		if err != nil {
			return ligoNil, err
		}
		var tp Type
		vars := make([]Variable, 0)
		for i, val := range tkns {
			v, err := vm.GetVariable(val)
			if err != nil {
				return ligoNil, err
			}
			vars = append(vars, v)
			if i == 0 {
				tp = v.Type
				continue
			}
			if tp != v.Type {
				tp = TYPE_PolyTypeArray
			}
		}
		if tp != TYPE_PolyTypeArray {
			tp = TYPE_MonoTypeArray | tp
		}
		retVars := Variable{Type: tp, Value: vars}
		return retVars, nil
	case r_expression.MatchString(token) || token[0] == '(':
		var err error
		v, err = vm.Eval(token)
		if err != nil {
			return ligoNil, err
		}
	case r_integer.MatchString(token):
		num, err := strconv.ParseInt(token, 10, 64)
		if err != nil {
			fmt.Println(err)
		}
		v = Variable{Type: TYPE_Int, Value: num}
	case r_float.MatchString(token):
		num, _ := strconv.ParseFloat(token, 64)
		v = Variable{Type: TYPE_Float, Value: num}
	case r_string.MatchString(token) || token[0] == '"':
		v = Variable{Type: TYPE_String, Value: token[1 : len(token)-1]}
	case token == "true":
		v = Variable{Type: TYPE_Bool, Value: true}
	case token == "false":
		v = Variable{Type: TYPE_Bool, Value: false}
	case r_variable.MatchString(token):
		varFromVM, ok := vm.Vars[token]
		if ok {
			v = Variable{Type: varFromVM.Type, Value: varFromVM.Value}
			break
		}
		if fnc, ok := vm.Funcs[token]; ok {
			v = Variable{Type: TYPE_IFunc, Value: fnc}
			break
		}
		if fnc, ok := vm.LFuncs[token]; ok {
			v = Variable{Type: TYPE_DFunc, Value: fnc}
			break
		}
		v = ligoNil
	default:
		if fnc, ok := vm.Funcs[token]; ok {
			v = Variable{Type: TYPE_IFunc, Value: fnc}
			break
		}
		if fnc, ok := vm.LFuncs[token]; ok {
			v = Variable{Type: TYPE_DFunc, Value: fnc}
			break
		}
	}
	if v == ligoNil && vm.global != nil {
		return vm.global.GetVariable(token)
	}
	return v, nil
}

func (vm *VM) setFn(tokens []string) (Variable, error) {
	if len(tokens) != 4 {
		return ligoNil, LigoError("A function construct can only have a single returning function")
	}
	fnName := tokens[1]
	if _, ok := vm.Funcs[fnName]; ok {
		fmt.Printf("Warning : function \"%s\" has already been declared as an InBuilt function.\n", fnName)
	}
	if _, ok := vm.LFuncs[fnName]; ok {
		fmt.Printf("Warning : function \"%s\" has already been declared as an Ligo function.\n", fnName)
	}
	if !r_closure.MatchString(tokens[2]) {
		return ligoNil,
			LigoError("Expected parameter name in the function definition " + fnName + " closure : " + tokens[2])
	}
	varNames := getVarsFromClosure(tokens[2])
	fn := Defined{scopevars: varNames, eval: tokens[3]}
	vm.LFuncs[fnName] = fn
	return ligoNil, nil
}

func (vm *VM) setVar(tokens []string) (Variable, error) {
	if len(tokens) != 3 {
		return ligoNil, LigoError("Wrong number of arguments to the keyword.")
	}
	if !r_variable.MatchString(tokens[1]) {
		return ligoNil, LigoError("Wrong token found in the variable name")
	}
	v, err := vm.GetVariable(tokens[2])
	if err != nil {
		return ligoNil, err
	}

	switch v.Type {
	case TYPE_IFunc:
		_, ok := vm.Funcs[tokens[1]]
		if !ok {
			return ligoNil, LigoError("Variable not defined. Try \"var\" for creating a new variable")
		}
		vm.Funcs[tokens[1]] = v.Value.(InBuilt)
		return ligoNil, nil
	case TYPE_DFunc:
		_, ok := vm.LFuncs[tokens[1]]
		if !ok {
			return ligoNil, LigoError("Variable not defined. Try \"var\" for creating a new variable")
		}
		vm.LFuncs[tokens[1]] = v.Value.(Defined)
		return ligoNil, nil
	}

	_, ok := vm.Vars[tokens[1]]
	if ok {
		vm.Vars[tokens[1]] = v
		return ligoNil, nil
	}
	if vm.global == nil {
		return ligoNil, LigoError("Variable '" + tokens[1] + "' not defined. Try \"var\" for creating a new variable")
	}
	return vm.global.setVar(tokens)
}

func (vm *VM) newVar(tokens []string) (Variable, error) {
	if len(tokens) != 3 {
		return ligoNil, LigoError("Wrong number of arguments to the keyword.")
	}
	if !r_variable.MatchString(tokens[1]) {
		return ligoNil, LigoError("Wrong token found in the variable name")
	}
	v, err := vm.GetVariable(tokens[2])
	if err != nil {
		return ligoNil, err
	}
	switch v.Type {
	case TYPE_IFunc:
		vm.Funcs[tokens[1]] = v.Value.(InBuilt)
		return ligoNil, nil
	case TYPE_DFunc:
		vm.LFuncs[tokens[1]] = v.Value.(Defined)
		return ligoNil, nil
	}
	_, ok := vm.Vars[tokens[1]]
	if ok {
		return ligoNil, LigoError("Variable '" + tokens[1] + "' already defined. Try \"set\" for updating variables")
	}
	vm.Vars[tokens[1]] = v
	return ligoNil, nil
}

func (vm *VM) getInBuiltFunction(fnName string) (InBuilt, bool) {
	fn, found := vm.Funcs[fnName]
	return fn, found
}

func (vm *VM) runInBuiltFunction(function InBuilt, vars []Variable) (Variable, error) {
	return function(vm, vars...), nil
}

func (vm *VM) getDefinedFunction(fnName string) (Defined, bool) {
	fn, found := vm.LFuncs[fnName]
	return fn, found
}

func (vm *VM) runDefinedFunction(function Defined, fnName string, vars []Variable) (Variable, error) {
	if len(vars) != len(function.scopevars) {
		return ligoNil, LigoError(fmt.Sprintf("Expected %d arguments, got %d for the %s function",
			len(function.scopevars),
			len(vars),
			fnName,
		))
	}
	nvm := vm.NewScope()
	for i, val := range function.scopevars {
		switch vars[i].Type {
		case TYPE_IFunc:
			nvm.Funcs[val] = vars[i].Value.(InBuilt)
		case TYPE_DFunc:
			nvm.LFuncs[val] = vars[i].Value.(Defined)
		default:
			nvm.Vars[val] = vars[i]
		}
	}
	return nvm.Eval(function.eval)
}

func (vm *VM) run(tkns []string) (Variable, error) {
	vars := make([]Variable, 0)
	fnName := tkns[0]
	for i := 1; i < len(tkns); i++ {
		v, err := vm.GetVariable(tkns[i])
		if err != nil {
			return ligoNil, err
		}

		vars = append(vars, v)
	}
	if function, ok := vm.getInBuiltFunction(fnName); ok {
		return vm.runInBuiltFunction(function, vars)
	}
	if function, ok := vm.getDefinedFunction(fnName); ok {
		return vm.runDefinedFunction(function, fnName, vars)
	}
	if vm.global == nil {
		return ligoNil, LigoError("Function '" + fnName + "' not found")
	}
	if function, ok := vm.global.getInBuiltFunction(fnName); ok {
		return vm.runInBuiltFunction(function, vars)
	}
	function, ok := vm.global.getDefinedFunction(fnName)
	if !ok {
		return ligoNil, LigoError("Function '" + fnName + "' not found")
	}
	return vm.runDefinedFunction(function, fnName, vars)
}

func (vm *VM) runLoop(tkns []string) (Variable, error) {
	if len(tkns) != 3 {
		return ligoNil, LigoError("Illegal loop construct. Can take 3 arguments only.")
	}
	condition := tkns[1]
	runExp := tkns[2]
	result, err := vm.Eval(condition)
	if err != nil {
		return ligoNil, err
	}
	if result.Type != TYPE_Bool {
		return ligoNil, LigoError("Expected boolean return from the expression : " + condition)
	}
	for result.Value.(bool) {
		_, err := vm.Eval(runExp)
		if err != nil {
			return ligoNil, err
		}
		result, err = vm.Eval(condition)
		if err != nil {
			return ligoNil, err
		}
		if result.Type != TYPE_Bool {
			return ligoNil, LigoError("Expected boolean return from the expression : " + condition)
		}
	}
	return ligoNil, err
}

func (vm *VM) ifClause(tkns []string) (Variable, error) {
	if len(tkns) > 4 || len(tkns) < 3 {
		return ligoNil, LigoError("Illegal if construct. Can take 3 or 4 arguments.")
	}
	condition := tkns[1]
	boolVar, ok := vm.Vars[condition]
	if condition != "true" && condition != "false" && !r_expression.MatchString(condition) && !ok {
		return ligoNil,
			LigoError("Expected a boolean value or expression for the if clause condition, got : " + condition)
	}

	successClause := tkns[2]
	failureClause := ""
	if len(tkns) == 4 {
		failureClause = tkns[3]
	}
	var result Variable
	var err error
	if !ok {
		result, err = vm.Eval(condition)
		if err != nil {
			return ligoNil, err
		}
	} else {
		result = boolVar
	}
	if result.Type != TYPE_Bool {
		return ligoNil, LigoError("Expected boolean return from the expression : " + condition)
	}
	if !result.Value.(bool) {
		if failureClause == "" {
			return ligoNil, nil
		}
		return vm.Eval(failureClause)
	}
	return vm.Eval(successClause)
}

func (vm *VM) returnArg(tkns []string) (Variable, error) {
	if len(tkns) != 2 {
		panic("Cannot return more than 2 values. (Atleast for now.)")
	}
	return vm.GetVariable(tkns[1])
}

func (vm *VM) fork(tkns []string) (Variable, error) {
	if len(tkns) != 2 {
		return ligoNil, LigoError("Expected one expression, got " + fmt.Sprint(len(tkns)) + " arguments")
	}
	go vm.Eval(tkns[1])
	return ligoNil, nil
}

func (vm *VM) runExpressions(tkns []string) (Variable, error) {
	v := ligoNil
	for i, val := range tkns {
		if i == 0 || val == "" {
			continue
		}
		vl, err := vm.Eval(val)
		if err != nil {
			return ligoNil, err
		}
		if i == len(tkns)-1 {
			v = vl
		}
	}
	return v, nil
}

func (vm *VM) evalString(tkns []string) (Variable, error) {
	if len(tkns) != 2 {
		return ligoNil, LigoError("'eval' keyword only accepts 1 argument")
	}
	vl, err := vm.GetVariable(tkns[1])
	if err != nil {
		return ligoNil, err
	}
	if vl.Type != TYPE_String {
		return ligoNil, LigoError("'eval' keyword only expression string")
	}
	exps := make([]string, 0)
	line := 0
	ltxt := vl.Value.(string)
	for i := 0; i < len(ltxt); i++ {
		ch := string(ltxt[i])
		switch ch {
		case "(":
			off := MatchChars(ltxt, int64(i), '(', ')') + 1
			exps = append(exps, ltxt[i:off])
			i = int(off)
		case " ", "\n", "\r", "\t":
			if ch == "\n" || ch == "\r" {
				line += 1
			}
			continue
		default:
			return ligoNil, LigoError(fmt.Sprintf("Unexpected Character at line %d : %s\n", line, ch))
		}
	}
	var retVal Variable
	for _, val := range exps {
		var err error
		retVal, err = vm.Eval(val)
		if err != nil {
			return ligoNil, err
		}
	}
	return retVal, nil
}

func (vm *VM) Eval(stmt string) (Variable, error) {
	stmt = strings.TrimSpace(stmt)
	if len(stmt) < 2 {
		return ligoNil, LigoError("Expected atleast (), got : " + stmt)
	}
	if !r_expression.MatchString(stmt) && stmt[0] != '(' {
		return vm.GetVariable(stmt)
	}
	tkns, err := ScanTokens(stmt)
	if err != nil {
		return ligoNil, err
	}
	if len(tkns) < 1 {
		return ligoNil, nil
	}
	fnName := tkns[0]
	switch fnName {
	case "var":
		return vm.newVar(tkns)
	case "set":
		return vm.setVar(tkns)
	case "fn":
		return vm.setFn(tkns)
	case "return":
		return vm.returnArg(tkns)
	case "progn":
		return vm.runExpressions(tkns)
	case "loop":
		return vm.runLoop(tkns)
	case "if":
		return vm.ifClause(tkns)
	case "eval":
		return vm.evalString(tkns)
	case "fork":
		return vm.fork(tkns)
	default:
		return vm.run(tkns)
	}
	return ligoNil, nil
}

func (vm *VM) Clone() *VM {
	nvm := NewVM()
	for key, value := range vm.Funcs {
		nvm.Funcs[key] = value
	}
	for key, value := range vm.LFuncs {
		nvm.LFuncs[key] = value
	}
	for key, value := range vm.Vars {
		nvm.Vars[key] = value
	}
	return nvm
}

func (vm *VM) NewScope() *VM {
	nvm := NewVM()
	if vm.global == nil {
		nvm.global = vm
	} else {
		nvm.global = vm.global
	}
	return nvm
}
