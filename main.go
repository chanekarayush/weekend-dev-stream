package main

import (
	"bytes"
	"flag"
	"io"
	"log"
	"os"
)

const (
	Wr   = 0.299
	Wg   = 0.587
	Wb   = 0.114
	Umax = 0.436
	Vmax = 0.615
)

func main() {
	var height int
	var width int
	flag.IntVar(&height, "ht", 216, "Height of the Video Stream")
	flag.IntVar(&width, "wd", 384, "Width of the Video Stream")

	frames := readFrames(width, height)
	convertedFrames := convertRGBtoYUV(width, height, frames)
	log.Println(len(convertedFrames))
	writeFileToDisk(convertedFrames, "encodedVideo.yuv")
}

func readFrames(width int, height int) [][]byte {
	frames := make([][]byte, 0)
	for {
		// RGB video has R, G, B channels, allocate space for them
		frame := make([]byte, width*height*3)
		_, err := io.ReadFull(os.Stdin, frame)
		if err != nil {
			break
		}
		frames = append(frames, frame)
	}
	// TODO: show size of frames
	return frames
}

func convertRGBtoYUV(width int, height int, frames [][]byte) [][]byte {
	// convertDSFrames := make([][]byte, len(frames))
	// copy(convertDSFrames, frames)
	for i, frame := range frames {
		// RGB to YUV conversion
		Y := make([]byte, width*height)
		U := make([]float64, width*height)
		V := make([]float64, width*height)
		for j := 0; j < width*height; j++ {
			r := float64(frame[3*j])
			g := float64(frame[3*j+1])
			b := float64(frame[3*j+2])

			y := Wr*r + Wg*g + Wb*b
			u := Umax * ((b - y) / (1 - Wb))
			v := Vmax * ((r - y) / (1 - Wr))

			Y[j] = uint8(y)
			U[j] = u
			V[j] = v
		}
		// Chroma Sampling ratio 4:2:0
		uDSample := make([]byte, width*height/4)
		vDSample := make([]byte, width*height/4)

		// Downsampling
		// ||=======||=======||
		// ||	Y1	||	 Y2	 ||
		// ||	U1	||	 U2	 ||
		// ||	V1	||	 V2	 ||
		// ||=======||=======||
		// ||	Y3	||	 Y4	 ||
		// ||	U3	||	 U4	 ||
		// ||	V3	||	 V4	 ||
		// ||=======||=======||
		//
		// Keep the Y as is and average out U and V so that
		// Ui = avg(U1, U2, U3, U4)
		// Vi = avg(V1, V2, V3, V4)
		// Basically a grid with a average pooling for U and V

		for r := 0; r < height; r += 2 {
			for c := 0; c < width; c += 2 {
				ui := (U[r*width+c] + U[r*width+c+1] + U[(r+1)*width+c] + U[(r+1)*width+c+1]) / 4
				vi := (V[r*width+c] + V[r*width+c+2] + V[(r+1)*width+c] + V[(r+1)*width+c+1]) / 4
				uDSample[r/2+width/2+c/2] = uint8(ui)
				vDSample[r/2+width/2+c/2] = uint8(vi)

			}
		}

		yuv420Frame := make([]byte, len(Y)+len(uDSample)+len(vDSample))

		copy(yuv420Frame, Y)
		copy(yuv420Frame[len(Y):], uDSample)
		copy(yuv420Frame[len(Y)+len(uDSample):], vDSample)

		frames[i] = yuv420Frame

	}
	// TODO: show size
	return frames
}

func writeFileToDisk(file [][]byte, outputName string) {
	err := os.WriteFile(outputName, bytes.Join(file, nil), 0o644)
	if err != nil {
		log.Print("WriteFileToDiskFailed")
		log.Fatalln(err)
	}
}
