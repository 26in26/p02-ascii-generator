package imageio

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

func LoadImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decodedImg, _, err := image.Decode(file)

	return decodedImg, err

}
