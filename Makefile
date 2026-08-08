
FRAMERATE := 25
HEIGHT := 832
WIDTH := 624
RGBTOYUV := ./assets/rgbtoyuv.yuv
YUVTORGB := ./assets/yuvtorgb.rgb24
DECODED := ./assets/original.rgb

# Default target / all target runs even if you just write `make`
all: convert run
play: YUV RGB Decoded

run: 
	@echo "Running the Video Encoder Program..."
	@cat ./assets/fort_video.rgb24 | go run constants.go decoder.go encoder.go utils.go main.go -wt $(WIDTH) -ht $(HEIGHT) 

build: 
	@echo "Building the Video Encoder Program..."
	@go build -o codec constants.go decoder.go encoder.go utils.go main.go 

clean:
	@echo "Deleting Program Generated Files..."
	@rm ./assets/encoded.yuv ./assets/original.rgb ./assets/rgbtoyuv.yuv ./assets/yuvtorgb.rgb24

YUV:
	@echo "Playing YUV Space Coverted file..."
	@ffplay -f rawvideo -pixel_format yuv420p -video_size $(WIDTH)x$(HEIGHT) -framerate $(FRAMERATE) $(RGBTOYUV)

RGB:
	@echo "Playing RGB Space Re-coverted file..."
	@ffplay -f rawvideo -pixel_format rgb24 -video_size $(WIDTH)x$(HEIGHT) -framerate $(FRAMERATE) $(YUVTORGB)

Decoded:
	@echo "Playing fully decoded file..."
	@ffplay -f rawvideo -pixel_format rgb24 -video_size $(WIDTH)x$(HEIGHT) -framerate $(FRAMERATE) $(DECODED)

convert:
	@ffmpeg  -i ./assets/fort_video.mp4 -f rawvideo -pix_fmt rgb24 ./assets/fort_video.rgb24
