package ascii

func pixelToASCII(densityCharset []byte, gray byte, invert bool) byte {
	index := int((float64(gray) / 255) * float64(len(densityCharset)-1))
	if invert {
		index = (len(densityCharset) - 1) - index
	}
	return densityCharset[index]
}
