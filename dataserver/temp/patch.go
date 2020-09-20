package temp

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func patch(w http.ResponseWriter, r *http.Request) {
	uuid := strings.Split(r.URL.EscapedPath(), "/")[2]
	tempinfo, err := readFromFile(uuid)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	infoFile := os.Getenv("STORAGE_ROOT") + "/temp/" + uuid
	datFile := infoFile + ".dat"
	f, err := os.OpenFile(datFile, os.O_WRONLY|os.O_APPEND, 0)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()
	_, err = io.Copy(f, r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	info, err := f.Stat()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	actual := info.Size()
	if actual > tempinfo.Size {
		_ = os.Remove(datFile)
		_ = os.Remove(infoFile)
		log.Println("actual size", actual, "exceeds", tempinfo.Size)
	}
}

func readFromFile(uuid string) (*tempInfo, error) {
	f, err := os.Open(os.Getenv("STORAGE_ROOT") + "/temp/" + uuid)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	bytes, _ := ioutil.ReadAll(f)
	var info tempInfo
	_ = json.Unmarshal(bytes, &info)
	return &info, nil
}
