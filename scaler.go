package main

import (
	"bytes"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"log"

	"golang.org/x/image/bmp"
	"golang.org/x/image/draw"
)

// ScaleImage converts an incoming image provided by Reader to a scaled version provided by the returned reader
func ScaleImage(in io.Reader, targetScale scalingTargetConf) (io.Reader, int, string, error) {
	src, contentType, err := image.Decode(in)

	if err != nil {
		log.Fatal(err)
	}

	dstBounds, err := computeDstBounds(src.Bounds(), targetScale.Width)
	if err != nil {
		log.Fatal(err)
	}

	dst := image.NewRGBA(dstBounds)

	draw.BiLinear.Scale(dst, dstBounds, src, src.Bounds(), draw.Over, nil)

	buff := new(bytes.Buffer)
	var encodeErr error

	if contentType == "jpeg" {
		encodeErr = jpeg.Encode(buff, dst, nil)
	} else if contentType == "png" {
		encodeErr = png.Encode(buff, dst)
	} else if contentType == "bmp" {
		encodeErr = bmp.Encode(buff, dst)
	} else if contentType == "gif" {
		encodeErr = gif.Encode(buff, dst, nil)
	} else {
		log.Printf("unknown image format %s", contentType)
	}
	if encodeErr != nil {
		return nil, -1, "nope", err
	}

	return bytes.NewReader(buff.Bytes()), len(buff.Bytes()), contentType, nil
}

//computeDstBounds returns
func computeDstBounds(srcBounds image.Rectangle, width int) (image.Rectangle, error) {

	dstX := width

	if dstX >= srcBounds.Dx() {
		return srcBounds, nil //nothing to do. we do not up-scale atm
	}

	//this not rounding the dimensions but cutting of fraction
	//digits. good enough for me. ;)
	scaleFactor := float64(dstX) / float64(srcBounds.Dx())
	dstY := int(float64(srcBounds.Dy()) * scaleFactor)

	return image.Rect(0, 0, dstX, dstY), nil
}
