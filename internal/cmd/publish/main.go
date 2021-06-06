package publish

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/1995parham/pakhshi/pkg/client"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	TopicFlag    = "topic"
	QoSFlag      = "qos"
	RetainedFlag = "retained"
	PayloadFlag  = "payload"

	Timeout = 10
)

type Config struct {
	QoS      byte
	Retained bool
	Brokers  []string
	Topic    string
	Payload  string
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

	c.Publish(cfg.Topic, cfg.QoS, cfg.Retained, cfg.Payload)

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
			Use:   "publish",
			Short: "publish on multiple brokers at the same time",
			Run: func(cmd *cobra.Command, args []string) {
				qos, _ := cmd.Flags().GetInt8(QoSFlag)
				topic, _ := cmd.Flags().GetString(TopicFlag)
				retained, _ := cmd.Flags().GetBool(RetainedFlag)
				payload, _ := cmd.Flags().GetString(PayloadFlag)

				cfg := Config{
					Brokers:  *brokers,
					QoS:      byte(qos),
					Retained: retained,
					Payload:  payload,
					Topic:    topic,
				}

				main(cfg)
			},
		}
	cmd.Flags().StringP(TopicFlag, "t", "hello", "topic to publish on")
	cmd.Flags().BoolP(RetainedFlag, "r", false, "")
	cmd.Flags().Int8P(QoSFlag, "q", 0, "publish qos")
	cmd.Flags().StringP(PayloadFlag, "p", "Hello World", "the message")
	root.AddCommand(cmd)
}
