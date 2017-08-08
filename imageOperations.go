package main

import (
	"bytes"
	"image"
	"image/png"
)

func dataToPNG(data []byte, imageName string) (image.Image, error) {
	reader := bytes.NewReader(data)
	img, err := png.Decode(reader)
	if err != nil {
		return img, err
	}
	return img, err
}
