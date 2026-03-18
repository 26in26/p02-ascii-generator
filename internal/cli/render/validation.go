package render

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

func commandValidation(cmd *cobra.Command, args []string) error {
	if _, err := validateOutputPathAndFormat(outputPath, outputFormat); err != nil {
		return err
	}
	if _, err := validateCharset(charset); err != nil {
		return err
	}
	if resizeFlag != "" {
		w, h, err := validateDimentions(resizeFlag)
		if err != nil {
			return err
		}
		width = w
		height = h
	} else {
		if err := validateDimention(width); err != nil {
			return fmt.Errorf("width: %w", err)
		}
		if err := validateDimention(height); err != nil {
			return fmt.Errorf("height: %w", err)
		}
	}

	if _, _, err := validateDimentions(aspectRatio); err != nil {
		return err
	}

	return nil
}

var outputFormatOptions = map[string]string{"text": "", "image": ""}

func validateOutputPathAndFormat(outputPath, outputFormat string) (string, error) {
	val, ok := outputFormatOptions[outputFormat]
	if !ok {
		return "", errors.New("invalid output format")
	}

	if outputPath == "" && outputFormat == "image" {
		return "", errors.New("can't write image to terminal")
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
