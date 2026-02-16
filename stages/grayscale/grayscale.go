package grayscale

import "github.com/26in26/p02-ascii-generator/image"

type grayscaleStage struct{}

func (s *grayscaleStage) Process(src *image.Buffer) *image.Buffer {
	if src == nil {
		panic("grayscale stage must receive non-nil image buffer")
	}

	return src.ToGray()
}

func NewGrayscaleStage() *grayscaleStage {
	return &grayscaleStage{}
}
