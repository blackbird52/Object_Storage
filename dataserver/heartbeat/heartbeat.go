package heartbeat

import (
	"lib/rabbitmq"
	"os"
	"time"
)

func StartHeartbeat() {
	rabbitMQ := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer rabbitMQ.Close()

	for {
		rabbitMQ.Publish("apiServers", os.Getenv("LISTEN_ADDRESS"))
		time.Sleep(5 * time.Second)
	}
}
