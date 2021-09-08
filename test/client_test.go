package main_test

import (
	"testing"

	"github.com/snapp-incubator/pakhshi/pkg/client"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/stretchr/testify/assert"
)

// nolint: cyclop
func TestMultiSubscriberSinglePublisher(t *testing.T) {
	t.Parallel()

	ch := make(chan string)

	{
		opts := mqtt.NewClientOptions()
		opts.AddBroker("tcp://127.0.0.1:1883")
		opts.AddBroker("tcp://127.0.0.1:1884")

		c := client.NewClient(opts)

		if token := c.Connect(); token.Wait() && token.Error() != nil {
			assert.NoError(t, token.Error())
		}

		if token := c.Subscribe("hello", 0, func(c mqtt.Client, m mqtt.Message) {
			ch <- string(m.Payload())
		}); token.Wait() && token.Error() != nil {
			assert.NoError(t, token.Error())
		}
	}

	{
		opts := mqtt.NewClientOptions()
		opts.AddBroker("tcp://127.0.0.1:1883")

		c := mqtt.NewClient(opts)

		if token := c.Connect(); token.Wait() && token.Error() != nil {
			assert.NoError(t, token.Error())
		}

		msg := "b1"
		if token := c.Publish("hello", 0, false, msg); token.Wait() && token.Error() != nil {
			assert.NoError(t, token.Error())
		}

		assert.Equal(t, msg, <-ch)
	}

	{
		opts := mqtt.NewClientOptions()
		opts.AddBroker("tcp://127.0.0.1:1884")

		c := mqtt.NewClient(opts)

		if token := c.Connect(); token.Wait() && token.Error() != nil {
			assert.NoError(t, token.Error())
		}

		msg := "b2"
		if token := c.Publish("hello", 0, false, msg); token.Wait() && token.Error() != nil {
			assert.NoError(t, token.Error())
		}

		assert.Equal(t, msg, <-ch)
	}
}
