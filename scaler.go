package main

import (
	"bytes"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"

	"golang.org/x/image/bmp"
	"golang.org/x/image/draw"
)

type base int

const (
	// THUMBNAIL is the size for thumbnails
	THUMBNAIL base = iota
	// WEB is the size for web usage
	WEB
)

// ScaleImage converts an incoming image provided by Reader to a scaled version provided by the returned reader
func ScaleImage(in io.Reader) io.Reader {
	src, imageType, err := image.Decode(in)

	if err != nil {
		log.Fatal(err)
	}
	dst := image.NewRGBA(image.Rect(0, 0, src.Bounds().Dx(), src.Bounds().Dy()))

	draw.BiLinear.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Over, nil)

	buff := new(bytes.Buffer)
	var encodeErr error

	if imageType == "jpeg" {
		encodeErr = jpeg.Encode(buff, dst, nil)
	} else if imageType == "png" {
		encodeErr = png.Encode(buff, dst)
	} else if imageType == "bmp" {
		encodeErr = bmp.Encode(buff, dst)
	} else {
		log.Printf("unknown image format %s", imageType)
	}
	if encodeErr != nil {
		log.Fatal(encodeErr)
	}

	return bytes.NewReader(buff.Bytes())
}

//DstBounds returns
func DstBounds(srcBounds image.Rectangle) (image.Rectangle, error) {

	var dstBounds image.Rectangle

	//do some math

	return dstBounds, errors.New("failed to compute destination bounds")
}
