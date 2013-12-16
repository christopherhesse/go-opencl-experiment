package opencl

/*
#cgo LDFLAGS: -framework OpenCL
#include <OpenCL/opencl.h>
*/
import "C"

import (
	"unsafe"
)

type Memory C.cl_mem
type MemoryFlags C.cl_mem_flags

const (
	MEM_READ_WRITE MemoryFlags = C.CL_MEM_READ_WRITE
	MEM_WRITE_ONLY MemoryFlags = C.CL_MEM_WRITE_ONLY
	MEM_READ_ONLY  MemoryFlags = C.CL_MEM_READ_ONLY
)

func CreateBuffer(context Context, flags MemoryFlags, size int) (Memory, error) {
	var clErr C.cl_int
	clMemory := C.clCreateBuffer(context, C.CL_MEM_READ_WRITE, C.size_t(size), nil, &clErr)
	return Memory(clMemory), convertError(clErr)
}

func EnqueueWriteBuffer(commandQueue CommandQueue, buffer Memory, blocking bool, data []byte) error {
	clErr := C.clEnqueueWriteBuffer(C.cl_command_queue(commandQueue), C.cl_mem(buffer), convertBoolToClBool(blocking), C.size_t(0), C.size_t(len(data)), unsafe.Pointer(&data[0]), 0, nil, nil)
	return convertError(clErr)
}

// err = clEnqueueReadBuffer(queue, output, CL_FALSE, 0, bytes, output_pixels, 0, NULL, NULL)
func EnqueueReadBuffer(commandQueue CommandQueue, buffer Memory, blocking bool, data []byte) error {
	clErr := C.clEnqueueReadBuffer(C.cl_command_queue(commandQueue), C.cl_mem(buffer), convertBoolToClBool(blocking), C.size_t(0), C.size_t(len(data)), unsafe.Pointer(&data[0]), 0, nil, nil)
	return convertError(clErr)
}

func ReleaseMemObject(memory Memory) error {
	clErr := C.clReleaseMemObject(C.cl_mem(memory))
	return convertError(clErr)
}
