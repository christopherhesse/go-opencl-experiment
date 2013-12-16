package opencl

/*
#cgo LDFLAGS: -framework OpenCL
#include <OpenCL/opencl.h>
*/
import "C"

func convertBoolToClBool(value bool) C.cl_bool {
	if value {
		return C.CL_TRUE
	} else {
		return C.CL_FALSE
	}
}

func convertClBoolToBool(value C.cl_bool) bool {
	return value == C.CL_TRUE
}
