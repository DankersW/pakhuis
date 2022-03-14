package main

import (
	"os"

	"github.com/DankersW/pakhuis/kafka"
	log "github.com/sirupsen/logrus"
)

func work(shutdown chan os.Signal, msgPipe kafka.MsgChan) {
	for {
		select {
		case <-shutdown:
			return
		case msg := <-msgPipe:
			log.Infof("handling msg on topic %q", msg.Topic)
		}
	}
}
