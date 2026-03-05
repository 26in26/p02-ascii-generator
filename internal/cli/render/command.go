package render

import (
	"context"
	"fmt"

	"github.com/26in26/p02-ascii-generator/imageio"
	"github.com/26in26/p02-ascii-generator/pipeline"
	"github.com/26in26/p02-ascii-generator/stages/ascii"
	"github.com/26in26/p02-ascii-generator/stages/edge"
	"github.com/26in26/p02-ascii-generator/stages/grayscale"
	"github.com/26in26/p02-ascii-generator/stages/resize"
	"github.com/spf13/cobra"
)

var (
	output       string
	outputFormat string
	resizeFlag   string
	width        int
	height       int
	aspectRatio  string
	edgeFlag     bool
	color        string
	charset      string
	invert       bool
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

			img, err := imageio.LoadImageFromFile(imgPath)
			if err != nil {
				return fmt.Errorf("couldn't load %s, %w", imgPath, err)
			}

			src, err := imageio.ConvertToRGBBuffer(img)
			if err != nil {
				return fmt.Errorf("internal error: \n%w", err)
			}

			stages := make([]pipeline.Stage, 0, 5)

			resizeStage, err := resize.NewResizeStage(resize.WithWidth(190),
				resize.WithAspectRatio(src.Width, src.Height, true),
			)
			if err != nil {
				return fmt.Errorf("internal error, %w", err)
			}

			grayscaleStage := grayscale.NewGrayscaleStage()
			stages = append(stages, resizeStage, grayscaleStage)

			if edgeFlag {
				edgeDetection := edge.NewSobelEdgeDetectionStage()
				stages = append(stages, edgeDetection)
			}

			asciiStage := ascii.NewAsciiStage(ascii.WithEdgeThreshold(23),
				ascii.WithColorMode(ascii.FullColor),
				ascii.WithInvert(false),
				ascii.WithEdge(edgeFlag),
			)
			stages = append(stages, asciiStage)

			// Create pipeline
			p, err := pipeline.NewPipeline(
				stages...,
			)

			p.Enqueue(context.Background(), src)
			result := <-p.Results()
			asciiArt, err := result.AsciiArt.Get(context.Background())
			if err != nil {
				return fmt.Errorf("internal error: %w", err)
			}
			print(asciiArt)
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

	cmd.Flags().StringVar(&color, "color", "none", "Select color mode")
	cmd.Flags().StringVar(&charset, "charset", "standard", "")
	cmd.Flags().BoolVar(&invert, "invert", false, "Invert brightness to character mapping")

	return cmd
}
