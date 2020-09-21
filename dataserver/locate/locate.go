package locate

import (
	"lib/rabbitmq"
	"lib/types"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

var objects = make(map[string]int)
var mutex sync.Mutex

func Locate(hash string) int {
	mutex.Lock()
	id, ok := objects[hash]
	mutex.Unlock()
	if !ok {
		return -1
	} else {
		return id
	}
}

func Add(hash string, id int) {
	mutex.Lock()
	objects[hash] = id
	mutex.Unlock()
}

func Del(hash string) {
	mutex.Lock()
	delete(objects, hash)
	mutex.Unlock()
}

func StartLocate() {
	rabbitMQ := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer rabbitMQ.Close()

	rabbitMQ.Bind("dataServers")
	channel := rabbitMQ.Consume()
	for msg := range channel {
		hash, err := strconv.Unquote(string(msg.Body))
		if err != nil {
			panic(err)
		}
		id := Locate(hash)
		if id != -1 {
			rabbitMQ.Send(msg.ReplyTo, types.LocateMessage{Addr: os.Getenv("LISTEN_ADDRESS"), Id: id})
		}
	}
}

func CollectObjects() {
	files, _ := filepath.Glob(os.Getenv("STORAGE_ROOT") + "/objects/*")
	for _, f := range files {
		file := strings.Split(filepath.Base(f), ".")
		if len(file) != 3 {
			panic(f)
		}
		hash := file[0]
		id, err := strconv.Atoi(file[1])
		if err != nil {
			panic(err)
		}
		objects[hash] = id
	}
}
