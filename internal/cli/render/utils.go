package render

import (
	"fmt"
	"image/png"
	"io"
	"os"

	internalImage "github.com/26in26/p02-ascii-generator/image"
	"github.com/26in26/p02-ascii-generator/stages/ascii"
	"github.com/26in26/p02-ascii-generator/utils/imageio"
	"golang.org/x/term"
)

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

func writePngImage(w io.Writer, ascii *internalImage.AsciiBuffer, useColor bool) error {
	asciiImg := ascii.ToImage(useColor)

	return png.Encode(w, asciiImg)
}

func getDensityCharSet(charsetName string) ascii.DensityCharset {
	switch charsetName {
	case "standard":
		return ascii.StandardCharset
	case "dense":
		return ascii.DenseCharset
	case "blocks":
		return ascii.BlocksCharset
	default:
		return ascii.StandardCharset
	}
}

func getTerminalWidth() (int, error) {
	fd := int(os.Stdout.Fd())

	// Check if the file descriptor is a terminal before attempting to get the size
	if !term.IsTerminal(fd) {
		return 0, fmt.Errorf("Not running in a terminal, cannot get size.")
	}

	width, _, err := term.GetSize(fd)
	if err != nil {

		return 0, fmt.Errorf("Error getting terminal size: %w", err)
	}

	return width, nil
}
