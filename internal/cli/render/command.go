package render

import (
	"context"
	"fmt"
	"image/png"
	"os"
	"sync"

	internalImage "github.com/26in26/p02-ascii-generator/image"
	"github.com/26in26/p02-ascii-generator/imageio"
	"github.com/26in26/p02-ascii-generator/pipeline/flow"
	"github.com/26in26/p02-ascii-generator/stages/ascii"
	"github.com/26in26/p02-ascii-generator/stages/edge"
	"github.com/26in26/p02-ascii-generator/stages/grayscale"
	"github.com/26in26/p02-ascii-generator/stages/resize"
	"github.com/spf13/cobra"
)

var (
	output        string
	outputFormat  string
	resizeFlag    string
	width         int
	height        int
	aspectRatio   string
	edgeFlag      bool
	edgeThreshold int
	color         string
	charset       string
	invert        bool
)

// cmd represents the render command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "render",
		Short: "Render an image/video to ASCII art",
		Long:  "Render an image/video to ASCII art",
		Args:  cobra.MaximumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if _, err := validateOutputFormat(outputFormat); err != nil {
				return err
			}
			if _, err := validateCharset(charset); err != nil {
				return err
			}
			if _, err := validateColor(color); err != nil {
				return err
			}
			if resizeFlag != "" {
				w, h, err := validateDimentions(resizeFlag)
				if err != nil {
					return err
				}
				width = w
				height = h
			} else {
				if err := validateDimention(width); err != nil {
					return fmt.Errorf("width: %w", err)
				}
				if err := validateDimention(height); err != nil {
					return fmt.Errorf("height: %w", err)
				}
			}

			if _, _, err := validateDimentions(aspectRatio); err != nil {
				return err
			}

			return nil

		},
		RunE: func(cmd *cobra.Command, args []string) error {
			imgPath := "./input.png"

			if len(args) >= 1 {
				imgPath = args[0]
			}

			src, err := loadImage(imgPath)
			if err != nil {
				return err
			}

			// stages
			resizeStage, err := resize.NewResizeStage(resize.WithWidth(190),
				resize.WithAspectRatio(src.Width, src.Height, true),
			)
			if err != nil {
				return fmt.Errorf("internal error, %w", err)
			}

			grayscaleStage := grayscale.NewGrayscaleStage()
			edgeStage := edge.NewSobelEdgeDetectionStage()
			asciiStage := ascii.NewAsciiStage(ascii.WithInvert(invert))

			// filters
			edgeFilter := ascii.NewEdgeFilter(ascii.WithEdgeThreshold(edgeThreshold))
			colorFilter := ascii.NewColorFilter()

			// stream
			streamCtx := context.Background()

			rawStream := flow.NewOutlet[*internalImage.RGBBuffer](10)
			resizeStream := flow.Map(streamCtx, &rawStream, resizeStage)
			resizeBranches := resizeStream.Branch(streamCtx, 2)

			// Path for grayscale -> ascii characters
			grayscaleStream := flow.Map(streamCtx, &resizeBranches[0], grayscaleStage)
			grayBranches := grayscaleStream.Branch(streamCtx, 2)
			edgeStream := flow.Map(streamCtx, &grayBranches[0], edgeStage)
			asciiStream := flow.Map(streamCtx, &grayBranches[1], asciiStage)

			// Combine streams
			edgeZip := flow.Zip(streamCtx, asciiStream, edgeStream)
			asciiStream = flow.Mask(streamCtx, edgeZip, edgeFilter)
			colorZip := flow.Zip(streamCtx, asciiStream, resizeBranches[1])
			asciiStream = flow.Mask(streamCtx, colorZip, colorFilter)

			var wg sync.WaitGroup
			wg.Add(1)

			asciiStream.Sink(func(ab *internalImage.AsciiBuffer, err error) {
				if err != nil {
					fmt.Println(err)
					return
				}

				err = writeToFile(ab)
				if err != nil {
					fmt.Println("Couldn't write ascii art to file")
				}

				wg.Done()

			})

			RGBBufferChan := make(chan *internalImage.RGBBuffer)
			defer close(RGBBufferChan)
			rawStream.Feed(streamCtx, RGBBufferChan)

			RGBBufferChan <- src

			wg.Wait()

			return nil
		},
	}

	cmd.Flags().StringVarP(&output, "output", "o", "", "Output file path")
	cmd.Flags().StringVar(&outputFormat, "output-format", "terminal", "Format for output")

	cmd.Flags().StringVar(&resizeFlag, "resize", "", "Output new dimensions, resize overide width and height")
	cmd.Flags().IntVar(&width, "width", 100, "Output width dimension")
	cmd.Flags().IntVar(&height, "height", 100, "Output height dimension")
	cmd.Flags().StringVar(&aspectRatio, "aspect-ratio", "1X1", "")

	cmd.Flags().BoolVar(&edgeFlag, "edge", true, "Enable edge detection")
	cmd.Flags().BoolVar(&edgeFlag, "no-edge", true, "Disable edge detection")
	cmd.Flags().Lookup("no-edge").NoOptDefVal = "false"
	cmd.Flags().IntVarP(&edgeThreshold, "threshold", "t", 23, "Edge detection threshold")

	cmd.Flags().StringVar(&color, "color", "none", "Select color mode")
	cmd.Flags().StringVar(&charset, "charset", "standard", "")
	cmd.Flags().BoolVar(&invert, "invert", false, "Invert brightness to character mapping")

	return cmd
}

func loadImage(imgPath string) (*internalImage.RGBBuffer, error) {
	img, err := imageio.LoadImageFromFile(imgPath)
	if err != nil {
		return nil, fmt.Errorf("Couldn't load %s, %w", imgPath, err)
	}

	src, err := imageio.ConvertToRGBBuffer(img)
	if err != nil {
		return nil, fmt.Errorf("Internal error: \n%w", err)
	}
	return src, nil
}

func writeToFile(ascii *internalImage.AsciiBuffer) error {
	asciiImg := ascii.ToImage()
	f, err := os.Create("ascii.png")

	if err != nil {
		return err
	}

	defer f.Close()

	return png.Encode(f, asciiImg)
}
