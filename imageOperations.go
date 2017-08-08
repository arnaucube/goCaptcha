package main

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
)

func dataToImage(data []byte, imageExtension string) (image.Image, error) {
	reader := bytes.NewReader(data)
	//img, err := png.Decode(reader)
	var img image.Image
	var err error
	switch imageExtension {
	case "png":
		img, err = png.Decode(reader)
	case "jpg", "jpeg":
		img, err = jpeg.Decode(reader)
	default:
		img = nil
	}
	if err != nil {
		fmt.Println(err)
		return img, err
	}
	return img, err
}
