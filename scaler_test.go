package main

import (
	"image"
	"testing"
)

func TestCorrectThumbnailScaling(t *testing.T) {

	srcImage := image.NewRGBA(image.Rect(0, 0, 1200, 200))
	srcRectangle := srcImage.Bounds()

	dstRectangle, err := computeDstBounds(srcRectangle, THUMBNAIL)
	if err != nil {
		t.Errorf("failed to convert: %s", err)
	}
	if dstRectangle.Dx() != 100 || dstRectangle.Dy() != 16 {
		t.Errorf("failed to scale image to exepcted value %d - got %d instead.", 100, dstRectangle.Dy())
	}
}

func TestCorrectWebScaling(t *testing.T) {
	srcImage := image.NewRGBA(image.Rect(0, 0, 1200, 200))
	srcRectangle := srcImage.Bounds()

	dstRectangle, err := computeDstBounds(srcRectangle, WEB)
	if err != nil {
		t.Errorf("failed to convert: %s", err)
	}
	if dstRectangle.Dx() != 1000 || dstRectangle.Dy() != 166 {
		t.Errorf("failed to scale image to exepcted value %d - got %d instead.", 100, dstRectangle)
	}
}
