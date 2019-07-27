package main

import (
	"log"
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
	go loadAndScaleImage(incomingImageUpdates, outgoingImageUpdates)

	log.Print("hallo")

	<-forever // hammer time!
}

func loadAndScaleImage(incomingImageUpdates <-chan ImageUpdate, outgoingImageUpdates chan<- ImageUpdate) {
	for imageUpdate := range incomingImageUpdates {
		//ScaleImage(nil, THUMBNAIL)
		log.Printf("got image update %s", imageUpdate.UserUUID)
		outgoingImageUpdates <- imageUpdate
	}
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
