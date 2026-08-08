package main

import (
	"flag"
)

func main() {
	var height int
	var width int
	flag.IntVar(&width, "wt", WIDTH, "Width of the Video Stream")
	flag.IntVar(&height, "ht", HEIGHT, "Height of the Video Stream")
	flag.Parse()

	const ROOT = "assets/"

	frames := readFrames(width, height)
	convertedFrames := convertRGBtoYUV(width, height, frames)
	deflatedBytes := findKeyframeDeflate(convertedFrames)
	writeFileToDisk(convertedFrames, ROOT+RGBTOYUVCONVERTED)
	writeBytesToDisk(deflatedBytes, ROOT+ENCODED)
	bytesToInflate := readBytesFromDisk(ROOT + ENCODED)
	decodedFrames := inflateKeyframeDelta(bytesToInflate, width, height)
	reconvertedFrames := convertYUVtoRGB(width, height, decodedFrames)
	writeFileToDisk(decodedFrames, ROOT+DECODED)
	writeFileToDisk(reconvertedFrames, ROOT+YUVTORGBCONVERTED)
}
