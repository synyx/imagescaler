package main

import (
	"fmt"
	"log"

	"time"

	"github.com/streadway/amqp"
)

func main() {
	config := readConfig()
	connection := connectRabbit(config)
	defer connection.Close()

	channel, err := connection.Channel()
	defer channel.Close()
	failOnError(err, "failed to create channel from connectoin")

	rabbitArtifacts := setupRabbitMqTopicsAndQueues(channel, "user.event", "user.image.event.dev", "user.image.created.#")

	ScaleImage(nil, THUMBNAIL)

	_, deliveryErr := channel.Consume(rabbitArtifacts.queriesQueueName, "what?", false, false, false, false, nil)
	failOnError(deliveryErr, "failed to deliver messages")

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
