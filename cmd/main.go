package main

import (
	"fmt"
	"time"

	"github.com/26in26/p02-ascii-generator/mapping"
	"github.com/26in26/p02-ascii-generator/stages/edge"
	"github.com/26in26/p02-ascii-generator/stages/grayscale"
	"github.com/26in26/p02-ascii-generator/stages/resize"

	"github.com/26in26/p02-ascii-generator/image"
	"github.com/26in26/p02-ascii-generator/imageio"
	"github.com/26in26/p02-ascii-generator/pipeline"
)

func NewTestImage(w, h int) *image.Buffer {
	img := image.NewBuffer(w, h, image.FormatRGB)
	bpp := img.Channels

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			idx := (y*w + x) * bpp

			img.Data[idx+0] = uint8(x * 255 / w) // R
			img.Data[idx+1] = uint8(y * 255 / h) // G
			img.Data[idx+2] = 0                  // B
		}
	}

	return img
}

func main() {
	// Run pipeline and count iterations in 30 seconds
	pngImage, _ := imageio.LoadPNG("./phinix.png")
	src := imageio.ConvertToBuffer(pngImage, image.FormatRGB)

	resizeStage := resize.NewResizeStage(190, 100).PreserveAspectRatio(src.Width, src.Height, true, false)
	grayscaleStage := grayscale.NewGrayscaleStage()

	edgeDetection := edge.NewSobelEdgeDetectionStage()

	// Create pipeline
	p := pipeline.New(
		resizeStage,
		grayscaleStage,
		edgeDetection,
	)

	mapping.PrintAsASCII(p.Run(src), edgeDetection.Gradient)
	// benchmark()

}

func benchmark() {
	start := time.Now()
	iterations := 0
	pngImage, _ := imageio.LoadPNG("./screenshot.png")
	src := imageio.ConvertToBuffer(pngImage, image.FormatRGB)

	resizeStage := resize.NewResizeStage(190, 100).PreserveAspectRatio(src.Width, src.Height, true, false)
	grayscaleStage := grayscale.NewGrayscaleStage()
	edgeDetection := edge.NewSobelEdgeDetectionStage()

	for time.Since(start) < 10*time.Second {

		// Create pipeline
		p := pipeline.New(
			resizeStage,
			grayscaleStage,
			edgeDetection,
		)

		mapping.PrintAsASCII(p.Run(src), edgeDetection.Gradient)
		iterations++
	}

	fmt.Printf("Completed %d iterations in 10 seconds. FPS: %f\n", iterations, float64(iterations)/10)
}
