package ascii

import (
	"strconv"
	"strings"
)

const DENSITY = " .-=+*x#$&X@"
const RESET = "\x1b[0m"

func pixelToASCII(gray byte, invert bool) byte {
	index := int((float64(gray) / 255) * float64(len(DENSITY)-1))
	if invert {
		index = (len(DENSITY) - 1) - index
	}
	return DENSITY[index]
}

func getAngleChar(gx, gy int) byte {
	// Fast angle approximation using integer math to avoid expensive Atan2.
	// We classify the gradient into one of four directions.
	abs_gx := gx
	if abs_gx < 0 {
		abs_gx = -abs_gx
	}
	abs_gy := gy
	if abs_gy < 0 {
		abs_gy = -abs_gy
	}

	// If gradient is mostly horizontal -> vertical edge
	if abs_gx > abs_gy<<1 { // abs_gx > 2 * abs_gy
		return '|'
	}
	// If gradient is mostly vertical -> horizontal edge
	if abs_gy > abs_gx<<1 { // abs_gy > 2 * abs_gx
		return '_'
	}
	// Otherwise, it's a diagonal edge
	if (gx > 0) == (gy > 0) { // gx and gy have the same sign
		return '\\'
	}
	return '/' // gx and gy have different signs
}

type ColorMode uint8

const (
	FullColor ColorMode = iota
	RetroColor
	Monochrome
)

// renderFullColor writes a 24-bit ANSI color escape code if the color has changed.
func renderFullColor(builder *strings.Builder, r, g, b byte, data []byte, rgbIndex int) (byte, byte, byte) {
	curR, curG, curB := data[rgbIndex], data[rgbIndex+1], data[rgbIndex+2]
	if r != curR || g != curG || b != curB {
		r, g, b = curR, curG, curB
		builder.WriteString("\x1b[38;2;")
		builder.WriteString(strconv.Itoa(int(r)))
		builder.WriteByte(';')
		builder.WriteString(strconv.Itoa(int(g)))
		builder.WriteByte(';')
		builder.WriteString(strconv.Itoa(int(b)))
		builder.WriteByte('m')
	}
	return r, g, b
}

// renderMonochrome does nothing, effectively producing monochrome output.
func renderMonochrome(builder *strings.Builder, r, g, b byte, data []byte, rgbIndex int) (byte, byte, byte) {
	return r, g, b
}
