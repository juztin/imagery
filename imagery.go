// Copyright 2013 Justin Wilson. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package imagery

import (
	"bufio"
	"bytes"
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"

	"code.google.com/p/graphics-go/graphics"
)

type ImgType int

const (
	IMG_UNKNOWN ImgType = iota
	IMG_PNG
	IMG_JPEG
	IMG_GIF
)

func Decode(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		//return nil, "", err
		return nil, err
	}
	defer f.Close()

	i, _, err := image.Decode(bufio.NewReader(f))
	return i, err
}

func ResizeWidth(img image.Image, width int) (image.Image, error) {
	s := img.Bounds().Size()
	h := int((float32(width) / float32(s.X)) * float32(s.Y))
	r := image.NewRGBA(image.Rect(0, 0, width, h))

	if err := graphics.Scale(r, img); err != nil {
		return nil, err
	}
	return r, nil
}

func WriteTo(t ImgType, m image.Image, path, filename string) error {
	var buf bytes.Buffer
	switch t {
	default:
		return errors.New("Invalid image type")
	case IMG_PNG:
		png.Encode(&buf, m)
	case IMG_JPEG:
		//jpeg.Encode(&buf, m, &jpeg.Options{Quality: 75})
		jpeg.Encode(&buf, m, nil)
	case IMG_GIF:
		gif.Encode(&buf, m, nil)
	}
	return ioutil.WriteFile(filepath.Join(path, filename), buf.Bytes(), 0644)
}

func WriteToPng(img image.Image, path, filename string) error {
	return WriteTo(IMG_PNG, img, path, filename)
}

func WriteToJpg(img image.Image, path, filename string) error {
	return WriteTo(IMG_JPEG, img, path, filename)
}

func WriteToGif(img image.Image, path, filename string) error {
	return WriteTo(IMG_GIF, img, path, filename)
}

func ConvertToJPG(path, filename string, deleteOrig bool) (image.Image, error) {
	//img, _, err := Decode(fpath)
	img, err := Decode(path)
	if err != nil {
		return img, err
	} else if WriteToJpg(img, filepath.Dir(path), filename); err != nil {
		return img, err
	}

	if deleteOrig {
		os.Remove(path)
	}

	return img, nil
}

func ResizeWidthToJPG(path, filename string, deleteOrig bool, width int) (image.Image, error) {
	img, err := Decode(path)
	if err != nil {
		return img, err
	} else if img, err = ResizeWidth(img, width); err != nil {
		return img, err
	} else if err = WriteToJpg(img, filepath.Dir(path), filename); err != nil {
		return img, err
	}

	if deleteOrig {
		os.Remove(path)
	}

	return img, nil
}
