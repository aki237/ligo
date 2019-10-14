package ligo

import (
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

// escape sequences to be replaced with the counterpart in a string
var escapeSequences = map[string]int{
	"\\\\": 0x5C,
	"\\a":  0x07,
	"\\b":  0x08,
	"\\e":  0x1B,
	"\\f":  0x0C,
	"\\n":  0x0A,
	"\\r":  0x0D,
	"\\t":  0x09,
	"\\v":  0x0B,
	"\\'":  0x27,
	"\"":   0x22,
	"\\ ":  0x20,
	"\\!":  0x21,
	"\\#":  0x23,
	"\\$":  0x24,
	"\\%":  0x25,
	"\\&":  0x26,
	"\\*":  0x2A,
	"\\<":  0x3C,
	"\\=":  0x3D,
	"\\>":  0x3E,
	"\\?":  0x3F,
	"\\|":  0x7C,
	"\\~":  0x7E,
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
var rInteger = regexp.MustCompile(`^[+-]?[0-9]+$`)
var rFloat = regexp.MustCompile(`^[+-]?[0-9]*\.[0-9]+$`)
var rString = regexp.MustCompile(`^\".*\"$`)
var rVariable = regexp.MustCompile(`^[[:alpha:]]+[[:alnum:]]*$`)
var rExpression = regexp.MustCompile(`^\\(.*\\)$`)
var rClosure = regexp.MustCompile(`^\\|([[:alpha:]]+[[:alnum:]]*\\s*)*\\.\\.\\.([[:alpha:]]+[[:alnum:]]*\\s*){0,1}|$`)
var rArray = regexp.MustCompile(`^\\[(.*\\s*)*\\]$`)

// Variable is a struct denoting a value in the VM
type Variable struct {
	Type  Type
	Value interface{}
}

// GetTypeString method returns a string value corresponding to the type of it's value
func (v Variable) GetTypeString() (tp string) {
	tp = ""
	switch v.Type {
	case TypeInt:
		tp = "int"
	case TypeFloat:
		tp = "float"
	case TypeBool:
		tp = "bool"
	case TypeString:
		tp = "string"
	case TypeNil:
		tp = "nil"
	case TypeArray:
		tp = "array"
	case TypeMap:
		tp = "map"
	case TypeIFunc:
		tp = "inbuilt function"
	case TypeDFunc:
		tp = "defined function"
	}
	return
}

// String method implements the Stringer interface for the Variable type
func (v Variable) String() string {
	typeString := "Variable {Type : <" + v.GetTypeString()
	return typeString + fmt.Sprint("> ,Value : ", v.Value, "}")
}

// Defined struct contains variables needed for storing a function defined in ligo script itself
type Defined struct {
	scopevars []string
	eval      string
}

// InBuilt type is a function format that is callable from the ligo script
type InBuilt func(*VM, ...Variable) Variable

// Map type is a ligo equivalent for dictionaly or hash maps
type Map map[Variable]Variable

// ProcessCommon is a struct type for process control and signal dispatch
type ProcessCommon struct {
	interrupt bool
	*sync.Mutex
}

// VM struct is a State Struct contains all the variable maps,
// defined function maps, in-built function maps and a global
// scope pointing to the global Scope VM
type VM struct {
	global         *VM
	exception      string
	Vars           map[string]Variable
	Funcs          map[string]InBuilt
	LFuncs         map[string]Defined
	namespaces     map[string]*VM
	pc             *ProcessCommon
	keywordHandler map[string]func([]string) (Variable, error)
	isNamespace    bool
}

// NewVM returns a new VM object pointer after initializing the values
func NewVM() *VM {
	vm := &VM{}
	vm.Vars = make(map[string]Variable)
	vm.Funcs = make(map[string]InBuilt)
	vm.LFuncs = make(map[string]Defined)
	vm.global = nil
	vm.pc = &ProcessCommon{Mutex: &sync.Mutex{}, interrupt: false}
	vm.keywordHandler = map[string]func([]string) (Variable, error){
		"var":       vm.newVar,
		"set":       vm.setVar,
		"fn":        vm.setFn,
		"return":    vm.returnArg,
		"progn":     vm.runExpressions,
		"loop":      vm.runLoop,
		"in":        vm.runIn,
		"if":        vm.ifClause,
		"match":     vm.matchClause,
		"eval":      vm.evalString,
		"fork":      vm.fork,
		"delete":    vm.deleteVar,
		"namespace": vm.namespaceEval,
		"lambda":    vm.lambdaEval,
		"struct":    vm.structEval,
	}
	vm.namespaces = make(map[string]*VM)
	vm.isNamespace = false
	return vm
}

// Stop method is used to stop the current process and return an error value.
func (vm *VM) Stop() {
	vm.pc.Lock()
	vm.pc.interrupt = true
}

// Resume method is used to resume the normal evaluation by releasing the lock
// on the mutex of the process control. This should never be called in this package
// itself. Resume should be used only when a error returned is ErrSignalRecieved
// in the main package. See the sample interpreter implementation in
// https://github.com/aki237/ligo/tree/master/cmd/ligo.
func (vm *VM) Resume() {
	vm.pc.Unlock()
	vm.pc.interrupt = false
}

// parseToArray is used to parse the given string into ligo.TypeArray
func (vm *VM) parseToArray(token string) (Variable, error) {
	ar := token[1:MatchChars(token, 0, '[', ']')]
	tkns, err := ScanTokens("(" + ar + ")")
	if err != nil {
		return ligoNil, err
	}
	vars := make([]Variable, 0)
	for _, val := range tkns {
		v, err := vm.GetVariable(val)
		if err != nil {
			return ligoNil, err
		}
		vars = append(vars, v)
	}
	retVars := Variable{Type: TypeArray, Value: vars}
	return retVars, nil
}

// parseToInt method is used to parse a given string to ligo.TypeInt
func (vm *VM) parseToInt(token string) (Variable, error) {
	num, err := strconv.ParseInt(token, 10, 64)
	if err != nil {
		return ligoNil, err
	}
	return Variable{Type: TypeInt, Value: num}, nil
}

// parseToFloat method is used to parse a given string to ligo.TypeFloat
func (vm *VM) parseToFloat(token string) (Variable, error) {
	num, err := strconv.ParseFloat(token, 64)
	if err != nil {
		return ligoNil, err
	}
	return Variable{Type: TypeFloat, Value: num}, nil
}

// parseToString method is used to reform escape sequences and
// return a ligo.TypeString
func (vm *VM) parseToString(token string) (Variable, error) {
	token = reformEscapes(token)
	return Variable{Type: TypeString, Value: token[1 : len(token)-1]}, nil
}

func getStructVar(strct Variable, key string) (Variable, error) {
	keys, ok := strct.Value.(map[string]Variable)
	if strct.Type != TypeStruct || !ok {
		return ligoNil, Error("passed variable is not a struct and doesn't have a member named '" + key + "'")
	}
	varName := key
	if strings.Contains(key, ":") {
		varName = strings.Split(key, ":")[0]
		v, ok := keys[varName]
		if !ok {
			return ligoNil, Error("no such key found in the struct : \"" + key + "\"")
		}
		return getStructVar(v, strings.Join(strings.Split(key, ":")[1:], ":"))
	}
	v, ok := keys[varName]
	if !ok {
		return ligoNil, Error("no such key found in the struct : \"" + key + "\"")
	}
	return v, nil
}

// parseToSymbol method is used to fetch the variable from the VM.
// If the variable is not found, it checks whether a function is found
// and returns it.
func (vm *VM) parseToSymbol(token string) (Variable, error) {
	varFromVM, ok := vm.Vars[token]
	if ok {
		return Variable{Type: varFromVM.Type, Value: varFromVM.Value}, nil
	}

	if strings.Contains(token, ":") {
		varName := strings.Split(token, ":")[0]
		v, err := vm.parseToSymbol(varName)
		if err != nil {
			return v, err
		}
		return getStructVar(v, strings.Join(strings.Split(token, ":")[1:], ":"))
	}

	nss := strings.Split(token, ".")

	if len(strings.Split(token, ".")) > 1 {
		ns := nss[0]
		fn := strings.Join(nss[1:], ".")
		namespace, ok := vm.namespaces[ns]
		if ok {
			return namespace.parseToSymbol(fn)
		}
	}

	function, err := vm.parseToFunc(token)
	if err == nil {
		return function, nil
	}

	if vm.global == nil {
		return ligoNil, ErrNoVariable + Error(" : "+token)
	}

	return vm.global.parseToSymbol(token)
}

// parseToFunc method is used to find the function from the vm and return it.
func (vm *VM) parseToFunc(token string) (Variable, error) {
	if fnc, ok := vm.Funcs[token]; ok {
		return Variable{Type: TypeIFunc, Value: fnc}, nil
	}
	if fnc, ok := vm.LFuncs[token]; ok {
		return Variable{Type: TypeDFunc, Value: fnc}, nil
	}

	nss := strings.Split(token, ".")

	if len(strings.Split(token, ".")) > 1 {
		ns := nss[0]
		fn := strings.Join(nss[1:], ".")
		namespace, ok := vm.namespaces[ns]
		if ok {
			return namespace.parseToFunc(fn)
		}
	}

	return ligoNil, ErrFuncNotFound + Error(" : "+token)
}

// GetVariable method is used to process the token string passed and get the
// corresponding value from the VM's memory. This is a crucial function
// as, if the token passed is a sub expression this method knows to evaluate and
// return the value of that sub expression.
func (vm *VM) GetVariable(token string) (Variable, error) {
	if len(token) < 1 {
		return ligoNil, Error("invalid Token passed")
	}
	switch true {
	case MatchChars(token, 0, '[', ']') > 0:
		return vm.parseToArray(token)
	case MatchChars(token, 0, '(', ')') > 0:
		return vm.Eval(token)
	case rInteger.MatchString(token):
		return vm.parseToInt(token)
	case rFloat.MatchString(token):
		return vm.parseToFloat(token)
	case rString.MatchString(token) || token[0] == '"':
		return vm.parseToString(token)
	case token == "true":
		return Variable{Type: TypeBool, Value: true}, nil
	case token == "false":
		return Variable{Type: TypeBool, Value: false}, nil
	}
	return vm.parseToSymbol(token)
}

// setFn is used to parse a ligo function construct and store it in the
// current scope. It also warns if the function is already declared.
func (vm *VM) setFn(tokens []string) (Variable, error) {
	if len(tokens) != 4 {
		return ligoNil, Error("A function construct can only have a single returning function")
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
			Error("Expected parameter name in the function definition " + fnName + " closure : " + tokens[2])
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
		return ligoNil, Error("Wrong number of arguments to the keyword.")
	}
	if !rVariable.MatchString(tokens[1]) {
		return ligoNil, Error("Wrong token found in the variable name")
	}
	v, err := vm.GetVariable(tokens[2])
	if err != nil {
		return ligoNil, err
	}

	switch v.Type {
	case TypeIFunc:
		_, ok := vm.Funcs[tokens[1]]
		if !ok {
			return ligoNil, Error("Variable not defined. Try \"var\" for creating a new variable")
		}
		vm.Funcs[tokens[1]] = v.Value.(InBuilt)
		return ligoNil, nil
	case TypeDFunc:
		_, ok := vm.LFuncs[tokens[1]]
		if !ok {
			return ligoNil, Error("Variable not defined. Try \"var\" for creating a new variable")
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
		return ligoNil, Error("Variable '" + tokens[1] + "' not defined. Try \"var\" for creating a new variable")
	}
	return vm.global.setVar(tokens)
}

// newVar method is used to declare a new variable in the VM and set a value to it.
func (vm *VM) newVar(tokens []string) (Variable, error) {
	if len(tokens) != 3 {
		return ligoNil, Error("Wrong number of arguments to the keyword.")
	}
	if !rVariable.MatchString(tokens[1]) {
		return ligoNil, Error("Wrong token found in the variable name")
	}
	v, err := vm.GetVariable(tokens[2])
	if err != nil {
		return ligoNil, err
	}
	switch v.Type {
	case TypeIFunc:
		vm.Funcs[tokens[1]] = v.Value.(InBuilt)
		return ligoNil, nil
	case TypeDFunc:
		vm.LFuncs[tokens[1]] = v.Value.(Defined)
		return ligoNil, nil
	}
	_, ok := vm.Vars[tokens[1]]
	if ok {
		return ligoNil, Error("Variable '" + tokens[1] + "' already defined. Try \"set\" for updating variables")
	}
	vm.Vars[tokens[1]] = v
	return ligoNil, nil
}

// getInBuiltFunction method is a small helper method to get the inbuilt function.
// If found it returns the inbuilt and true, else nil and false.
func (vm *VM) getInBuiltFunction(fnName string) (InBuilt, bool) {
	fn, found := vm.Funcs[fnName]
	if found {
		return fn, found
	}
	v, err := vm.parseToSymbol(fnName)
	if err != nil {
		return nil, false
	}
	if v.Type != TypeIFunc {
		return nil, false
	}
	return v.Value.(InBuilt), true
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
	if found {
		return fn, found
	}
	v, err := vm.parseToSymbol(fnName)
	if err != nil {
		return Defined{}, false
	}
	if v.Type != TypeDFunc {
		return Defined{}, false
	}
	return v.Value.(Defined), true
}

// RunDefined method is an outlet of the runDefinedFunction function
func (vm *VM) RunDefined(function Defined, vars []Variable) (Variable, error) {
	return vm.runDefinedFunction(function, "<defined function call>", vars)
}

// runDefinedFunction method is a helper method used to run a passed defined function with passed vars
func (vm *VM) runDefinedFunction(function Defined, fnName string, vars []Variable) (Variable, error) {
	if len(vars) < len(function.scopevars)-1 {
		return ligoNil, Error(fmt.Sprintf("Expected %d arguments, got %d for the %s function",
			len(function.scopevars),
			len(vars),
			fnName,
		))
	}

	if len(function.scopevars) > 0 && !isVariate(function.scopevars[len(function.scopevars)-1]) {
		if len(vars) != len(function.scopevars) {
			return ligoNil, Error(fmt.Sprintf("Expected %d arguments, got %d for the %s function",
				len(function.scopevars),
				len(vars),
				fnName,
			))
		}
	}

	nvm := vm.NewScope()
	for i, val := range function.scopevars {
		if len(vars)-1 < i {
			if isVariate(val) {
				val = val[3:]
				nvm.Vars[val] = Variable{Type: TypeArray, Value: make([]Variable, 0)}
				break
			}
			return ligoNil, Error("Not enough arguments to call the function")
		}
		switch vars[i].Type {
		case TypeIFunc:
			nvm.Funcs[val] = vars[i].Value.(InBuilt)
		case TypeDFunc:
			nvm.LFuncs[val] = vars[i].Value.(Defined)
		default:
			if isVariate(val) {
				if len(function.scopevars) != i+1 {
					panic(fmt.Sprintf("In function %s, the variate parameter should be at the end", fnName))
				}
				val = val[3:]
				if len(vars) >= i+1 {
					nvm.Vars[val] = Variable{TypeArray, vars[i:]}
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
	fnName := tkns[0]
	nspaces := strings.Split(fnName, ".")
	if len(nspaces) >= 2 {
		ns, ok := vm.namespaces[nspaces[0]]
		if ok {
			ntkns := tkns
			ntkns[0] = strings.Join(nspaces[1:], ".")
			return ns.run(ntkns)
		}
	}

	vars := make([]Variable, 0)
	for i := 1; i < len(tkns); i++ {

		if len(tkns[i]) > 3 && tkns[i][:3] == "..." && tkns[i][3] != '.' {
			v, err := vm.GetVariable(tkns[i][3:])
			if err != nil {
				return ligoNil, err
			}
			if v.Type == TypeArray {
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
		return ligoNil, Error("Function '" + fnName + "' not found in any of the namespaces")
	}
	if function, ok := vm.global.getInBuiltFunction(fnName); ok {
		return vm.runInBuiltFunction(function, vars)
	}
	if function, ok := vm.global.getDefinedFunction(fnName); ok {
		return vm.runDefinedFunction(function, fnName, vars)
	}
	return ligoNil, Error("Function '" + fnName + "' not found")
}

// runLoop method is used to run the "loop" construct
func (vm *VM) runLoop(tkns []string) (Variable, error) {
	if len(tkns) != 3 {
		return ligoNil, Error("Illegal loop construct. Can take 3 arguments only.")
	}
	condition := tkns[1]
	runExp := tkns[2]
	result, err := vm.Eval(condition)
	if err != nil {
		return ligoNil, err
	}
	if result.Type != TypeBool {
		return ligoNil, Error("Expected boolean return from the expression : " + condition)
	}
	for result.Value.(bool) {
		if vm.pc.interrupt {
			return ligoNil, ErrSignalRecieved
		}
		_, err := vm.Eval(runExp)
		if err != nil {
			return ligoNil, err
		}
		result, err = vm.Eval(condition)
		if err != nil {
			return ligoNil, err
		}
		if result.Type != TypeBool {
			return ligoNil, Error("Expected boolean return from the expression : " + condition)
		}
	}
	return ligoNil, err
}

// runIn method is used to run the "in" construct
func (vm *VM) runIn(tkns []string) (Variable, error) {

	if len(tkns) != 4 {
		return ligoNil, Error("Illegal in loop construct. Can take 4 arguments only.")
	}
	iterVar := tkns[2]
	runExp := tkns[3]

	array, err := vm.GetVariable(tkns[1])
	if err != nil {
		return ligoNil, err
	}

	if array.Type != TypeString && array.Type != TypeArray {
		return ligoNil, Error("in : can only iterate thorugh arrays or strings")
	}

	v, ok := vm.Vars[iterVar]
	if array.Type == TypeString {
		for _, val := range array.Value.(string) {
			vm.Vars[iterVar] = Variable{Type: TypeString, Value: string(val)}
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

// structEval method is used to evaluate the struct construct
// and return the corresponding variable
func (vm *VM) structEval(tkns []string) (Variable, error) {
	if len(tkns) < 3 && len(tkns)%2 == 0 {
		return ligoNil, Error("illegal struct construct. Should take in atleast 3 arguments")
	}

	mapVar := make(map[string]Variable)

	for i := 0; i < (len(tkns) / 2); i++ {
		index := 1 + (2 * i)
		key := tkns[index]
		val, err := vm.GetVariable(tkns[index+1])
		if err != nil {
			return ligoNil, err
		}
		mapVar[key] = val
	}
	return Variable{Type: TypeStruct, Value: mapVar}, nil
}

// matchClause is used to evaluate the match case construct
func (vm *VM) matchClause(tkns []string) (Variable, error) {
	if len(tkns) < 4 && len(tkns)%2 != 0 {
		return ligoNil, Error("illegal match construct. Should take in atleast 4 arguments")
	}

	matchVariable, err := vm.GetVariable(tkns[1])
	if err != nil {
		return ligoNil, err
	}

	for i := 1; i <= (len(tkns)/2)-1; i++ {
		if tkns[2*i] == "_" {
			if (2 * i) != len(tkns)-2 {
				return ligoNil, Error("default case '_' should be placed at last")
			}
			return vm.Eval(tkns[(2*i)+1])
		}

		caseVariable, err := vm.GetVariable(tkns[2*i])
		if err != nil {
			return ligoNil, err
		}
		if caseVariable == matchVariable {
			return vm.Eval(tkns[(2*i)+1])
		}
	}
	return ligoNil, nil

}

// ifClause is used to evaluate the "if" / "if...else" clause
// The if or else clause can be another subexp or can be just a variable.
// This variable is returned and can be passed directly to functions.
// See the samples/basic.lg file for more details.
func (vm *VM) ifClause(tkns []string) (Variable, error) {
	if len(tkns) > 4 || len(tkns) < 3 {
		return ligoNil, Error("Illegal if construct. Can take 3 or 4 arguments.")
	}
	condition := tkns[1]
	boolVar, ok := vm.Vars[condition]
	if condition != "true" && condition != "false" && MatchChars(condition, 0, '(', ')') < 0 && !ok {
		return ligoNil,
			Error("Expected a boolean value or expression for the if clause condition, got : " + condition)
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
	if result.Type != TypeBool {
		return ligoNil, Error("Expected boolean return from the expression : " + condition)
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
		return Variable{Type: TypeBool, Value: false}, Error("nothing passed to delete")
	}
	for _, variable := range tkns[1:] {
		_, ok := vm.Vars[variable]
		if !ok {
			return Variable{Type: TypeBool, Value: false}, Error("variable not found")
		}
		delete(vm.Vars, variable)
	}
	return Variable{Type: TypeBool, Value: true}, nil
}

// fork method is used to run the passed sub-expression in a separate go-routine
func (vm *VM) fork(tkns []string) (Variable, error) {
	if len(tkns) != 2 {
		return ligoNil, Error("Expected one expression, got " + fmt.Sprint(len(tkns)) + " arguments")
	}
	go vm.Eval(tkns[1])
	return ligoNil, nil
}

// namespaceEval method is used to run the code in a namespace environment
func (vm *VM) namespaceEval(tkns []string) (Variable, error) {
	if len(tkns) < 3 {
		return ligoNil, Error("Expected atleast 3 expressions, got " + fmt.Sprint(len(tkns)))
	}

	ns := tkns[1]
	splitted := strings.Split(ns, ".")
	nss, ok := vm.namespaces[splitted[0]]
	if !ok {
		nss = vm.CreateNamespace(splitted[0])
	}
	if len(splitted) < 2 {
		for i, val := range tkns[2:] {
			v, err := nss.Eval(val)
			if err != nil {
				return ligoNil, err
			}
			if i == len(tkns[2:])-1 {
				return v, nil
			}
		}
	}
	newTkns := make([]string, 0)
	newTkns = append(newTkns, []string{tkns[0], strings.Join(splitted[2:], ".")}...)
	for _, val := range newTkns[2:] {
		newTkns = append(newTkns, val)
	}

	return nss.namespaceEval(newTkns)
}

// lambdaEval is used to evaluate a lambda expression and return a
// ligo function
func (vm *VM) lambdaEval(tkns []string) (Variable, error) {
	if len(tkns) != 3 {
		return ligoNil, Error("Error in the lambda construct")
	}

	closure := tkns[1]
	body := tkns[2]

	if !rClosure.MatchString(tkns[2]) {
		return ligoNil, Error("malformed closure in the lambda, near " + closure)
	}
	varNames := getVarsFromClosure(closure)
	fn := Variable{Type: TypeDFunc, Value: Defined{scopevars: varNames, eval: body}}
	return fn, nil
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
		return ligoNil, Error("'eval' keyword only accepts 1 argument")
	}
	vl, err := vm.GetVariable(tkns[1])
	if err != nil {
		return ligoNil, err
	}
	if vl.Type != TypeString {
		return ligoNil, Error("'eval' keyword only expression string")
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
				line++
			}
			continue
		default:
			return ligoNil, Error(fmt.Sprintf("Unexpected Character at line %d : %s\n", line, ch))
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

// catchException method is used to manage the uncaught exception
func (vm *VM) catchException(tkns []string) (Variable, error) {
	if len(tkns) != 3 {
		return ligoNil, Error("catch expects 3 arguments")
	}

	if vm.exception == "" {
		return ligoNil, nil
	}

	scope := vm.Clone()

	scope.Vars[tkns[1]] = Variable{Value: vm.exception, Type: TypeString}
	vm.exception = ""
	return scope.Eval(tkns[2])
}

// Throw method is used to throw an exception in the VM.
func (vm *VM) Throw(exception string) Variable {
	vm.exception = exception
	return ligoNil
}

// Eval method is used to parse a passed string and evaluate it.
// This is the entry point for any proper execution.
func (vm *VM) Eval(stmt string) (Variable, error) {
	if vm.pc.interrupt {
		return ligoNil, ErrSignalRecieved
	}
	stmt = strings.TrimSpace(stmt)
	if len(stmt) < 1 {
		return ligoNil, Error("Expected atleast a token, got : " + stmt)
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

	if vm.exception != "" && fnName != "catch" {
		return ligoNil, ErrExceptionNotHandled
	}

	if fnName == "catch" {
		return vm.catchException(tkns)
	}
	return vm.evalKeyword(fnName, tkns)
}

// evalKeyword is used to run the corresponding function for the given keyword
func (vm *VM) evalKeyword(fnName string, tkns []string) (Variable, error) {
	handler, ok := vm.keywordHandler[fnName]
	if ok {
		return handler(tkns)
	}
	return vm.run(tkns)
}

// GetNameSpace method is used to get the namespace scope corresponding to the name passed
func (vm *VM) GetNameSpace(ns string) *VM {
	namespace, ok := vm.namespaces[ns]
	if !ok {
		return nil
	}
	return namespace
}

// CreateNamespace method is used to create a new namespace if it doesn't exist
func (vm *VM) CreateNamespace(ns string) *VM {
	if namespace := vm.GetNameSpace(ns); namespace != nil {
		return namespace
	}
	vm.namespaces[ns] = vm.NewScope()
	vm.namespaces[ns].isNamespace = true
	return vm.namespaces[ns]
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
	if vm.global == nil || vm.isNamespace {
		nvm.global = vm
	} else {
		nvm.global = vm.global
	}
	nvm.pc = vm.pc
	return nvm
}

// LoadReader method is used to load script from a io.Reader and evaluate it
func (vm *VM) LoadReader(input io.Reader) error {
	ltxtb, err := ioutil.ReadAll(input)
	if err != nil {
		return err
	}

	exps, err := vm.BreakChunk(string(ltxtb))
	if err != nil {
		return err
	}

	for _, val := range exps {
		_, err := vm.Eval(val)
		if err != nil {
			return fmt.Errorf("error : %s", err)
		}
	}

	return nil
}

// BreakChunk is used to break a chunk of ligo code into string list of subexps
func (vm VM) BreakChunk(ltxt string) ([]string, error) {
	ltxt = StripComments(ltxt)
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
				return nil, fmt.Errorf("Syntax error near %d:%d : %s", i, line, ltxt[i:])
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
			return nil, fmt.Errorf("unexpected Character at line %d : %s", line, ch)
		}
	}
	return exps, nil
}
