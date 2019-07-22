package main

import (
	"log"

	"github.com/streadway/amqp"
)

type rabbitArtifacts struct {
	queriesExchangeName string
	queriesQueueName    string
}

func setupRabbitMqTopicsAndQueues(channel *amqp.Channel, userImageEventExchangeName string, userImageEventQueueName string, userImageEventUpdateRoutingKey string) rabbitArtifacts {
	exchangeErr := channel.ExchangeDeclare(userImageEventExchangeName, "topic", true, false, false, false, nil)
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

	bindErr := channel.QueueBind(userImageEventQueueName, userImageEventUpdateRoutingKey, userImageEventExchangeName, false, nil)
	failOnError(bindErr, "Failed to bind queries queue to topic exchange")

	log.Printf("created topics and queues %s, %s", userImageEventQueueName, userImageEventExchangeName)

	return rabbitArtifacts{queriesExchangeName: userImageEventExchangeName, queriesQueueName: userImageEventQueueName}
}
