package objects

import (
	"io"
	"os"
)

func sendFile(w io.Writer, file string) {
	f, _ := os.Open(file)
	defer f.Close()
	_, _ = io.Copy(w, f)
}
