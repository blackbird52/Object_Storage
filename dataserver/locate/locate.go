package locate

import (
	"lib/rabbitmq"
	"os"
	"strconv"
)

func Locate(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

func StartLocate() {
	rabbitMQ := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer rabbitMQ.Close()

	rabbitMQ.Bind("dataServers")
	channel := rabbitMQ.Consume()
	for msg := range channel {
		object, err := strconv.Unquote(string(msg.Body))
		if err != nil {
			panic(err)
		}
		if Locate(os.Getenv("STORAGE_ROOT") + "/objects/" + object) {
			rabbitMQ.Send(msg.ReplyTo, os.Getenv("LISTEN_ADDRESS"))
		}
	}
}
