package main

import (
	"log"

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
