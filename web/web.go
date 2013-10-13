package web

import (
	"crypto/rand"
	"fmt"
	"image"
	"io"
	//"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"bitbucket.org/juztin/imagery"
)

const chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func RandName(l int) string {
	var name = make([]byte, l)
	rand.Read(name)
	for i, c := range name {
		name[i] = chars[c%byte(len(chars))]
	}
	return string(name)
}

func SaveFile(req *http.Request, resp http.ResponseWriter, filePath, filename string) (string, *os.File, error) {
	// create the project path
	if err := os.MkdirAll(filePath, 0755); err != nil {
		return "", nil, err
	}

	// replace spaces in image name
	filename = strings.Replace(filename, " ", "-", -1)

	// open image file for writing
	t, err := os.OpenFile(filepath.Join(filePath, filename), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return filename, nil, err
	}
	defer t.Close()

	// copy the bytes to the file
	r := http.MaxBytesReader(resp, req.Body, 2<<20)
	defer req.Body.Close()
	if _, err := io.Copy(t, r); err != nil {
		return filename, nil, err
	}

	// return the filename
	return filename, t, nil
}

func SaveFormFile(r *http.Request, filePath string) (string, *os.File, error) {
	// grab the file from the request
	f, h, err := r.FormFile("image")
	if err != nil {
		return "", nil, err
	}
	defer f.Close()

	// create the project path
	if err := os.MkdirAll(filePath, 0755); err != nil {
		return "", nil, err
	}

	// replace spaces in image name
	filename := strings.Replace(h.Filename, " ", "-", -1)

	// save the file
	t, err := os.OpenFile(filepath.Join(filePath, filename), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	//t, err := ioutil.TempFile(filePath, "."+filename)
	if err != nil {
		return filename, nil, err
	}
	defer t.Close()
	// copy the image to the file
	if _, err := io.Copy(t, f); err != nil {
		return filename, nil, err
	}

	return filename, t, nil
}

func ConvertToJpg(imgName string, f *os.File, isThumb bool) (p string, i image.Image, err error) {
	x := filepath.Ext(imgName)
	//n = imgName[:len(imgName)-len(x)]+".jpg"
	s := ""
	n := imgName[:len(imgName)-len(x)]
	if isThumb {
		s = ".thumb"
	}
	p = fmt.Sprintf("%s%s.jpg", n, s)

	if isThumb {
		i, err = imagery.ResizeWidthToJPG(f.Name(), p, true, 200)
	} else {
		i, err = imagery.ConvertToJPG(f.Name(), p, true)
	}

	// delete the temporary image file on error
	if err != nil {
		os.Remove(f.Name())
		f = nil
	}

	return
}
