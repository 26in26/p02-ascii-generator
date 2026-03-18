package render

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/26in26/p02-ascii-generator/filters/charcolor"
	"github.com/26in26/p02-ascii-generator/filters/drawedge"
	internalImage "github.com/26in26/p02-ascii-generator/image"
	"github.com/26in26/p02-ascii-generator/pipeline/flow"
	"github.com/26in26/p02-ascii-generator/stages/ascii"
	"github.com/26in26/p02-ascii-generator/stages/edge"
	"github.com/26in26/p02-ascii-generator/stages/grayscale"
	"github.com/26in26/p02-ascii-generator/stages/resize"
	"github.com/spf13/cobra"
)

var (
	outputPath    string
	outputFormat  string
	resizeFlag    string
	width         int
	height        int
	aspectRatio   string
	edgeFlag      bool
	edgeThreshold int
	colotFlag     bool
	charset       string
	invert        bool
)

// cmd represents the render command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "render",
		Short:   "Render an image to ASCII art",
		Long:    "Render an image to ASCII art with a variety of options.\nUnlocking youre ASCII dream! ",
		Args:    cobra.ExactArgs(1),
		PreRunE: commandValidation,
		RunE:    run,
	}

	cmd.Flags().StringVarP(&outputPath, "output", "o", "", "Output file path, if not specified prints to the terminal")
	cmd.Flags().StringVar(&outputFormat, "output-format", "text", "Format for output (image, text)")

	cmd.Flags().StringVar(&resizeFlag, "resize", "", "Output new dimensions, resize overide width and height")
	cmd.Flags().IntVar(&width, "width", 100, "Output width dimension")
	cmd.Flags().IntVar(&height, "height", 100, "Output height dimension")
	cmd.Flags().StringVar(&aspectRatio, "aspect-ratio", "1X1", "")

	cmd.Flags().BoolFunc("edge", "Enable edge detection", func(s string) error {
		edgeFlag = true
		return nil
	})
	cmd.Flags().Lookup("edge").NoOptDefVal = "true"
	cmd.Flags().BoolFunc("no-edge", "Disable edge detection", func(s string) error {
		edgeFlag = false
		return nil
	})
	cmd.Flags().Lookup("no-edge").NoOptDefVal = "false"
	cmd.Flags().IntVarP(&edgeThreshold, "threshold", "t", 23, "Edge detection threshold")

	cmd.Flags().BoolVar(&colotFlag, "color", true, "Enable full color")
	cmd.Flags().BoolVar(&colotFlag, "no-color", true, "Disable edge detection")
	cmd.Flags().Lookup("no-color").NoOptDefVal = "false"

	cmd.Flags().StringVar(&charset, "charset", "standard", "")
	cmd.Flags().BoolVar(&invert, "invert", false, "Invert brightness to character mapping")

	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	imgPath := args[0]

	src, err := loadImage(imgPath)
	if err != nil {
		return err
	}

	streamCtx := context.Background()

	input, output, err := buildPipeline(streamCtx, src.Width, src.Height)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(1)

	handleOutput(output, &wg)

	bufferChan := make(chan *internalImage.RGBBuffer, 1)
	input.Feed(streamCtx, bufferChan)
	bufferChan <- src
	close(bufferChan)

	wg.Wait()

	return nil
}

func buildPipeline(ctx context.Context, srcWidth, srcHeight int) (flow.Outlet[*internalImage.RGBBuffer], flow.Outlet[*internalImage.AsciiBuffer], error) {
	resizeStage, err := resize.NewResizeStage(resize.WithWidth(width),
		resize.WithAspectRatio(srcWidth, srcHeight, true),
	)
	if err != nil {
		return flow.Outlet[*internalImage.RGBBuffer]{}, flow.Outlet[*internalImage.AsciiBuffer]{}, fmt.Errorf("internal error, %w", err)
	}

	grayscaleStage := grayscale.NewGrayscaleStage()
	edgeStage := edge.NewSobelEdgeDetectionStage()
	asciiStage := ascii.NewAsciiStage(ascii.WithInvert(invert), ascii.WithDensityCharset(getDensityCharSet(charset)))

	edgeFilter := drawedge.NewEdgeFilter(drawedge.WithEdgeThreshold(edgeThreshold))
	colorFilter := charcolor.NewColorFilter()

	rawStream := flow.NewOutlet[*internalImage.RGBBuffer](10)
	resizeStream := flow.Map(ctx, &rawStream, resizeStage)
	resizeBranches := resizeStream.Branch(ctx, 2)

	grayscaleStream := flow.Map(ctx, &resizeBranches[0], grayscaleStage)
	var asciiStream flow.Outlet[*internalImage.AsciiBuffer]

	if edgeFlag {
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

func handleOutput(stream flow.Outlet[*internalImage.AsciiBuffer], wg *sync.WaitGroup) {
	stream.Sink(func(ab *internalImage.AsciiBuffer, err error) {
		defer wg.Done()
		if err != nil {
			fmt.Println(err)
			return
		}

		if outputFormat == "text" {
			var str strings.Builder

			ab.ToString(&str, colotFlag)
			fmt.Println(str.String())
		} else {
			if err := writeToFile("ascii.png", ab); err != nil {
				fmt.Println("Couldn't write ascii art to file")
			}
		}
	})
}
