package opencl

/*
#cgo LDFLAGS: -framework OpenCL
#include <OpenCL/opencl.h>
*/
import "C"

func Finish(commandQueue CommandQueue) error {
	return convertError(C.clFinish(commandQueue))
}
