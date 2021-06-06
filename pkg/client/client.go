package client

import (
	"net/url"

	"github.com/1995parham/pakhshi/pkg/token"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// NewClient creates a pakhshi client based on given paho options.
// it uses the servers array for finding out about the clusters and also use
// their host name to name them.
func NewClient(opts *mqtt.ClientOptions) mqtt.Client {
	servers := make(map[string]*mqtt.ClientOptions)

	for _, server := range opts.Servers {
		lopts := new(mqtt.ClientOptions)
		*lopts = *opts
		lopts.Servers = []*url.URL{server}
		servers[server.String()] = lopts
	}

	return NewClientWithOptions(servers)
}

// NewClientWithOptions creates a pakhshi client based on given paho options and broker names.
func NewClientWithOptions(opts map[string]*mqtt.ClientOptions) mqtt.Client {
	clients := make(map[string]mqtt.Client)

	for name, opt := range opts {
		clients[name] = mqtt.NewClient(opt)
	}

	return &Client{
		Clients: clients,
	}
}

// Client handles an array for clients to the available cluster
// by using this client you can control all avaialbe client at the same time.
// you also can use each client separately by their cluster name.
type Client struct {
	Clients map[string]mqtt.Client
}

// IsConnected returns a bool signifying whether
// the client is connected to all mqtt brokers or not.
func (c *Client) IsConnected() bool {
	result := true

	for _, client := range c.Clients {
		result = result && client.IsConnected()
	}

	return result
}

// IsConnectionOpen return a bool signifying whether the client has an active
// connection to all mqtt brokers, i.e not in disconnected or reconnect mode.
func (c *Client) IsConnectionOpen() bool {
	result := true

	for _, client := range c.Clients {
		result = result && client.IsConnectionOpen()
	}

	return result
}

// Connect will create a connection to the message broker, by default
// it will attempt to connect at v3.1.1 and auto retry at v3.1 if that
// fails.
func (c *Client) Connect() mqtt.Token {
	token := token.NewTokens()

	for name, client := range c.Clients {
		token.Append(name, client.Connect())
	}

	return token
}

// Disconnect will end the connection with the server, but not before waiting
// the specified number of milliseconds to wait for existing work to be
// completed.
// Please note that this time will be given to each broker.
func (c *Client) Disconnect(quiesce uint) {
	for _, client := range c.Clients {
		client.Disconnect(quiesce)
	}
}

// Publish will publish a message with the specified QoS and content
// to the specified topic.
// Returns a token to track delivery of the message to the broker.
func (c *Client) Publish(topic string, qos byte, retained bool, payload interface{}) mqtt.Token {
	token := token.NewTokens()

	for name, client := range c.Clients {
		token.Append(name, client.Publish(topic, qos, retained, payload))
	}

	return token
}

// Subscribe starts a new subscription. Provide a MessageHandler to be executed when
// a message is published on the topic provided, or nil for the default handler.
//
// If options.OrderMatters is true (the default) then callback must not block or
// call functions within this package that may block (e.g. Publish) other than in
// a new go routine.
// callback must be safe for concurrent use by multiple goroutines.
func (c *Client) Subscribe(topic string, qos byte, callback mqtt.MessageHandler) mqtt.Token {
	token := token.NewTokens()

	for name, client := range c.Clients {
		token.Append(name, client.Subscribe(topic, qos, callback))
	}

	return token
}

// SubscribeMultiple starts a new subscription for multiple topics. Provide a MessageHandler to
// be executed when a message is published on one of the topics provided, or nil for the
// default handler.
//
// If options.OrderMatters is true (the default) then callback must not block or
// call functions within this package that may block (e.g. Publish) other than in
// a new go routine.
// callback must be safe for concurrent use by multiple goroutines.
func (c *Client) SubscribeMultiple(filters map[string]byte, callback mqtt.MessageHandler) mqtt.Token {
	token := token.NewTokens()

	for name, client := range c.Clients {
		token.Append(name, client.SubscribeMultiple(filters, callback))
	}

	return token
}

// Unsubscribe will end the subscription from each of the topics provided.
// Messages published to those topics from other clients will no longer be
// received.
func (c *Client) Unsubscribe(topics ...string) mqtt.Token {
	token := token.NewTokens()

	for name, client := range c.Clients {
		token.Append(name, client.Unsubscribe(topics...))
	}

	return token
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
func (c *Client) AddRoute(topic string, callback mqtt.MessageHandler) {
	for _, client := range c.Clients {
		client.AddRoute(topic, callback)
	}
}

// OptionsReader returns a ClientOptionsReader which is a copy of the clientoptions
// in use by the client.
func (*Client) OptionsReader() mqtt.ClientOptionsReader {
	return mqtt.ClientOptionsReader{}
}
