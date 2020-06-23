package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"
)

type Philote struct {
	ID               string
	AccessKey        *AccessKey
	hive             *hive
	ws               *websocket.Conn
	IncomingMessages chan *Message
}

func NewPhilote(ak *AccessKey, ws *websocket.Conn) *Philote {
	p := &Philote{
		ws:               ws,
		ID:               uuid.NewV4().String(),
		AccessKey:        ak,
		IncomingMessages: make(chan *Message),
	}

	go p.DistributeIncomingMessages()

	return p
}

func (p *Philote) SetHive(hive *hive) {
	p.hive = hive
}

func (p *Philote) DistributeIncomingMessages() {
	var message *Message

	for {
		message = <-p.IncomingMessages
		p.ws.WriteJSON(message)
	}
}

func (p *Philote) Listen() {
	log.WithFields(log.Fields{"philote": p.ID}).Debug("Listening to Philote")
	for {
		message := &Message{}
		err := p.ws.ReadJSON(&message)
		if err != nil {
			log.WithFields(log.Fields{
				"philote": p.ID,
				"error":   err.Error()}).Warn("Error reading from socket, disconnecting")

			p.hive.Disconnect(p)
			break
		}

		// Ensure no tampering with message data
		message.IssuerID = p.ID

		log.WithFields(log.Fields{"philote": p.ID, "channel": message.Channel}).Debug("Received message from socket")

		if p.AccessKey.CanWrite(message.Channel) {
			p.publish(message)
		} else {
			log.WithFields(log.Fields{
				"philote": p.ID,
				"channel": message.Channel,
				"data":    message.Data,
			}).Info("Message dropped due to insufficient write permissions")
		}
	}
}

func (p *Philote) disconnect() {
	log.WithFields(log.Fields{"philote": p.ID}).Debug("Closing Philote")
	p.ws.Close()
}

func (p *Philote) publish(message *Message) {
	message.IssuerID = p.ID
	p.hive.Publish(message)
}
