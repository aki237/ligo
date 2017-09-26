package ligo

import (
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

// escape sequences to be replaced with the counterpart in a string
var escapeSequences = map[string]int{
	"\\\\": 0x5C,
	"\\a":  0x07,
	"\\b":  0x08,
	"\\f":  0x0C,
	"\\n":  0x0A,
	"\\r":  0x0D,
	"\\t":  0x09,
	"\\v":  0x0B,
	"\\'":  0x27,
	"\"":   0x22,
}

// TODO : add escape sequence handling for hex and octal sequences...
func reformEscapes(str string) string {
	ret := ""
	isEscape := false
	for _, val := range str {
		switch val {
		case '\\':
			if !isEscape {
				isEscape = true
			} else {
				ret += string(0x5C)
				isEscape = false
			}
		default:
			if !isEscape {
				ret += string(val)
			} else {
				es := "\\" + string(val)
				num, ok := escapeSequences[es]
				if !ok {
					panic("in :\n\t" + str + "\nUnknown Escape sequence. : '\\" + string(val) + "'")
				}
				ret += string(num)
				isEscape = false
			}
		}
	}
	return ret
}

// regexp variables for matching the syntax of the script
var rInteger = regexp.MustCompile("^[+-]?[0-9]+$")
var rFloat = regexp.MustCompile("^[+-]?[0-9]*\\.[0-9]+$")
var rString = regexp.MustCompile("^\".*\"$")
var rVariable = regexp.MustCompile("^[[:alpha:]]+[[:alnum:]]*$")
var rExpression = regexp.MustCompile("^\\(.*\\)$")
var rClosure = regexp.MustCompile("^\\|([[:alpha:]]+[[:alnum:]]*\\s*)*\\.\\.\\.([[:alpha:]]+[[:alnum:]]*\\s*){0,1}|$")
var rArray = regexp.MustCompile("^\\[(.*\\s*)*\\]$")

// Variable is a struct denoting a value in the VM
type Variable struct {
	Type  Type
	Value interface{}
}

// GetTypeString method returns a string value corresponding to the type of it's value
func (v Variable) GetTypeString() (tp string) {
	tp = ""
	switch v.Type {
	case TYPE_Int:
		tp = "int"
	case TYPE_Float:
		tp = "float"
	case TYPE_Bool:
		tp = "bool"
	case TYPE_String:
		tp = "string"
	case TYPE_Nil:
		tp = "nil"
	case TYPE_DFunc, TYPE_IFunc:
		tp = "func"
	case TYPE_MonoTypeArray:
		tp = "array<mono type>"
	case TYPE_PolyTypeArray:
		tp = "array<poly type>"
	case TYPE_Map:
		tp = "map"
	}
	return
}

// String method implements the Stringer interface for the Variable type
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

// Defined struct contains variables needed for storing a function defined in ligo script itself
type Defined struct {
	scopevars []string
	eval      string
}

// InBuilt type is a function format that is callable from the ligo script
type InBuilt func(*VM, ...Variable) Variable

// VM struct is a State Struct contains all the variable maps,
// defined function maps, in-built function maps and a global
// scope pointing to the global Scope VM
type VM struct {
	global *VM
	Vars   map[string]Variable
	Funcs  map[string]InBuilt
	LFuncs map[string]Defined
}

// NewVM returns a new VM object pointer after initializing the values
func NewVM() *VM {
	vm := &VM{}
	vm.Vars = make(map[string]Variable, 0)
	vm.Funcs = make(map[string]InBuilt, 0)
	vm.LFuncs = make(map[string]Defined, 0)
	vm.global = nil
	return vm
}

// GetVariable method is used to process the token string passed and get the
// corresponding value from the VM's memory. This is a crucial function
// as, if the token passed is a sub expression this method knows to evaluate and
// return the value of that sub expression.
func (vm *VM) GetVariable(token string) (Variable, error) {
	v := ligoNil
	if len(token) < 1 {
		return ligoNil, LigoError("Invalid Token passed")
	}
	switch true {
	case rArray.MatchString(token):
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
	case rExpression.MatchString(token) || MatchChars(token, 0, '(', ')') > 0:
		var err error
		v, err = vm.Eval(token)
		if err != nil {
			return ligoNil, err
		}
	case rInteger.MatchString(token):
		num, err := strconv.ParseInt(token, 10, 64)
		if err != nil {
			fmt.Println(err)
		}
		v = Variable{Type: TYPE_Int, Value: num}
	case rFloat.MatchString(token):
		num, _ := strconv.ParseFloat(token, 64)
		v = Variable{Type: TYPE_Float, Value: num}
	case rString.MatchString(token) || token[0] == '"':
		token = reformEscapes(token)
		v = Variable{Type: TYPE_String, Value: token[1 : len(token)-1]}
	case token == "true":
		v = Variable{Type: TYPE_Bool, Value: true}
	case token == "false":
		v = Variable{Type: TYPE_Bool, Value: false}
	case rVariable.MatchString(token):
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

// setFn is used to parse a ligo function construct and store it in the
// current scope. It also warns if the function is already declared.
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
	if !rClosure.MatchString(tokens[2]) {
		return ligoNil,
			LigoError("Expected parameter name in the function definition " + fnName + " closure : " + tokens[2])
	}
	varNames := getVarsFromClosure(tokens[2])
	fn := Defined{scopevars: varNames, eval: tokens[3]}
	vm.LFuncs[fnName] = fn
	return ligoNil, nil
}

// setVar method is used to set a value to a variable.
// If the variable is not defined already, this will throw an error.
func (vm *VM) setVar(tokens []string) (Variable, error) {
	if len(tokens) != 3 {
		return ligoNil, LigoError("Wrong number of arguments to the keyword.")
	}
	if !rVariable.MatchString(tokens[1]) {
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

// newVar method is used to declare a new variable in the VM and set a value to it.
func (vm *VM) newVar(tokens []string) (Variable, error) {
	if len(tokens) != 3 {
		return ligoNil, LigoError("Wrong number of arguments to the keyword.")
	}
	if !rVariable.MatchString(tokens[1]) {
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

// getInBuiltFunction method is a small helper method to get the inbuilt function.
// If found it returns the inbuilt and true, else nil and false.
func (vm *VM) getInBuiltFunction(fnName string) (InBuilt, bool) {
	fn, found := vm.Funcs[fnName]
	return fn, found
}

// runInBuiltFunction method is a small helper method to run the passed inbuilt function
// with the passed variables.
func (vm *VM) runInBuiltFunction(function InBuilt, vars []Variable) (Variable, error) {
	return function(vm, vars...), nil
}

// getDefinedFunction method is a small helper method to get the defined function.
// If found it returns the inbuilt and true, else nil and false.
func (vm *VM) getDefinedFunction(fnName string) (Defined, bool) {
	fn, found := vm.LFuncs[fnName]
	return fn, found
}

// runDefinedFunction method is a helper method used to run a passed defined function with passed vars
func (vm *VM) runDefinedFunction(function Defined, fnName string, vars []Variable) (Variable, error) {
	if len(vars) < len(function.scopevars)-1 {
		return ligoNil, LigoError(fmt.Sprintf("Expected %d arguments, got %d for the %s function",
			len(function.scopevars),
			len(vars),
			fnName,
		))
	}

	if len(function.scopevars) > 0 && !isVariate(function.scopevars[len(function.scopevars)-1]) {
		if len(vars) != len(function.scopevars) {
			return ligoNil, LigoError(fmt.Sprintf("Expected %d arguments, got %d for the %s function",
				len(function.scopevars),
				len(vars),
				fnName,
			))
		}
	}

	nvm := vm.NewScope()
	for i, val := range function.scopevars {
		switch vars[i].Type {
		case TYPE_IFunc:
			nvm.Funcs[val] = vars[i].Value.(InBuilt)
		case TYPE_DFunc:
			nvm.LFuncs[val] = vars[i].Value.(Defined)
		default:
			if isVariate(val) {
				if len(function.scopevars) != i+1 {
					panic(fmt.Sprintf("In function %s, the variate parameter should be at the end", fnName))
				}
				val = val[3:]
				if len(vars) >= i+1 {
					nvm.Vars[val] = Variable{TYPE_MonoTypeArray, vars[i:]}
				}
				break
			}
			nvm.Vars[val] = vars[i]
		}
	}
	return nvm.Eval(function.eval)
}

// run is the method used to call the functions (defined or in-built) with the arguments
func (vm *VM) run(tkns []string) (Variable, error) {
	vars := make([]Variable, 0)
	fnName := tkns[0]
	for i := 1; i < len(tkns); i++ {

		if len(tkns[i]) > 3 && tkns[i][:3] == "..." && tkns[i][3] != '.' {
			v, err := vm.GetVariable(tkns[i][3:])
			if err != nil {
				return ligoNil, err
			}
			if v.Type == TYPE_MonoTypeArray || v.Type == TYPE_PolyTypeArray {
				vars = append(vars, v.Value.([]Variable)...)
				continue
			}
			vars = append(vars, v)
			continue
		}

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

// runLoop method is used to run the "loop" construct
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

// runIn method is used to run the "in" construct
func (vm *VM) runIn(tkns []string) (Variable, error) {
	if len(tkns) != 4 {
		return ligoNil, LigoError("Illegal in loop construct. Can take 4 arguments only.")
	}
	iterVar := tkns[2]
	runExp := tkns[3]

	array, err := vm.GetVariable(tkns[1])
	if err != nil {
		return ligoNil, err
	}

	if array.Type != TYPE_String && array.Type != TYPE_PolyTypeArray && array.Type != TYPE_MonoTypeArray {
		return ligoNil, LigoError("in : can only iterate thorugh arrays or strings")
	}

	v, ok := vm.Vars[iterVar]
	if array.Type == TYPE_String {
		for _, val := range array.Value.(string) {
			vm.Vars[iterVar] = Variable{Type: TYPE_String, Value: string(val)}
			_, err = vm.Eval(runExp)
			if err != nil {
				return ligoNil, err
			}
		}
	} else {
		for _, val := range array.Value.([]Variable) {
			vm.Vars[iterVar] = val
			_, err = vm.Eval(runExp)
			if err != nil {
				return ligoNil, err
			}
		}
	}
	if ok {
		vm.Vars[iterVar] = v
	} else {
		delete(vm.Vars, iterVar)
	}
	return ligoNil, nil
}

// ifClause is used to evaluate the "if" / "if...else" clause
// The if or else clause can be another subexp or can be just a variable.
// This variable is returned and can be passed directly to functions.
// See the samples/basic.lg file for more details.
func (vm *VM) ifClause(tkns []string) (Variable, error) {
	if len(tkns) > 4 || len(tkns) < 3 {
		return ligoNil, LigoError("Illegal if construct. Can take 3 or 4 arguments.")
	}
	condition := tkns[1]
	boolVar, ok := vm.Vars[condition]
	if condition != "true" && condition != "false" && MatchChars(condition, 0, '(', ')') < 0 && !ok {
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

// returnArg method is used to return a variable or a value.
func (vm *VM) returnArg(tkns []string) (Variable, error) {
	if len(tkns) != 2 {
		panic("Cannot return more than 2 values. (Atleast for now.)")
	}
	return vm.GetVariable(tkns[1])
}

// deleteVar method is used to delete a variable from the VM
func (vm *VM) deleteVar(tkns []string) (Variable, error) {
	if len(tkns) < 2 {
		return Variable{Type: TYPE_Bool, Value: false}, LigoError("nothing passed to delete")
	}
	for _, variable := range tkns[1:] {
		_, ok := vm.Vars[variable]
		if !ok {
			return Variable{Type: TYPE_Bool, Value: false}, LigoError("variable not found")
		}
		delete(vm.Vars, variable)
	}
	return Variable{Type: TYPE_Bool, Value: true}, nil
}

// fork method is used to run the passed sub-expression in a separate go-routine
func (vm *VM) fork(tkns []string) (Variable, error) {
	if len(tkns) != 2 {
		return ligoNil, LigoError("Expected one expression, got " + fmt.Sprint(len(tkns)) + " arguments")
	}
	go vm.Eval(tkns[1])
	return ligoNil, nil
}

// runExpressions method is used to run the passed sub-expressions
// Generally this is used inside a loop, function or condition clauses
// as then can only take one sub-expression for execution.
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

// evalString method is used to evaluate a passed string as a ligo expression and
// pass back it's return
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

// Eval method is used to parse a passed string and evaluate it.
// This is the entry point for any proper execution.
func (vm *VM) Eval(stmt string) (Variable, error) {
	stmt = strings.TrimSpace(stmt)
	if len(stmt) < 2 {
		return ligoNil, LigoError("Expected atleast (), got : " + stmt)
	}
	if !rExpression.MatchString(stmt) && MatchChars(stmt, 0, '(', ')') < 0 {
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
	case "in":
		return vm.runIn(tkns)
	case "if":
		return vm.ifClause(tkns)
	case "eval":
		return vm.evalString(tkns)
	case "fork":
		return vm.fork(tkns)
	case "delete":
		return vm.deleteVar(tkns)
	default:
		return vm.run(tkns)
	}
	return ligoNil, nil
}

// Clone method is used to clone the VM and return the clone one.
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

// NewScope method is used to create a new vm with global scope set from the
// current VM.
// (if the current vm is the parent vm, then it is set as the global, else the global of the current vm is set )
func (vm *VM) NewScope() *VM {
	nvm := NewVM()
	if vm.global == nil {
		nvm.global = vm
	} else {
		nvm.global = vm.global
	}
	return nvm
}

// LoadReader method is used to load script from a io.Reader and evaluate it
func (vm *VM) LoadReader(input io.Reader) error {
	ltxtb, err := ioutil.ReadAll(input)
	if err != nil {
		return err
	}
	ltxt := StripComments(string(ltxtb))
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
			off := MatchChars(ltxt, int64(i), '(', ')') + 1
			if off < int64(i) {
				return fmt.Errorf("Syntax error near %d:%d : %s", i, line, ltxt[i:])
			}
			exps = append(exps, ltxt[i:off])
			i = int(off) - 1
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
