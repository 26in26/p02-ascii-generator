package image

type Format uint8

const (
	FormatRGB Format = iota
	FormatRGBA
	FormatGray
	FormatFloatGray
)

func (f Format) BytesPerPixel() int {
	switch f {
	case FormatRGB:
		return 3
	case FormatRGBA:
		return 4
	case FormatGray:
		return 1
	case FormatFloatGray:
		return 4
	default:
		panic("image: unknown format")
	}
}
