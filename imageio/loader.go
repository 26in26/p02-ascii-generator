package imageio

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
)

func LoadImageFromFile(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decodedImg, _, err := image.Decode(file)

	return decodedImg, err

}

func LoadImageFromStream(r io.Reader) (image.Image, error) {
	decodedImg, _, err := image.Decode(r)

	return decodedImg, err
}
