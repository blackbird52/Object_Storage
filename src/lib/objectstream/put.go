package objectstream

import (
	"fmt"
	"io"
	"net/http"
)

type PutStream struct {
	writer  *io.PipeWriter
	chanErr chan error
}

func NewPutStream(server string, object string) *PutStream {
	reader, writer := io.Pipe()
	chanErr := make(chan error)
	go func() {
		request, _ := http.NewRequest("PUT", "http://"+server+"/objects/"+object, reader)
		client := http.Client{}
		response, err := client.Do(request)
		if err == nil && response.StatusCode != http.StatusOK {
			err = fmt.Errorf("dataServer return http code %d", response.StatusCode)
			chanErr <- err
		}
	}()
	return &PutStream{writer, chanErr}
}

func (w *PutStream) Write(p []byte) (n int, err error) {
	return w.writer.Write(p)
}

func (w *PutStream) Close() error {
	_ = w.writer.Close()
	return <-w.chanErr
}
