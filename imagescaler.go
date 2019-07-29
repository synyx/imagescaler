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
		//ScaleImage(nil, THUMBNAIL)
		log.Printf("got image update %s", imageUpdate.UserUUID)

		imageAsBytes, scaleErr := loadImageFromObjectStorage(imageUpdate.URL)
		if scaleErr != nil {
			log.Printf("failed to load image from image update %v: %v\n", imageUpdate, scaleErr)
		}

		//TODO: optimize me. look at the repeated code
		webReader, scaleErr := scaleImageToTarget(imageAsBytes, WEB)
		webImageUpdate, writeErr := writeImageToObjectStorage(webReader, WEB)
		outgoingImageUpdates <- webImageUpdate

		thumbnailReader, scaleErr := scaleImageToTarget(imageAsBytes, THUMBNAIL)
		thumbnailImageUpdate, writeErr := writeImageToObjectStorage(thumbnailReader, THUMBNAIL)
		outgoingImageUpdates <- thumbnailImageUpdate

		if writeErr != nil {
			log.Printf("failed to write scaled image to object storage: %v ", writeErr)
		}
		if scaleErr != nil {
			log.Printf("failed to scale image to target scale: %v", scaleErr)
		}

	}
}

func scaleImageToTarget(sourceImageBytes []byte, scale Scale) (io.Reader, error) {

	scaledReader, scaleErr := ScaleImage(bytes.NewReader(sourceImageBytes), scale)
	if scaleErr != nil {
		return nil, scaleErr
	}
	return scaledReader, scaleErr

}

func writeImageToObjectStorage(scaledReader io.Reader, scale Scale) (ImageUpdate, error) {
	var imageUpdate ImageUpdate
	return imageUpdate, nil
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
