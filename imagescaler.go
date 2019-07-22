package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

type ImageUpdate struct {
	imageUUID string
	userUUID  string
}

func main() {
	config := readConfig()
	connection := connectRabbit(config)
	defer connection.Close()

	channel, err := connection.Channel()
	defer channel.Close()
	failOnError(err, "failed to create channel from connectoin")

	rabbitArtifacts := setupRabbitMqTopicsAndQueues(channel, "user.event", "user.image.event.dev", "user.image.created.#")

	ScaleImage(nil, THUMBNAIL)

	msgs, deliveryErr := channel.Consume(rabbitArtifacts.userImageUpdateQueueName, "what?", false, false, false, false, nil)
	failOnError(deliveryErr, "failed to deliver messages")

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	go func() {
		for msg := range msgs {

			var imageUpdate ImageUpdate
			jsonErr := json.Unmarshal(msg.Body, &imageUpdate)

			if jsonErr != nil {
				log.Println("failed to consume image update message")
				msg.Nack(false, false)
			} else {
				log.Println("successfully consumed image update message")
				msg.Ack(false)
			}
		}
	}()

	log.Print("hallo")
}

func connectRabbit(conf rabbitConf) *amqp.Connection {
	for {
		conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/", conf.username, conf.password, conf.hostname, conf.port))
		if err == nil && conn != nil {
			log.Println("connected to rabbitmq")
			return conn
		}
		log.Println(fmt.Sprintf("failed to connect to rabbitmq will retry in %d. current cause: %s", conf.timeout, err))
		time.Sleep(conf.timeout)
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
