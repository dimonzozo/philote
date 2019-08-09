package main

import (
	"net/http"
	"runtime"

	log "github.com/sirupsen/logrus"
)

var Config = LoadConfig()

func main() {
	log.WithFields(log.Fields{
		"port":  Config.port,
		"cores": runtime.NumCPU()}).Info("Initializing Philotic Network")

	log.WithFields(log.Fields{
		"read-buffer-size":  Config.readBufferSize,
		"write-buffer-size": Config.writeBufferSize,
		"max-connections":   Config.maxConnections}).Debug("Configuration options:")

	h := NewHive()
	http.HandleFunc("/", h.ServeNewConnection)

	err := http.ListenAndServe(":"+Config.port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
