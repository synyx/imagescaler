package main

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestCorrectMarshaling(t *testing.T) {
	testJSONString := `
		{
			"ImageUUID": "yolo",
			"userUUID": "user",
			"url": "url",
			"imageScale" : "ORIGINAL"
		}	
	`

	var imageUpdate ImageUpdate
	err := json.Unmarshal([]byte(testJSONString), &imageUpdate)
	if err != nil {
		t.Errorf("failed to unmarshal json string: %v", err)
	}
	if imageUpdate.URL != "url" {
		t.Errorf("url did not have expected value 'url': %v", imageUpdate.URL)
	}
	if imageUpdate.ImageUUID != "yolo" {
		t.Errorf("imageUUID did not have expected value 'yolo': %v", imageUpdate.ImageUUID)
	}
	if imageUpdate.UserUUID != "user" {
		t.Errorf("userUUID did not have expected value 'user': %v", imageUpdate.UserUUID)
	}

}

func TestCorrectUnmarshalling(t *testing.T) {
	testImageUpdate := ImageUpdate{ImageUUID: "yolo", UserUUID: "user", URL: "url", ImageScale: "ORIGINAL"}

	byteArr, err := json.Marshal(testImageUpdate)
	if err != nil {
		t.Errorf("failed to marshal imageUpdate to JSON string: %v", err)
	}

	testJSONString := `{"imageUUID":"yolo","userUUID":"user","url":"url","imageScale":"ORIGINAL"}`

	actualOutput := strings.ToLower(string(byteArr))
	expectedOutput := strings.ToLower(testJSONString)
	if actualOutput != expectedOutput {
		t.Errorf("output does not equal expected output (case insensitive): %s vs %s", actualOutput, expectedOutput)
	}
}
