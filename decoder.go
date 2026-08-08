package main

import (
	"bytes"
	"compress/flate"
	"io"
	"log"
)

func inflateKeyframeDelta(bytesToInflate []byte, width int, height int) [][]byte {
	var keyframeInflate bytes.Buffer
	keyframeInflate.Write(bytesToInflate)
	var inflated bytes.Buffer
	log.Printf("[Extracting] Size: %d", keyframeInflate.Len())
	// extract
	r := flate.NewReader(&keyframeInflate)
	// copy
	if _, err := io.Copy(&inflated, r); err != nil {
		log.Fatalln(err)
	}
	if err := r.Close(); err != nil {
		log.Fatalln(err)
	}
	log.Printf("[Extracted] Size: %d", inflated.Len())

	// outer loop ran for width * height * 3
	// inner loop ran for (width * height) / 2
	// resultant -> width * height * 3 / 2
	decodedFrames := make([][]byte, 0)
	for {
		frame := make([]byte, width*height*3/2)
		if _, err := io.ReadFull(&inflated, frame); err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalln(err)
		}
		decodedFrames = append(decodedFrames, frame)
	}

	log.Printf("[Decoded Data] Size: %d", size(decodedFrames))

	for i := range decodedFrames {
		if i == 0 {
			continue
		}

		for j := range len(decodedFrames[i]) {
			decodedFrames[i][j] += decodedFrames[i-1][j]
		}
	}
	return decodedFrames
}

func convertYUVtoRGB(width int, height int, decodedFrames [][]byte) [][]byte {
	for i, frame := range decodedFrames {
		// Subsampling ratio was 4:2:0
		Y := frame[:height*width]
		U := frame[height*width : height*width+(height*width/4)]
		V := frame[height*width+(height*width/4):]

		originalFrame := make([]byte, 0, height*width*3)
		for j := range height {
			for k := range width {
				y, u, v := float64(Y[j*width+k]), float64(U[(j/2)*(width/2)+(k/2)])-128, float64(V[(j/2)*(width/2)+(k/2)])-128

				r := y + v*(1-Wr)/Vmax
				g := y - u*(Wb*(1-Wb))/Umax*Wg - v*(Wr*(1-Wr))/Vmax*Wg
				b := y + u*(1-Wb)/Umax

				r = clamp(r, 0, 255)
				g = clamp(g, 0, 255)
				b = clamp(b, 0, 255)

				originalFrame = append(originalFrame, uint8(r), uint8(g), uint8(b))
			}
		}
		decodedFrames[i] = originalFrame
	}
	return decodedFrames
}
