package main

import (
	"image"
	"testing"
)

func TestAbs(t *testing.T) {

	srcImage := image.NewRGBA(image.Rect(0, 0, 200, 200))
	srcRectangle := srcImage.Bounds()

	dstRectangle, err := DstBounds(srcRectangle)
	if err != nil {
		t.Errorf("failed to convert: %s", err)
	}
	if dstRectangle.Dy() != 800 {
		t.Errorf("failed to scale image to exepcted value %d - got %d instead.", dstRectangle.Dy(), 800)
	}
}
