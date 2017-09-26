package ligo

// LigoError is a type string used to denote errors from the VM
type LigoError string

// Error method implements the error interface for the type LigoError
func (le LigoError) Error() string {
	return string(le)
}
