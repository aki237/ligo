package ligo

// Error contants
const (
	ErrSyntaxError  LigoError = "Syntax Error"
	ErrNoVariable   LigoError = "Variable not found in scope"
	ErrFuncNotFound LigoError = "Function not defined in scope"
)

// Type is a type to denote the type of Variables in the VM
type Type int

// Required constants for the variable type
const (
	TYPE_Err             Type = -0x00
	TYPE_Int             Type = 0x000
	TYPE_Float           Type = 0x001
	TYPE_Bool            Type = 0x002
	TYPE_String          Type = 0x003
	TYPE_Nil             Type = 0x004
	TYPE_IFunc           Type = 0x005
	TYPE_DFunc           Type = 0x006
	TYPE_Exp             Type = 0x007
	TYPE_MonoTypeArray   Type = 0x100
	TYPE_PolyTypeArray   Type = 0x200
	TYPE_Map             Type = 0x300
	TYPE_Reader          Type = 0x400
	TYPE_Writer          Type = 0x500
	TYPE_Seeker          Type = 0x600
	TYPE_Closer          Type = 0x700
	TYPE_ReadCloser      Type = 0x701
	TYPE_ReadSeeker      Type = 0x702
	TYPE_ReadWriter      Type = 0x703
	TYPE_ReadWriteSeeker Type = 0x704
	TYPE_ReadWriteCloser Type = 0x705
)

var ligoNil = Variable{TYPE_Nil, nil}
