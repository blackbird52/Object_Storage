package objects

import (
	"compress/gzip"
	"io"
	"log"
	"os"
)

func sendFile(w io.Writer, file string) {
	f, err := os.Open(file)
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()

	gzipStream, err := gzip.NewReader(f)
	if err != nil {
		log.Println(err)
		return
	}
	_, _ = io.Copy(w, gzipStream)
	_ = gzipStream.Close()
}
