package opencl

/*
#cgo LDFLAGS: -framework OpenCL
#include <OpenCL/opencl.h>
*/
import "C"

type CommandQueue C.cl_command_queue

func CreateCommandQueue(context Context, deviceId DeviceId) (CommandQueue, error) {
	var clErr C.cl_int
	clCommandQueue := C.clCreateCommandQueue(C.cl_context(context), C.cl_device_id(deviceId), C.CL_QUEUE_PROFILING_ENABLE, &clErr)
	if clErr != C.CL_SUCCESS {
		return nil, convertError(clErr)
	}
	return CommandQueue(clCommandQueue), nil
}

func ReleaseCommandQueue(commandQueue CommandQueue) error {
	clErr := C.clReleaseCommandQueue(C.cl_command_queue(commandQueue))
	return convertError(clErr)
}
