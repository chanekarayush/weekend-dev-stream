package main

import (
	"bytes"
	"io"
	"log"
	"os"
)

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
