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
	ImageUUID  string `json:"imageUUID"`
	UserUUID   string `json:"userUUID"`
	URL        string `json:"url"`
	ImageScale string `json:"imageScale"`
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
	failOnError(err, "failed to create channel from connection")

	rabbitArtifacts := setupRabbitMqTopicsAndQueues(channel, config.imageExchange, config.imageUpdateQueue, config.imageUpdateRoutingKey)

	msgs, deliveryErr := channel.Consume(rabbitArtifacts.userImageUpdateQueueName, "what?", false, false, false, false, nil)
	failOnError(deliveryErr, "failed to deliver messages")

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

	go handleIncomingImageUpdateMessages(msgs, incomingImageUpdates, config)
	go handleOutgoingImageUpdateMessages(outgoingImageUpdates, channel, config)
	go handleImageUpdates(incomingImageUpdates, outgoingImageUpdates, config)

	<-forever // hammer time!
}

func handleImageUpdates(incomingImageUpdates <-chan ImageUpdate, outgoingImageUpdates chan<- ImageUpdate, config imageScalerConfig) {
	for imageUpdate := range incomingImageUpdates {

		for _, targetScale := range config.scalingTarget {
			outGoingImageUpdate, scaleErr := loadScaleAndWriteImage(imageUpdate, targetScale, config)
			if scaleErr != nil {
				log.Printf("failed to handle image update %v for %s scale: %v\n", imageUpdate, targetScale.Factor, scaleErr)
			} else {
				outgoingImageUpdates <- outGoingImageUpdate
			}
			log.Printf("handled image update image update %v", imageUpdate)
		}

	}
}

func loadScaleAndWriteImage(incomingImageUpdate ImageUpdate, targetScale scalingTargetConf, config imageScalerConfig) (ImageUpdate, error) {

	var imageUpdate ImageUpdate

	imageAsBytes, loadErr := loadImageFromObjectStorage(incomingImageUpdate.URL)
	if loadErr != nil {
		log.Printf("failed to load image from image update %v: %v\n", incomingImageUpdate, loadErr)
		return imageUpdate, loadErr
	}

	thumbnailReader, scaledLength, contentType, scaleErr := scaleImageToTarget(imageAsBytes, targetScale)
	if scaleErr != nil {
		log.Printf("failed to scale image to target scale %v: %v", targetScale, scaleErr)
		return imageUpdate, scaleErr
	}

	imageUpdate, writeErr := writeImageToObjectStorage(thumbnailReader, scaledLength, contentType, targetScale.Factor, config)
	if writeErr != nil {
		log.Printf("failed to write scaled image to object storage: %v ", writeErr)
		return imageUpdate, writeErr
	}
	imageUpdate.UserUUID = incomingImageUpdate.UserUUID // don't forget the userUUID

	log.Printf("wrote new image to min.io: %v", imageUpdate)
	return imageUpdate, nil
}

func scaleImageToTarget(sourceImageBytes []byte, targetScale scalingTargetConf) (io.Reader, int, string, error) {

	scaledReader, scaledLength, contentType, scaleErr := ScaleImage(bytes.NewReader(sourceImageBytes), targetScale)
	if scaleErr != nil {
		return nil, -1, "nope", scaleErr
	}
	return scaledReader, scaledLength, contentType, scaleErr

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
