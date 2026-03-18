package ascii

import (
	"context"

	"github.com/26in26/p02-ascii-generator/image"
	"github.com/26in26/p02-ascii-generator/pipeline"
)

type AsciiStage struct {
	*pipeline.BaseStage[*image.GrayBuffer, *image.AsciiBuffer]
	asciiArtPool   *asciiArtPool
	invert         bool
	densityCharset []byte
}

func NewAsciiStage(densityCharset []byte, invert bool) pipeline.Stage[*image.GrayBuffer, *image.AsciiBuffer] {
	return &AsciiStage{
		BaseStage:      pipeline.NewBaseStage[*image.GrayBuffer, *image.AsciiBuffer]("Ascii"),
		asciiArtPool:   NewAsciiArtPool(),
		invert:         invert,
		densityCharset: densityCharset,
	}
}

func (s *AsciiStage) Kernal(ctx context.Context, input *image.GrayBuffer) (*image.AsciiBuffer, error) {
	asciiArt, err := s.asciiArtPool.Get(size{input.Width, input.Height})
	if err != nil {
		return nil, err
	}

	grayIndex := 0
	ArtIndex := 0

	for y := 0; y < asciiArt.Height; y++ {
		for x := 0; x < asciiArt.Width; x++ {
			asciiArt.Data[ArtIndex] = pixelToASCII(s.densityCharset, input.Data[grayIndex], s.invert)

			grayIndex++
			ArtIndex += 4
		}
	}

	return asciiArt, nil
}

func (s *AsciiStage) Release(asciiArt *image.AsciiBuffer) {
	if asciiArt == nil {
		return
	}

	s.asciiArtPool.Put(size{asciiArt.Width, asciiArt.Height}, asciiArt)
}
