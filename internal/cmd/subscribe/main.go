package subscribe

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/snapp-incubator/pakhshi/pkg/client"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	// TopicFlag used with cobra to read the topic flag from user.
	TopicFlag = "topic"
	// QoSFlag used with cobra to read the qos flag from user.
	QoSFlag = "qos"

	// Timeout for disconnecting from brokers.
	Timeout = 10
)

type Config struct {
	Brokers []string
	QoS     byte
	Topic   string
}

func main(cfg Config) {
	opts := mqtt.NewClientOptions()

	for _, broker := range cfg.Brokers {
		opts.AddBroker(broker)
	}

	c := client.NewClient(opts)

	if token := c.Connect(); token.Wait() && token.Error() != nil {
		pterm.Error.Println("mqtt connection failed:", token.Error().Error())

		return
	}

	pterm.Info.Printf("subscribing on %s\n", cfg.Topic)

	if token := c.Subscribe(cfg.Topic, cfg.QoS, func(c mqtt.Client, m mqtt.Message) {
		pterm.Info.Printf("received: %s\n", string(m.Payload()))
	}); token.Wait() && token.Error() != nil {
		pterm.Error.Println("mqtt subcription failed:", token.Error().Error())

		return
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	c.Disconnect(Timeout)
}

// Register migrate command.
func Register(root *cobra.Command, brokers *[]string) {
	// nolint: exhaustivestruct
	cmd :=
		&cobra.Command{
			Use:   "subscribe",
			Short: "subscribe on multiple brokers at the same time",
			Run: func(cmd *cobra.Command, args []string) {
				qos, _ := cmd.Flags().GetInt8(QoSFlag)
				topic, _ := cmd.Flags().GetString(TopicFlag)

				cfg := Config{
					Brokers: *brokers,
					QoS:     byte(qos),
					Topic:   topic,
				}

				main(cfg)
			},
		}
	cmd.Flags().StringP(TopicFlag, "t", "hello", "topic to subscribe on")
	cmd.Flags().Int8P(QoSFlag, "q", 0, "subscription qos")
	root.AddCommand(cmd)
}
