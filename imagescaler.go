package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

//ImageUpdate holds the information encoded by a received image update message
type ImageUpdate struct {
	ImageUUID  string
	UserUUID   string
	URL        string
	ImageScale string
}

func main() {
	config := readConfig()
	connection := connectRabbit(config)
	defer connection.Close()

	//used later to keep the process alive
	forever := make(chan bool)

	incomingImageUpdates := make(chan ImageUpdate)
	outgoingImageUpdates := make(chan ImageUpdate)

	channel, err := connection.Channel()
	defer channel.Close()
	failOnError(err, "failed to create channel from connectoin")

	rabbitArtifacts := setupRabbitMqTopicsAndQueues(channel, "user.event", "user.image.event.dev", "user.image.created.#")

	msgs, deliveryErr := channel.Consume(rabbitArtifacts.userImageUpdateQueueName, "what?", false, false, false, false, nil)
	failOnError(deliveryErr, "failed to deliver messages")

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

	go handleIncomingImageUpdateMessages(msgs, incomingImageUpdates)
	go handleOutgoingImageUpdateMessages(outgoingImageUpdates)
	go handleImageUpdates(incomingImageUpdates, outgoingImageUpdates)

	log.Print("hallo")

	<-forever // hammer time!
}

func handleImageUpdates(incomingImageUpdates <-chan ImageUpdate, outgoingImageUpdates chan<- ImageUpdate) {
	for imageUpdate := range incomingImageUpdates {

		outGoingImageUpdateForWeb, webErr := loadScaleAndWriteImage(imageUpdate, WEB)
		if webErr != nil {
			log.Printf("failed to handle image update %v for WEB scale: %v\n", imageUpdate, webErr)
		} else {
			outgoingImageUpdates <- outGoingImageUpdateForWeb
		}

		outGoingImageUpdateForThumbnail, thumbnailErr := loadScaleAndWriteImage(imageUpdate, THUMBNAIL)
		if thumbnailErr != nil {
			log.Printf("failed to handle image update %v for WEB scale: %v\n", imageUpdate, thumbnailErr)
		} else {
			outgoingImageUpdates <- outGoingImageUpdateForThumbnail
		}

		log.Printf("got image update %s", imageUpdate.UserUUID)
	}
}

func loadScaleAndWriteImage(incomingImageUpdate ImageUpdate, targetScale Scale) (ImageUpdate, error) {

	var imageUpdate ImageUpdate

	imageAsBytes, loadErr := loadImageFromObjectStorage(incomingImageUpdate.URL)
	if loadErr != nil {
		log.Printf("failed to load image from image update %v: %v\n", incomingImageUpdate, loadErr)
		return imageUpdate, loadErr
	}

	thumbnailReader, scaleErr := scaleImageToTarget(imageAsBytes, THUMBNAIL)
	if scaleErr != nil {
		log.Printf("failed to scale image to target scale: %v", scaleErr)
		return imageUpdate, scaleErr
	}

	writeErr := writeImageToObjectStorage(thumbnailReader, THUMBNAIL, &imageUpdate)
	if writeErr != nil {
		log.Printf("failed to write scaled image to object storage: %v ", writeErr)
		return imageUpdate, writeErr
	}

	return imageUpdate, nil
}

func scaleImageToTarget(sourceImageBytes []byte, scale Scale) (io.Reader, error) {

	scaledReader, scaleErr := ScaleImage(bytes.NewReader(sourceImageBytes), scale)
	if scaleErr != nil {
		return nil, scaleErr
	}
	return scaledReader, scaleErr

}


func loadImageFromObjectStorage(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
func logOnError(err error, msg string) {
	if err != nil {
		log.Printf("%s: %s\n", msg, err)
	}
}
