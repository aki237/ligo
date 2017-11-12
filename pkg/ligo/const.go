package ligo

// Error contants
const (
	ErrSyntaxError         Error = "Syntax Error"
	ErrNoVariable          Error = "Variable not found in scope"
	ErrFuncNotFound        Error = "Function not defined in scope"
	ErrSignalRecieved      Error = "Caught cancellation amidst evaluation"
	ErrExceptionNotHandled Error = "Exception not handled"
)

// Type is a type to denote the type of Variables in the VM
type Type int

// Required constants for the variable type
const (
	TypeErr    Type = -0x00
	TypeInt    Type = 0x000
	TypeFloat  Type = 0x001
	TypeBool   Type = 0x002
	TypeString Type = 0x003
	TypeNil    Type = 0x004
	TypeIFunc  Type = 0x005
	TypeDFunc  Type = 0x006
	TypeExp    Type = 0x007
	TypeArray  Type = 0x100
	TypeMap    Type = 0x300
	TypeStruct Type = 0x400
)

var ligoNil = Variable{TypeNil, nil}
