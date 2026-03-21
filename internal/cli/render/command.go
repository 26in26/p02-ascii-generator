package render

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/26in26/p02-ascii-generator/filters/charcolor"
	"github.com/26in26/p02-ascii-generator/filters/drawedge"
	internalImage "github.com/26in26/p02-ascii-generator/image"
	"github.com/26in26/p02-ascii-generator/pipeline"
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
	rWidth        int
	rHeight       int
	edgeFlag      bool = true
	edgeThreshold int
	colorFlag     bool = true
	charset       string
	invert        bool
)

// cmd represents the render command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "render",
		Short:   "Render an image to ASCII art",
		Long:    "Render an image to ASCII art with a variety of options.\nUnlocking your ASCII dream! ",
		Args:    cobra.ExactArgs(1),
		PreRunE: commandValidation,
		RunE:    run,
	}

	cmd.Flags().StringVarP(&outputPath, "output", "o", "", "Output file path, if not specified prints to the terminal")
	cmd.Flags().StringVar(&outputFormat, "output-format", "text", "Format for output (text, image)")

	cmd.Flags().StringVar(&resizeFlag, "resize", "", "Output dimensions as WIDTHxHEIGHT. Exclusive with --width, --height, and --aspect-ratio")
	cmd.Flags().IntVar(&width, "width", 0, "Target output width (in characters). Use with --aspect-ratio or automatically preserve image ratio")
	cmd.Flags().IntVar(&height, "height", 0, "Target output height (in characters). Use with --aspect-ratio or automatically preserve image ratio")
	cmd.Flags().StringVar(&aspectRatio, "aspect-ratio", "", "Output aspect ratio in WIDTHxHEIGHT. If not set, uses source image aspect ratio")

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

	cmd.Flags().BoolFunc("color", "Enable full color support", func(s string) error {
		colorFlag = true
		return nil
	})
	cmd.Flags().Lookup("color").NoOptDefVal = "true"
	cmd.Flags().BoolFunc("no-color", "Disable color support", func(s string) error {
		colorFlag = false
		return nil
	})
	cmd.Flags().Lookup("no-color").NoOptDefVal = "false"

	cmd.Flags().StringVar(&charset, "charset", "standard", "charset for ascii art (standard, dense, dots)")
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

	if aspectRatio == "" {
		rWidth = src.Width
		rHeight = src.Height
	}

	input, output, err := buildPipeline(streamCtx, rWidth, rHeight)
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

func buildPipeline(ctx context.Context, rWidth, rHeight int) (flow.Outlet[*internalImage.RGBBuffer], flow.Outlet[*internalImage.AsciiBuffer], error) {
	var resizeStage pipeline.Stage[*internalImage.RGBBuffer, *internalImage.RGBBuffer]
	var err error

	if width > 0 && height > 0 {
		resizeStage, err = resize.NewResizeStage(resize.With(width, height))
	} else if width > 0 {
		resizeStage, err = resize.NewResizeStage(resize.WithWidth(width), resize.WithAspectRatio(rWidth, rHeight, true))
	} else {
		resizeStage, err = resize.NewResizeStage(resize.WithHeight(height), resize.WithAspectRatio(rWidth, rHeight, false))
	}

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

		var asciiArtWriter io.Writer

		if outputPath != "" {
			f, err := os.Create(outputPath)
			if err != nil {
				fmt.Println("Couldn't create file")
				return
			}
			defer f.Close()
			asciiArtWriter = f
		} else {
			asciiArtWriter = os.Stdout
		}

		if outputFormat == "text" {
			var str strings.Builder

			ab.ToString(&str, colorFlag)
			if _, err := fmt.Fprint(asciiArtWriter, str.String()); err != nil {
				fmt.Printf("Couldn't write ascii art %s", outputFormat)
			}
		} else {
			if err := writePngImage(asciiArtWriter, ab, colorFlag); err != nil {
				fmt.Printf("Couldn't write ascii art %s", outputFormat)
			}
		}
	})
}
