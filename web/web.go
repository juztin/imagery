// Copyright 2013 Justin Wilson. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"crypto/rand"
	"image"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"minty.io/imagery"
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

func SaveFile(req *http.Request, resp http.ResponseWriter, path, filename string) (*os.File, error) {
	// create the project path
	if err := os.MkdirAll(path, 0755); err != nil {
		return nil, err
	}

	// open image file for writing
	t, err := os.OpenFile(filepath.Join(path, filename), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return nil, err
	}
	defer t.Close()

	// copy the bytes to the file
	r := http.MaxBytesReader(resp, req.Body, 2<<20)
	defer req.Body.Close()
	if _, err := io.Copy(t, r); err != nil {
		return nil, err
	}

	// return the filename
	return t, nil
}

// TODO set a defined max request size (currently set to 10MB in net/http/request.go)
func SaveFormFile(req *http.Request, resp http.ResponseWriter, path string) (*os.File, error) {
	// grab the file from the request
	f, h, err := req.FormFile("image")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// create the project path
	if err := os.MkdirAll(path, 0755); err != nil {
		return nil, err
	}

	// save the file
	t, err := os.OpenFile(filepath.Join(path, h.Filename), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return nil, err
	}
	defer t.Close()
	// copy the image to the file
	if _, err := io.Copy(t, f); err != nil {
		return nil, err
	}

	return t, nil
}

func SaveImage(req *http.Request, resp http.ResponseWriter, path, name string) (image.Image, error) {
	_, err := SaveFile(req, resp, path, name)
	if err != nil {
		return nil, err
	}
	return imagery.Decode(filepath.Join(path, name))
}

func SaveFormImage(req *http.Request, resp http.ResponseWriter, path, name string) (image.Image, error) {
	_, err := SaveFormFile(req, resp, path)
	if err != nil {
		return nil, err
	}
	return imagery.Decode(filepath.Join(path, name))
}

func ImageType(req *http.Request) imagery.ImgType {
	switch req.Header.Get("Content-Type") {
	default:
		return imagery.IMG_UNKNOWN
	case "image/gif":
		return imagery.IMG_GIF
	case "image/jpeg":
		return imagery.IMG_JPEG
	case "image/png":
		return imagery.IMG_PNG
	}
}
