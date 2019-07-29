package main

import (
	"io"
	"log"

	"github.com/minio/minio-go/v6"
)

func writeImageToObjectStorage(scaledReader io.Reader, length int, contentType string, scale Scale, imageUpdate *ImageUpdate, config imageScalerConfig) error {
	var minioOpts minio.PutObjectOptions
	minioClient, err := minio.New(config.minioURL, config.minioAccessKey, config.minioSecret, true)
	if err != nil {
		log.Printf("error while creating min.io client: %v", err)
		return err
	}
	minioClient.PutObject(config.minioBucketName, "name", scaledReader, -1, minioOpts)
	return nil
}
