package locate

import (
	"encoding/json"
	"lib/rabbitmq"
	"net/http"
	"os"
	"strconv"
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

func Locate(name string) string {
	rabbitMQ:= rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	rabbitMQ.Publish("dataServers", name)
	channel := rabbitMQ.Consume()
	go func() {
		time.Sleep(time.Second)
		rabbitMQ.Close()
	}()
	msg := <-channel
	server, _ := strconv.Unquote(string(msg.Body))
	return server
}

func Exist(name string) bool {
	return Locate(name) != ""
}
