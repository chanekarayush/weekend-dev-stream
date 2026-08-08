package main

import (
	"bytes"
	"compress/flate"
	"log"
)

func convertRGBtoYUV(width int, height int, frames [][]byte) [][]byte {
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

		// -ve -> green U and V
		// underflow -> -100 (green) to be come 246
		for r := 0; r < height; r += 2 {
			for c := 0; c < width; c += 2 {
				ui := (U[r*width+c] + U[r*width+c+1] + U[(r+1)*width+c] + U[(r+1)*width+c+1]) / 4
				vi := (V[r*width+c] + V[r*width+c+1] + V[(r+1)*width+c] + V[(r+1)*width+c+1]) / 4
				uDSample[r/2*width/2+c/2] = uint8(ui + 128)
				vDSample[r/2*width/2+c/2] = uint8(vi + 128)

			}
		}

		yuv420Frame := make([]byte, len(Y)+len(uDSample)+len(vDSample))

		copy(yuv420Frame, Y)
		copy(yuv420Frame[len(Y):], uDSample)
		copy(yuv420Frame[len(Y)+len(uDSample):], vDSample)

		frames[i] = yuv420Frame

	}
	log.Printf("[ConvertedRGBtoYUV] Size: %d", size(frames))
	return frames
}

func findKeyframeDeflate(frames [][]byte) []byte {

	var keyframeDeflate bytes.Buffer

	w, err := flate.NewWriter(&keyframeDeflate, flate.BestCompression)
	if err != nil {
		log.Fatalln(err)
	}
	for i := range frames {
		if i == 0 {
			keyframe := frames[i]
			_, err := w.Write(keyframe)
			if err != nil {
				log.Fatalln(err)
			}
		} else {
			difference := make([]byte, len(frames[i]))
			for j := range len(frames[i]) {
				// current frame - keyframe
				difference[j] = frames[i][j] - frames[i-1][j]
			}
			_, err := w.Write(difference)
			if err != nil {
				log.Fatalln(err)
			}
		}
	}
	err = w.Close()
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("[Deflated Data] Size: %d", keyframeDeflate.Len())
	return keyframeDeflate.Bytes()

}
