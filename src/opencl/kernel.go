package opencl

/*
#cgo LDFLAGS: -framework OpenCL
#include <OpenCL/opencl.h>
*/
import "C"

import (
	"unsafe"
)

type Kernel C.cl_kernel

func CreateKernel(program Program, kernelName string) (Kernel, error) {
	var clErr C.cl_int
	clKernelName := C.CString(kernelName)
	defer C.free(unsafe.Pointer(clKernelName))
	clKernel := C.clCreateKernel(program, clKernelName, &clErr)
	return Kernel(clKernel), convertError(clErr)
}

func SetKernelArg(kernel Kernel, position int, data Memory) error {
	var t C.cl_mem
	clErr := C.clSetKernelArg(kernel, C.cl_uint(position), C.size_t(unsafe.Sizeof(t)), unsafe.Pointer(&data))
	return convertError(clErr)
}

func EnqueueNDRangeKernel(commandQueue CommandQueue, kernel Kernel, workSize []int) error {
	globalWorkSize := make([]C.size_t, len(workSize))
	for index, size := range workSize {
		globalWorkSize[index] = C.size_t(size)
	}
	clErr := C.clEnqueueNDRangeKernel(C.cl_command_queue(commandQueue), C.cl_kernel(kernel), C.cl_uint(len(globalWorkSize)), nil, &globalWorkSize[0], nil, 0, nil, nil)
	return convertError(clErr)
}

func ReleaseKernel(kernel Kernel) error {
	clErr := C.clReleaseKernel(C.cl_kernel(kernel))
	return convertError(clErr)
}
