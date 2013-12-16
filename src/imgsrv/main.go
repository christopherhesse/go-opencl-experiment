package main

/*
Prototype: OpenCL image processing

The native libraries for encoding/decoding are not very fast, IPP would
probably be a good bet.

If this needs to be fast, go could be okay, but IPP + OpenCL C library would probably be the way to go.
*/

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"log"
	cl "opencl"
	"os"
	"time"
)

const source string = `
kernel void copy(global const uchar4* in, global uchar4* out) {
    size_t i = get_global_id(0);
    out[i] = in[i];
}

constant float3 luminanceWeighting = float3(0.2125, 0.7154, 0.0721);
kernel void greyscale(global const uchar4* in, global uchar4* out) {
    size_t i = get_global_id(0);
    float4 pixel = convert_float4(in[i]) / 255;
    float luminance = dot(pixel.rgb, luminanceWeighting);
    pixel.rgb = float3(luminance);
    out[i] = convert_uchar4_sat_rte(pixel * 255);
}

kernel void invert(global const uchar4* in, global uchar4* out) {
    size_t i = get_global_id(0);
    float4 pixel = convert_float4(in[i]) / 255;
    pixel.rgb = float3(1.0) - pixel.rgb;
    out[i] = convert_uchar4_sat_rte(pixel * 255);
}
`

func convertToNRGBA(input image.Image) *image.NRGBA {
	if output, ok := input.(*image.NRGBA); ok {
		return output
	}
	// the image isn't the right kind, convert it
	bounds := input.Bounds()
	output := image.NewNRGBA(bounds)
	draw.Draw(output, bounds, input, bounds.Min, draw.Src)
	return output
}

func readImage(fileName string) image.Image {
	file, err := os.Open(fileName)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	return img
}

func main() {
	deviceIds, _ := cl.GetDeviceIDs(cl.DEVICE_TYPE_GPU)
	deviceId := deviceIds[0]
	context, err := cl.CreateContext(deviceId)
	queue, err := cl.CreateCommandQueue(context, deviceId)
	program, err := cl.CreateProgramWithSource(context, source)
	err = cl.BuildProgram(program)
	kernel, err := cl.CreateKernel(program, "invert")

	var start time.Time

	start = time.Now()
	img := readImage("rethinkdb.jpg")
	fmt.Println("decode:", time.Now().Sub(start))
	start = time.Now()
	input := convertToNRGBA(img)
	fmt.Println("convert:", time.Now().Sub(start))
	inputData := input.Pix

	start = time.Now()
	inputBuffer, err := cl.CreateBuffer(context, cl.MEM_READ_WRITE, len(inputData))
	outputBuffer, err := cl.CreateBuffer(context, cl.MEM_READ_WRITE, len(inputData))
	err = cl.EnqueueWriteBuffer(queue, inputBuffer, false, inputData)
	err = cl.SetKernelArg(kernel, 0, inputBuffer)
	err = cl.SetKernelArg(kernel, 1, outputBuffer)
	err = cl.EnqueueNDRangeKernel(queue, kernel, []int{len(inputData)})
	outputData := make([]uint8, len(inputData))
	err = cl.EnqueueReadBuffer(queue, outputBuffer, false, outputData)
	err = cl.Finish(queue)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("opencl:", time.Now().Sub(start))

	output := image.NRGBA{
		Pix:    outputData,
		Stride: input.Stride,
		Rect:   input.Rect,
	}
	outputFile, _ := os.Create("rethinkdb.png")
	defer outputFile.Close()
	start = time.Now()
	png.Encode(outputFile, &output)
	fmt.Println("encode:", time.Now().Sub(start))

	cl.ReleaseMemObject(inputBuffer)
	cl.ReleaseMemObject(outputBuffer)
	cl.ReleaseProgram(program)
	cl.ReleaseKernel(kernel)
	cl.ReleaseCommandQueue(queue)
	cl.ReleaseContext(context)
}
