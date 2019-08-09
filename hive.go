package main

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"sync"
)

type hive struct {
	m          *sync.Mutex
	philotes   map[string]*Philote
	connect    chan *Philote
	disconnect chan *Philote
}

func NewHive() *hive {
	h := &hive{
		m:          &sync.Mutex{},
		philotes:   map[string]*Philote{},
		connect:    make(chan *Philote),
		disconnect: make(chan *Philote),
	}

	go h.MaintainPhiloteIndex()

	return h
}

func (h *hive) PhilotesCount() int {
	h.m.Lock()
	defer h.m.Unlock()
	return len(h.philotes)
}

func (h *hive) Disconnect(philote *Philote) {
	h.disconnect <- philote
}

func (h *hive) Publish(message *Message) {
	h.m.Lock()
	defer h.m.Unlock()
	for _, philote := range h.philotes {
		if message.IssuerID == philote.ID {
			continue
		}

		for _, channel := range philote.AccessKey.Read {
			if message.Channel == channel {
				philote.IncomingMessages <- message
				break
			}
		}

	}
}

func (h *hive) MaintainPhiloteIndex() {
	log.Debug("Starting bookeeper")

	for {
		select {
		case p := <-h.connect:
			if h.PhilotesCount() >= Config.maxConnections {
				log.WithFields(log.Fields{"philote": p.ID}).Warn("MAX_CONNECTIONS limit reached, dropping new connection")
				p.disconnect()
			}

			log.WithFields(log.Fields{"philote": p.ID}).Debug("Registering Philote")
			h.m.Lock()
			p.SetHive(h)
			h.philotes[p.ID] = p
			h.m.Unlock()
			go p.Listen()
		case p := <-h.disconnect:
			log.WithFields(log.Fields{"philote": p.ID}).Debug("Disconnecting Philote")
			h.m.Lock()
			delete(h.philotes, p.ID)
			h.m.Unlock()
			p.disconnect()
		}
	}
}

func (h *hive) ServeNewConnection(w http.ResponseWriter, r *http.Request) {
	auth := strings.TrimSpace(strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer"))
	if auth == "" {
		r.ParseForm()
		auth = r.Form.Get("auth")
		log.WithFields(log.Fields{"auth": auth}).Debug("Empty Authorization header, trying querystring #auth param")
	}

	accessKey, err := NewAccessKey(auth)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	connection, err := Config.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.WithFields(log.Fields{"error": err.Error()}).Warn("Can't upgrade connection")
		w.Write([]byte(err.Error()))
		return
	}

	philote := NewPhilote(accessKey, connection)
	h.connect <- philote
}
