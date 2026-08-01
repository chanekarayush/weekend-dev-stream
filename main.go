package main

import (
	"bytes"
	"compress/flate"
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
	flag.IntVar(&width, "wt", 384, "Width of the Video Stream")
	flag.IntVar(&height, "ht", 216, "Height of the Video Stream")
	flag.Parse()

	frames := readFrames(width, height)
	convertedFrames := convertRGBtoYUV(width, height, frames)
	deflatedBytes := findKeyframeDeflate(convertedFrames)
	writeFileToDisk(convertedFrames, "encodedVideo.yuv")
	writeBytesToDisk(deflatedBytes, "encoded.yuv")
	bytesToInflate := readBytesFromDisk("encoded.yuv")
	decodedFrames := inflateKeyframeDelta(bytesToInflate, width, height)
	reconvertedFrames := convertYUVtoRGB(width, height, decodedFrames)
	writeFileToDisk(decodedFrames, "decoded.yuv")
	writeFileToDisk(reconvertedFrames, "original.rgb24")
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

func writeFileToDisk(file [][]byte, outputName string) {
	err := os.WriteFile(outputName, bytes.Join(file, nil), 0o644)
	if err != nil {
		log.Print("WriteFileToDiskFailed")
		log.Fatalln(err)
	}
}

func writeBytesToDisk(file []byte, outputName string) {
	err := os.WriteFile(outputName, file, 0o644)
	if err != nil {
		log.Print("WriteBytesToDiskFailed")
		log.Fatalln(err)
	}
}
func readBytesFromDisk(outputName string) []byte {
	buffer, err := os.ReadFile(outputName)
	if err != nil {
		log.Print("WriteBytesToDiskFailed")
		log.Fatalln(err)
	}
	return buffer
}

func clamp(val, min, max float64) float64 {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

func size(input [][]byte) int {
	i := 0
	for _, x := range input {
		i += len(x)
	}
	return i
}
