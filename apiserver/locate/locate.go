package locate

import (
	"encoding/json"
	"lib/rabbitmq"
	"lib/rs"
	"lib/types"
	"net/http"
	"os"
	"strings"
	"time"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	if method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	info := Locate(strings.Split(r.URL.EscapedPath(), "/")[2])
	if len(info) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	bytes, _ := json.Marshal(info)
	_, _ = w.Write(bytes)
}

func Locate(name string) (locateInfo map[int]string) {
	rabbitMQ := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	rabbitMQ.Publish("dataServers", name)
	channel := rabbitMQ.Consume()
	go func() {
		time.Sleep(time.Second)
		rabbitMQ.Close()
	}()
	locateInfo = make(map[int]string)
	for i := 0; i < rs.ALL_SHARDS; i++ {
		msg := <-channel
		if len(msg.Body) == 0 {
			return
		}
		var message types.LocateMessage
		_ = json.Unmarshal(msg.Body, &message)
		locateInfo[message.Id] = message.Addr
	}
	return
}

func Exist(name string) bool {
	return len(Locate(name)) >= rs.DATA_SHARDS
}
