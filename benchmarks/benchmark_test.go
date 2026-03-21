package benchmarks_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/26in26/p02-ascii-generator/filters/charcolor"
	"github.com/26in26/p02-ascii-generator/filters/drawedge"
	internalImage "github.com/26in26/p02-ascii-generator/image"
	"github.com/26in26/p02-ascii-generator/pipeline/flow"
	"github.com/26in26/p02-ascii-generator/stages/ascii"
	"github.com/26in26/p02-ascii-generator/stages/edge"
	"github.com/26in26/p02-ascii-generator/stages/grayscale"
	"github.com/26in26/p02-ascii-generator/stages/resize"
)

func TestPipelineThroughput(t *testing.T) {
	tests := []struct {
		name     string
		withEdge bool
	}{
		{
			name:     "Complex Pipeline (With Edge Detection)",
			withEdge: true,
		},
		{
			name:     "Simple Pipeline (No Edge Detection)",
			withEdge: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Benchmark configuration
			const (
				duration  = 10 * time.Second
				minFPS    = 30.0
				imgWidth  = 1920
				imgHeight = 1080
				targetW   = 200
				targetH   = 100
			)

			// Setup context and pipeline
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			input, output, err := buildBenchmarkPipeline(ctx, tt.withEdge, targetW, targetH)
			if err != nil {
				t.Fatalf("Failed to build pipeline: %v", err)
			}

			// Prepare source image
			src, err := internalImage.NewRGBBuffer(imgWidth, imgHeight)
			// Fill with simple pattern to ensure we aren't just processing zeros
			for i := range src.Data {
				src.Data[i] = uint8(i % 255)
			}

			// Start input feeder
			inputChan := make(chan *internalImage.RGBBuffer, 50)
			go func() {
				for {
					select {
					case <-ctx.Done():
						close(inputChan)
						return
					case inputChan <- src:
					}
				}
			}()
			input.Feed(ctx, inputChan)

			// Start output consumer (Sink)
			var frameCount int64
			done := make(chan struct{})

			go func() {
				output.Sink(func(ab *internalImage.AsciiBuffer, err error) {
					if err == nil {
						atomic.AddInt64(&frameCount, 1)
					}
				})
				close(done)
			}()

			t.Logf("Running benchmark for %v...", duration)
			start := time.Now()

			// Run for fixed duration
			<-time.After(duration)
			cancel()

			// Wait for Sink to finish processing remaining items
			<-done

			elapsed := time.Since(start)
			if elapsed < duration {
				elapsed = duration
			}

			totalFrames := atomic.LoadInt64(&frameCount)
			fps := float64(totalFrames) / elapsed.Seconds()

			t.Logf("Result: Processed %d frames in %v", totalFrames, elapsed)
			t.Logf("FPS: %.2f", fps)

			if fps < minFPS {
				t.Errorf("Performance insufficient: %.2f FPS (expected > %.2f)", fps, minFPS)
			}
		})
	}
}

func buildBenchmarkPipeline(ctx context.Context, withEdge bool, w, h int) (flow.Outlet[*internalImage.RGBBuffer], flow.Outlet[*internalImage.AsciiBuffer], error) {
	// Replicating command.go logic structure
	resizeStage, err := resize.NewResizeStage(resize.With(w, h))
	if err != nil {
		return flow.Outlet[*internalImage.RGBBuffer]{}, flow.Outlet[*internalImage.AsciiBuffer]{}, err
	}

	grayscaleStage := grayscale.NewGrayscaleStage()
	asciiStage := ascii.NewAsciiStage(ascii.WithDensityCharset(ascii.StandardCharset))
	colorFilter := charcolor.NewColorFilter()

	// Pipeline construction
	rawStream := flow.NewOutlet[*internalImage.RGBBuffer](100)
	resizeStream := flow.Map(ctx, &rawStream, resizeStage)
	resizeBranches := resizeStream.Branch(ctx, 2)

	grayscaleStream := flow.Map(ctx, &resizeBranches[0], grayscaleStage)
	var asciiStream flow.Outlet[*internalImage.AsciiBuffer]

	if withEdge {
		edgeStage := edge.NewSobelEdgeDetectionStage()
		edgeFilter := drawedge.NewEdgeFilter(drawedge.WithEdgeThreshold(50))

		grayBranches := grayscaleStream.Branch(ctx, 2)
		edgeStream := flow.Map(ctx, &grayBranches[0], edgeStage)
		baseAsciiStream := flow.Map(ctx, &grayBranches[1], asciiStage)

		edgeZip := flow.Zip(ctx, baseAsciiStream, edgeStream)
		asciiStream = flow.Mask(ctx, edgeZip, edgeFilter)
	} else {
		asciiStream = flow.Map(ctx, &grayscaleStream, asciiStage)
	}

	colorZip := flow.Zip(ctx, asciiStream, resizeBranches[1])
	finalStream := flow.Mask(ctx, colorZip, colorFilter)

	return rawStream, finalStream, nil
}
