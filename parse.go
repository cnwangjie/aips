package aips

import (
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"strings"
)

func Parse(path string) (image.Image, error) {
	var img image.Image = nil
	var err error = nil
	ispng := strings.HasSuffix(path, ".png")
	if ispng {
		img, err = ParsePNG(path)
	}
	isjpeg := strings.HasSuffix(path, ".jpeg")
	isjpg := strings.HasSuffix(path, ".jpg")
	if isjpeg || isjpg {
		img, err = ParseJPEG(path)
	}
	isgif := strings.HasSuffix(path, ".gif")
	if isgif {
		img, err = ParseGIF(path)
	}
	return img, err
}

func ParsePNG(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, err := png.Decode(file)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func ParseJPEG(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, err := jpeg.Decode(file)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func ParseGIF(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, err := gif.Decode(file)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func SavePNG(src image.Image, path string) error {
	o, err := os.Create(path)
	defer o.Close()
	if err != nil {
		return err
	}
	png.Encode(o, src)
	return nil
}

func SaveJPEG(src image.Image, path string) error {
	o, err := os.Create(path)
	defer o.Close()
	if err != nil {
		return err
	}
	jpeg.Encode(o, src, nil)
	return nil
}

func SaveGIF(src image.Image, path string) error {
	o, err := os.Create(path)
	defer o.Close()
	if err != nil {
		return err
	}
	gif.Encode(o, src, nil)
	return nil
}
