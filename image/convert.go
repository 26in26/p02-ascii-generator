package image

func (b *Buffer) ToGray() (*Buffer, error) {
	dst, err := NewBuffer(b.Width, b.Height, FormatGray)

	if err != nil {
		return nil, err
	}

	bpp := b.Channels
	srcData := b.Data
	dstData := dst.Data

	dstIndex := 0
	end := len(srcData) - bpp

	for i := 0; i <= end; i += bpp {
		r := int(srcData[i])
		g := int(srcData[i+1])
		b := int(srcData[i+2])

		// Convert to grayscale using integer approximation (0.21 R + 0.72 G + 0.07 B)
		// Scaled by 256: 54 R + 184 G + 18 B. Sum = 256.
		gray := (r*54 + g*184 + b*18) >> 8
		dstData[dstIndex] = uint8(gray)

		dstIndex++
	}

	return dst, nil
}
