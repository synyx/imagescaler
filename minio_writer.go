package main

import (
	"fmt"
	"io"
	"log"

	"github.com/minio/minio-go/v6"
	uuid "github.com/nu7hatch/gouuid"
)

func writeImageToObjectStorage(scaledReader io.Reader, length int, imageType string, scale Scale, imageUpdate *ImageUpdate, config imageScalerConfig) error {
	var minioOpts minio.PutObjectOptions
	minioOpts.ContentType = fmt.Sprintf("image/%s", imageType)
	imageUUID, err := uuid.NewV4()
	if err != nil {
		log.Printf("error while creating image UUID: %v", err)
		return err
	}
	minioClient, err := minio.New(config.minioURL, config.minioAccessKey, config.minioSecret, true)
	if err != nil {
		log.Printf("error while creating min.io client: %v", err)
		return err
	}
	minioClient.PutObject(config.minioBucketName, "name", scaledReader, -1, minioOpts)

	imageUpdate.ImageUUID = imageUUID.String()
	imageUpdate.URL = fmt.Sprintf("%s/%s/%s", config.minioURL, config.minioBucketName, imageUUID)
	scaleString, err := scaleToString(scale) // things can go wrong here
	if err != nil {
		return err
	}
	imageUpdate.ImageScale = scaleString

	return nil
}
