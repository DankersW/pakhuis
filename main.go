package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/DankersW/pakhuis/config"
	"github.com/DankersW/pakhuis/kafka"
	"github.com/DankersW/pakhuis/models"
	log "github.com/sirupsen/logrus"
)

var conf models.Config

func init() {
	conf = config.Get()
	logLevel, err := log.ParseLevel(conf.Log.Level)
	if err != nil {
		log.Errorf("Failed to set minumimum log level, %s", err)
	}
	log.SetLevel(logLevel)
}

func main() {
	log.Info("Starting Pakhuis")

	mainCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	msgChan := make(chan kafka.ConsumerMsg, 10)
	topics := []string{"test", "wsn.sensor-data.telemetry"}

	consumer, err := kafka.NewConsumer(conf.Kafka.Brokers, topics, msgChan)
	if err != nil {
		log.Fatalf("Failed to setup kafka consumer, %s", err.Error())
	}
	go consumer.Serve(mainCtx)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	work(quit, msgChan)

	mainCtx.Done()
}

// TODO: MongoDB driver
