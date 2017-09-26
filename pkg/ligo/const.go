package ligo

// Error contants
const (
	ErrSyntaxError  Error = "Syntax Error"
	ErrNoVariable   Error = "Variable not found in scope"
	ErrFuncNotFound Error = "Function not defined in scope"
)

// Type is a type to denote the type of Variables in the VM
type Type int

// Required constants for the variable type
const (
	TypeErr             Type = -0x00
	TypeInt             Type = 0x000
	TypeFloat           Type = 0x001
	TypeBool            Type = 0x002
	TypeString          Type = 0x003
	TypeNil             Type = 0x004
	TypeIFunc           Type = 0x005
	TypeDFunc           Type = 0x006
	TypeExp             Type = 0x007
	TypeMonoTypeArray   Type = 0x100
	TypePolyTypeArray   Type = 0x200
	TypeMap             Type = 0x300
	TypeReader          Type = 0x400
	TypeWriter          Type = 0x500
	TypeSeeker          Type = 0x600
	TypeCloser          Type = 0x700
	TypeReadCloser      Type = 0x701
	TypeReadSeeker      Type = 0x702
	TypeReadWriter      Type = 0x703
	TypeReadWriteSeeker Type = 0x704
	TypeReadWriteCloser Type = 0x705
)

var ligoNil = Variable{TypeNil, nil}
