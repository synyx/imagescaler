package main

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
	"log"
	"time"
)

type rabbitConf struct {
	hostname              string
	port                  int
	username              string
	password              string
	filename              string
	imageExchange         string
	imageUpdateQueue      string
	imageUpdateRoutingKey string
	timeout               time.Duration
	minioUrl              string
	minioAccessKey        string
	minioSecret           string
	minioBucketName       string
}

func main() {
	rabbitConfig := readRabbitConf()
	conn := connectRabbit(rabbitConfig)
	defer conn.Close()

	log.Print("hallo")
}
func readRabbitConf() rabbitConf {
	viper.SetConfigFile("config.properties")
	viper.SetConfigType("properties")

	//default values suitable for vanilla rabbitmq docker container
	viper.SetDefault("rabbitmq.hostname", "localhost")
	viper.SetDefault("rabbitmq.port", "5672")
	viper.SetDefault("rabbitmq.username", "guest")
	viper.SetDefault("rabbitmq.password", "guest")
	viper.SetDefault("rabbitmq.timeout", "5s")
	viper.SetDefault("rabbitmq.image.exchange", "image.event")
	viper.SetDefault("minio.url", "http://localhost:9000")
	viper.SetDefault("minio.accesskey", "admin")
	viper.SetDefault("minio.secret", "secret")
	viper.SetDefault("minio.bucketname", "testbucket")

	//load config
	confErr := viper.ReadInConfig()
	logOnError(confErr, "No configuration file loaded - using defaults")

	return rabbitConf{
		hostname:              viper.GetString("rabbitmq.hostname"),
		port:                  viper.GetInt("rabbitmq.port"),
		username:              viper.GetString("rabbitmq.username"),
		password:              viper.GetString("rabbitmq.password"),
		timeout:               viper.GetDuration("rabbitmq.timeout"),
		filename:              viper.GetString("filename"),
		imageExchange:         viper.GetString("rabbitmq.image.exchange"),
		imageUpdateQueue:      viper.GetString("rabbitmq.image.udpate.queue"),
		imageUpdateRoutingKey: viper.GetString("rabbitmq.image.update.routingkey"),
		minioUrl:              viper.GetString("minio.url"),
		minioAccessKey:        viper.GetString("minio.accesskey"),
		minioSecret:           viper.GetString("minio.secret"),
		minioBucketName:       viper.GetString("minio.bucketname"),
	}
}

func connectRabbit(conf rabbitConf) *amqp.Connection {
	for {
		conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/", conf.username, conf.password, conf.hostname, conf.port))
		if err == nil && conn != nil {
			log.Println("connected to rabbitmq")
			return conn
		} else {
			log.Println(fmt.Sprintf("failed to connect to rabbitmq will retry in %d. current cause: %s", conf.timeout, err))
			time.Sleep(conf.timeout)
		}
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
