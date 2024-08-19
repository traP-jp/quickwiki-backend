package main

/*
#cgo pkg-config: python-3.12
#cgo LDFLAGS: -L. -lhello -lpython3.12
#include <Python.h>
#include <stdio.h>
#include <string.h>
#include "hello.h"
*/
import "C"
import "fmt"

func main() {
	name := C.CString("yattane")
	fileName := C.CString("mypy")
	funcName := C.CString("hello")
	result := C.pyHello(name, fileName, funcName)
	fmt.Printf("res= %s\n", C.GoString(result))
}
