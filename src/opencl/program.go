package opencl

/*
#cgo LDFLAGS: -framework OpenCL
#include <OpenCL/opencl.h>
*/
import "C"

import (
	"unsafe"
)

type Program C.cl_program

func CreateProgramWithSource(context Context, source string) (Program, error) {
	clSources := make([]*C.char, 1)
	clSources[0] = C.CString(source)
	defer C.free(unsafe.Pointer(clSources[0]))
	clLengths := make([]C.size_t, 1)
	clLengths[0] = C.size_t(len(source))
	var clErr C.cl_int
	clProgram := C.clCreateProgramWithSource(C.cl_context(context), C.cl_uint(len(clSources)), &clSources[0], &clLengths[0], &clErr)
	return Program(clProgram), convertError(clErr)
}

func BuildProgram(program Program) error {
	clErr := C.clBuildProgram(program, 0, nil, nil, nil, nil)
	return convertError(clErr)
}

func ReleaseProgram(program Program) error {
	clErr := C.clReleaseProgram(C.cl_program(program))
	return convertError(clErr)
}
