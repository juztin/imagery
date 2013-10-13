package imagery

import (
	"bufio"
	"bytes"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"io/ioutil"
	"os"
	"path/filepath"

	"code.google.com/p/graphics-go/graphics"
)

/*func jpegName(img string) string {
    d := filepath.Dir(img)
    f := filepath.Base(img)
    x := filepath.Ext(img)
    n := f[:len(x)+1]+".jpg"

    return filepath.Join(d, n)
}*/

func Decode(fpath string) (image.Image, string, error) {
	f, err := os.Open(fpath)
	if err != nil {
		return nil, "", err
	}
	defer f.Close()

	return image.Decode(bufio.NewReader(f))
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

func WriteToJpg(img image.Image, dstPath, name string) error {
	var buf bytes.Buffer
	jpeg.Encode(&buf, img, &jpeg.Options{Quality: 75})
	return ioutil.WriteFile(filepath.Join(dstPath, name), buf.Bytes(), 0644)
}

func ConvertToJPG(fpath, name string, deleteOrig bool) (image.Image, error) {
	img, _, err := Decode(fpath)
	if err != nil {
		return img, err
	} else if WriteToJpg(img, filepath.Dir(fpath), name); err != nil {
		return img, err
	}

	if deleteOrig {
		os.Remove(fpath)
	}

	return img, nil
}

func ResizeWidthToJPG(fpath, name string, deleteOrig bool, width int) (image.Image, error) {
	img, _, err := Decode(fpath)
	if err != nil {
		return img, err
	} else if img, err = ResizeWidth(img, width); err != nil {
		return img, err
	} else if err = WriteToJpg(img, filepath.Dir(fpath), name); err != nil {
		return img, err
	}

	if deleteOrig {
		os.Remove(fpath)
	}

	//return img.Bounds().Size(), nil
	return img, nil
}
