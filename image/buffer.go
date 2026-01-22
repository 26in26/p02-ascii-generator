package image

type Buffer struct {
	width  int
	height int
	stride int

	format Format
	data   []byte
}

func NewBuffer(width, height int, format Format) *Buffer {
	if width < 0 || height < 0 {
		panic("image: negative width or height")
	}

	bpp := format.BytesPerPixel()
	stride := width * bpp

	return &Buffer{
		width:  width,
		height: height,
		stride: stride,
		format: format,
		data:   make([]byte, stride*height),
	}

}

func (b *Buffer) Width() int {
	return b.width
}

func (b *Buffer) Height() int {
	return b.height
}

func (b *Buffer) Stride() int {
	return b.stride
}

func (b *Buffer) Data() []byte {
	return b.data
}

func (b *Buffer) Index(x, y int) int {
	return y*b.stride + x*b.format.BytesPerPixel()
}
