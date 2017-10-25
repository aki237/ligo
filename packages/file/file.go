package main

import (
	"fmt"
	"io"
	"os"

	"github.com/aki237/ligo/pkg/ligo"
)

// PluginInit function is the plugin initializer for the file package
func PluginInit(vm *ligo.VM) {
	vm.Funcs["file-open"] = vmFileOpen   // (file-open "filename.txt" "rw") => file handler   | panics
	vm.Funcs["file-read"] = vmFileRead   // (file-read fh nchars)           => string         | panics
	vm.Funcs["file-close"] = vmFileClose // (file-close fh)                 => nil            | error
	vm.Funcs["file-seek"] = vmFileSeek   // (file-seek fh amt from)         => current offset | panics
	vm.Funcs["file-write"] = vmFileWrite // (file-write fh string)          => written amount | panics
}

func vmFileSeek(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 3 {
		panic("file-seek : expected 3 arguments exactly")
	}
	if a[0].Type != 0x10080 {
		panic("file-seek : not a valid file handler")
	}

	if a[1].Type != ligo.TypeInt {
		panic("file-seek : position should be an integer")
	}

	if a[2].Type != ligo.TypeInt {
		panic("file-seek : whence should be an integer")
	}
	fh := a[0].Value.(*os.File)
	pos := a[1].Value.(int64)
	whence := a[2].Value.(int64)

	offset, err := fh.Seek(pos, int(whence))
	if err != nil {
		panic("Error while seeking : " + err.Error())
	}

	return ligo.Variable{Type: ligo.TypeInt, Value: offset}
}

func vmFileClose(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 1 {
		panic("file-close : expected 1 arguments exactly")
	}
	if a[0].Type != 0x10080 {
		panic("file-close : not a valid file handler")
	}

	fh := a[0].Value.(*os.File)
	err := fh.Close()
	if err != nil {
		return ligo.Variable{Type: ligo.TypeErr, Value: err}
	}
	return ligo.Variable{Type: ligo.TypeNil, Value: nil}

}

func vmFileOpen(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 2 {
		panic("file-open : expected 2 arguments exactly")
	}

	f := a[0]
	if f.Type != ligo.TypeString {
		panic("file-open : expects a filename as a string")
	}
	filename := f.Value.(string)
	if a[1].Type != ligo.TypeString {
		panic("file-open : expects a mode as a string")
	}
	mode := a[1].Value.(string)
	var fl *os.File
	var err error
	switch mode {
	case "r":
		fl, err = os.OpenFile(filename, os.O_CREATE|os.O_RDONLY, 0755)
		if err != nil {
			panic(fmt.Sprintf("file-open : error occurred : %s", err))
		}
	case "w":
		fl, err = os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0755)
		if err != nil {
			panic(fmt.Sprintf("file-open : error occurred : %s", err))
		}
	case "rw":
		fl, err = os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0755)
		if err != nil {
			panic(fmt.Sprintf("file-open : error occurred : %s", err))
		}
	case "a":
		fl, err = os.OpenFile(filename, os.O_CREATE|os.O_APPEND, 0755)
		if err != nil {
			panic(fmt.Sprintf("file-open : error occurred : %s", err))
		}
	default:
		panic("file-open : unrecogonized mode \"" + mode + "\"")
	}

	return ligo.Variable{Type: 0x10080, Value: fl}
}

func vmFileRead(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 2 {
		panic("file-read : expected 2 arguments exactly")
	}
	if a[0].Type != 0x10080 {
		panic("file-read : not a valid file handler")
	}

	fh := a[0].Value.(*os.File)

	if a[1].Type != ligo.TypeInt {
		panic("file-read : not a valid character count")
	}

	amt := a[1].Value.(int64)

	p := make([]byte, amt)

	read, err := fh.Read(p)
	if err != nil && err != io.EOF {
		panic(fmt.Sprintf("file-read : error occurred while reading : %s", err))
	}
	return ligo.Variable{Type: ligo.TypeString, Value: string(p[:read])}
}

func vmFileWrite(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if len(a) != 2 {
		panic("file-write : expected 2 arguments exactly")
	}
	if a[0].Type != 0x10080 {
		panic("file-write : not a valid file handler")
	}

	fh := a[0].Value.(*os.File)

	if a[1].Type != ligo.TypeString {
		panic("file-write : not a valid string")
	}

	str := a[1].Value.(string)

	written, err := fh.Write([]byte(str))
	if err != nil {
		panic(fmt.Sprintf("file-write : error occurred while writing : %s", err))
	}
	return ligo.Variable{Type: ligo.TypeInt, Value: int64(written)}
}
