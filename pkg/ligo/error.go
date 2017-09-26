package ligo

// Error is a type string used to denote errors from the VM
type Error string

// Error method implements the error interface for the type Error
func (le Error) Error() string {
	return string(le)
}
