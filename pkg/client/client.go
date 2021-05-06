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

	return &Client{}
}

// Client handles an array for clients to the available cluster
// by using this client you can control all availabe client at the same time.
// you also can use each client separately by their cluster name.
type Client struct {
	Clients map[string]mqtt.Client
}

// IsConnected returns a bool signifying whether
// the client is connected or not.
func (*Client) IsConnected() bool {
	return false
}

// IsConnectionOpen return a bool signifying whether the client has an active
// connection to mqtt broker, i.e not in disconnected or reconnect mode
func (*Client) IsConnectionOpen() bool {
	return false
}

// Connect will create a connection to the message broker, by default
// it will attempt to connect at v3.1.1 and auto retry at v3.1 if that
// fails
func (*Client) Connect() mqtt.Token {
	return nil
}

// Disconnect will end the connection with the server, but not before waiting
// the specified number of milliseconds to wait for existing work to be
// completed.
func (*Client) Disconnect(quiesce uint) {
}

// Publish will publish a message with the specified QoS and content
// to the specified topic.
// Returns a token to track delivery of the message to the broker
func (*Client) Publish(topic string, qos byte, retained bool, payload interface{}) mqtt.Token {
	return nil
}

// Subscribe starts a new subscription. Provide a MessageHandler to be executed when
// a message is published on the topic provided, or nil for the default handler.
//
// If options.OrderMatters is true (the default) then callback must not block or
// call functions within this package that may block (e.g. Publish) other than in
// a new go routine.
// callback must be safe for concurrent use by multiple goroutines.
func (*Client) Subscribe(topic string, qos byte, callback mqtt.MessageHandler) mqtt.Token {
	return nil
}

// SubscribeMultiple starts a new subscription for multiple topics. Provide a MessageHandler to
// be executed when a message is published on one of the topics provided, or nil for the
// default handler.
//
// If options.OrderMatters is true (the default) then callback must not block or
// call functions within this package that may block (e.g. Publish) other than in
// a new go routine.
// callback must be safe for concurrent use by multiple goroutines.
func (*Client) SubscribeMultiple(filters map[string]byte, callback mqtt.MessageHandler) mqtt.Token {
	return nil
}

// Unsubscribe will end the subscription from each of the topics provided.
// Messages published to those topics from other clients will no longer be
// received.
func (*Client) Unsubscribe(topics ...string) mqtt.Token {
	return nil
}

// AddRoute allows you to add a handler for messages on a specific topic
// without making a subscription. For example having a different handler
// for parts of a wildcard subscription or for receiving retained messages
// upon connection (before Sub scribe can be processed).
//
// If options.OrderMatters is true (the default) then callback must not block or
// call functions within this package that may block (e.g. Publish) other than in
// a new go routine.
// callback must be safe for concurrent use by multiple goroutines.
func (*Client) AddRoute(topic string, callback mqtt.MessageHandler) {
}

// OptionsReader returns a ClientOptionsReader which is a copy of the clientoptions
// in use by the client.
func (*Client) OptionsReader() mqtt.ClientOptionsReader {
	return nil
}
