package opencl

/*
#cgo LDFLAGS: -framework OpenCL
#include <OpenCL/opencl.h>
*/
import "C"

import (
	"unsafe"
)

type Context C.cl_context

func CreateContext(deviceId DeviceId) (Context, error) {
	clDeviceIds := make([]C.cl_device_id, 1)
	clDeviceIds[0] = C.cl_device_id(deviceId)

	var clErr C.cl_int
	pfnNotify := (*[0]byte)(unsafe.Pointer(C.clLogMessagesToStdoutAPPLE))
	context := C.clCreateContext(nil, C.cl_uint(len(clDeviceIds)), &clDeviceIds[0], pfnNotify, nil, &clErr)
	return Context(context), nil
}

func ReleaseContext(context Context) error {
	clErr := C.clReleaseContext(C.cl_context(context))
	return convertError(clErr)
}
