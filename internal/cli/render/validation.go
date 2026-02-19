package render

import (
	"errors"
	"fmt"
)

var outputFormatOptions = map[string]string{"terminal": "", "file": ""}

func validateOutputFormat(outputFormat string) (string, error) {
	val, ok := outputFormatOptions[outputFormat]
	if !ok {
		return "", errors.New("invalid output format")
	}

	return val, nil
}

var charsetOptions = map[string]string{"standard": "", "dense": "", "blocks": ""}

func validateCharset(charset string) (string, error) {
	val, ok := charsetOptions[charset]
	if !ok {
		return "", errors.New("invalid charset")
	}

	return val, nil
}

var colorOptions = map[string]string{"full": "", "Retro": "", "none": ""}

func validateColor(color string) (string, error) {
	val, ok := colorOptions[color]
	if !ok {
		return "", errors.New("invalid charset")
	}

	return val, nil
}

func validateDimention(d int) error {
	if d <= 0 {
		return errors.New("invalid dimention: dimention must be greater than 0")
	}

	return nil
}

func validateDimentions(d string) (int, int, error) {
	var w, h int
	_, err := fmt.Sscanf(d, "%dX%d", &w, &h)

	if err != nil {
		return 0, 0, err
	}
	if err := validateDimention(w); err != nil {
		return 0, 0, fmt.Errorf("width, %w", err)
	}
	if err := validateDimention(h); err != nil {
		return 0, 0, fmt.Errorf("height, %w", err)
	}

	return w, h, nil
}
