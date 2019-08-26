package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"log"

	"golang.org/x/image/bmp"
	"golang.org/x/image/draw"
)

// Scale defines a symbolic value for the target size of a scaling operation
type Scale int

func stringToScale(input string) (Scale, error) {
	switch input {
	case "WEB":
		return WEB, nil
	case "THUMBNAIL":
		return THUMBNAIL, nil
	case "ORIGINAL":
		return ORIGINAL, nil
	default:
		return ORIGINAL, errors.New("unknown scale type")
	}
}

func scaleToString(scale Scale) (string, error) {
	switch scale {
	case WEB:
		return "WEB", nil
	case THUMBNAIL:
		return "THUMBNAIL", nil
	case ORIGINAL:
		return "ORIGINAL", nil
	default:
		return "", errors.New("unknown scale type")
	}
}

const (
	// THUMBNAIL is the size for thumbnails
	THUMBNAIL Scale = iota
	// WEB is the size for web usage
	WEB
	// ORIGINAL is the original upload size
	ORIGINAL
)

// ScaleImage converts an incoming image provided by Reader to a scaled version provided by the returned reader
func ScaleImage(in io.Reader, scale Scale) (io.Reader, int, string, error) {
	src, contentType, err := image.Decode(in)

	if err != nil {
		log.Fatal(err)
	}

	dstBounds, err := computeDstBounds(src.Bounds(), scale)
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
func computeDstBounds(srcBounds image.Rectangle, scale Scale) (image.Rectangle, error) {

	var dstBounds image.Rectangle

	var dstX int
	switch scale {
	case THUMBNAIL:
		dstX = 100
		break
	case WEB:
		dstX = 1000
		break
	default:
		return dstBounds, fmt.Errorf("unknown scale: %d", scale)
	}

	if dstX >= srcBounds.Dx() {
		return srcBounds, nil //nothing to do. we do not up-scale atm
	}

	//this not rounding the dimensions but cutting of fraction
	//digits. good enough for me. ;)
	scaleFactor := float64(dstX) / float64(srcBounds.Dx())
	dstY := int(float64(srcBounds.Dy()) * scaleFactor)

	return image.Rect(0, 0, dstX, dstY), nil
}
