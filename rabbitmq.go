package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

type rabbitArtifacts struct {
	userEventExchangeName    string
	userImageUpdateQueueName string
}

func setupRabbitMqTopicsAndQueues(channel *amqp.Channel, userEventExchangeName string, userImageEventQueueName string, userImageEventUpdateRoutingKey string) rabbitArtifacts {
	exchangeErr := channel.ExchangeDeclare(userEventExchangeName, "topic", true, false, false, false, nil)
	failOnError(exchangeErr, "failed to declare exchange")

	_, queueDeclarationErr := channel.QueueDeclare(
		userImageEventQueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(queueDeclarationErr, "Failed to declare queue")

	bindErr := channel.QueueBind(userImageEventQueueName, userImageEventUpdateRoutingKey, userEventExchangeName, false, nil)
	failOnError(bindErr, "Failed to bind queries queue to topic exchange")

	log.Printf("created topics and queues %s, %s", userImageEventQueueName, userEventExchangeName)

	return rabbitArtifacts{userEventExchangeName: userEventExchangeName, userImageUpdateQueueName: userImageEventQueueName}
}

func handleIncomingImageUpdateMessages(inBound <-chan amqp.Delivery, outBound chan<- ImageUpdate) {
	for msg := range inBound {

		var imageUpdate ImageUpdate
		jsonErr := json.Unmarshal(msg.Body, &imageUpdate)

		if jsonErr != nil {
			log.Printf("failed to consume image update message %v\n", jsonErr)
			msg.Nack(false, false) // nack and don't requeue -> good bye!
		} else {
			if imageUpdate.ImageScale == "ORIGINAL" {
				outBound <- imageUpdate
				log.Println("successfully consumed image update message")
				msg.Ack(false)
			} else {
				log.Println("won't consume image update messages with scale other than ORIGINAL")
			}
		}
	}
}

func handleOutgoingImageUpdateMessages(inBound <-chan ImageUpdate, channel *amqp.Channel, config imageScalerConfig) {
	for imageUpdate := range inBound {

		messageBody, err := json.Marshal(imageUpdate)
		if err != nil {
			log.Printf("failed to marshal imageUpdate %v to JSON %v\n", imageUpdate, err)
		}

		sendErr := channel.Publish("user.event", fmt.Sprintf(config.imageUpdateRoutingKey, imageUpdate.UserUUID), false, false,
			amqp.Publishing{
				ContentType: "application/json",
				Body:        messageBody,
			})

		if sendErr != nil {
			log.Printf("failed to send reply message: %v", sendErr)
			continue
		}
		log.Printf("published image update event for %v", imageUpdate)
	}
}

func connectRabbit(conf imageScalerConfig) *amqp.Connection {
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
