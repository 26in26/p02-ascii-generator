package render

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func commandValidation(cmd *cobra.Command, args []string) error {
	if _, err := validateOutputPathAndFormat(outputPath, outputFormat); err != nil {
		return err
	}
	if _, err := validateCharset(charset); err != nil {
		return err
	}

	widthSet := cmd.Flags().Changed("width")
	heightSet := cmd.Flags().Changed("height")
	aspectRatioSet := cmd.Flags().Changed("aspect-ratio")
	resizeSet := cmd.Flags().Changed("resize")

	if resizeSet {
		if widthSet || heightSet {
			return errors.New("--resize cannot be combined with --width or --height")
		}
		if aspectRatioSet {
			return errors.New("--resize cannot be combined with --aspect-ratio")
		}

		w, h, err := validateDimentions(resizeFlag)
		if err != nil {
			return fmt.Errorf("resize: %w", err)
		}
		width = w
		height = h
		return nil
	}

	if !widthSet && !heightSet {
		terminalWidth, err := getTerminalWidth()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting terminal width: %v\nUsing default width 190.\nYou can specify --width, --height or --resize flags instead.", err)
			width = 190
		}
		width = terminalWidth
	}

	if widthSet {
		if err := validateDimention(width); err != nil {
			return fmt.Errorf("width: %w", err)
		}
	}

	if heightSet {
		if err := validateDimention(height); err != nil {
			return fmt.Errorf("height: %w", err)
		}
	}

	if aspectRatioSet {
		if widthSet && heightSet {
			return errors.New("--aspect-ratio cannot be combined with both --width and --height")
		}
		w, h, err := validateDimentions(aspectRatio)
		if err != nil {
			return fmt.Errorf("aspect-ratio: %w", err)
		}
		rWidth = w
		rHeight = h
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
