package client

import (
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// NewClient creates a pakhshi client based on given paho options.
// it uses the servers array for finding out about the clusters and also use
// their host name to name them.
func NewClient(opts *mqtt.ClientOptions) mqtt.Client {
	for _, server := range opts.Servers {
		log.Println(server.Host)
	}

	return nil
}
