package main

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"

	"golang.org/x/image/bmp"
	"golang.org/x/image/draw"
)

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
		log.Fatal("unknown image format ")
	}
	if encodeErr != nil {
		log.Fatal(encodeErr)
	}

	return bytes.NewReader(buff.Bytes())
}
