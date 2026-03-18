package render

import (
	"fmt"
	"image/png"
	"os"

	internalImage "github.com/26in26/p02-ascii-generator/image"
	"github.com/26in26/p02-ascii-generator/stages/ascii"
	"github.com/26in26/p02-ascii-generator/utils/imageio"
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

func writeToFile(path string, ascii *internalImage.AsciiBuffer) error {
	asciiImg := ascii.ToImage()
	f, err := os.Create(path)

	if err != nil {
		return err
	}

	defer f.Close()

	return png.Encode(f, asciiImg)
}

func getDensityCharSet(charsetName string) ascii.DensityCharset {
	switch charsetName {
	case "standard":
		return ascii.StandardCharset
	case "dense":
		return ascii.DenseCharset
	case "dots":
		return ascii.DotsCharset
	default:
		return ascii.StandardCharset
	}
}
