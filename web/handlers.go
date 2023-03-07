package web

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/nem0z/go-sharding-storage/types"
)

func HandleFile(ch_post chan *types.WrappedChan, ch_get chan *types.WrappedChan) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			PostFile(ch_post, w, r)
		case http.MethodGet:
			GetFile(ch_get, w, r)
		default:
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	}
}

func PostFile(ch chan *types.WrappedChan, w http.ResponseWriter, r *http.Request) {
	binary_file, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error when reading the binary file", http.StatusNotAcceptable)
		return
	}

	wrapped_ch := &types.WrappedChan{Data: binary_file, Chan: make(chan []byte)}
	ch <- wrapped_ch

	file_hash := <-wrapped_ch.Chan
	if file_hash != nil {
		w.Header().Set("Content-Type", "application/json")

		fmt.Fprintf(w, "%x", file_hash)
		if err != nil {
			http.Error(w, "Error writing response body", http.StatusInternalServerError)
			return
		}
	}
}

func GetFile(ch chan *types.WrappedChan, w http.ResponseWriter, r *http.Request) {

	parts := strings.Split(r.URL.Path, "/")

	if len(parts) < 3 {
		http.Error(w, "No hash found", http.StatusBadRequest)
		return
	}

	hash, err := hex.DecodeString(parts[2])
	if err != nil {
		http.Error(w, "Can't convert the hash to buffer representation", http.StatusBadRequest)
		return
	}

	wrapped_ch := &types.WrappedChan{Data: hash, Chan: make(chan []byte)}
	ch <- wrapped_ch

	binary_file := <-wrapped_ch.Chan

	if binary_file == nil {
		http.Error(w, "Error fetching the file for given hash", http.StatusBadRequest)
		return
	}

	contentType := http.DetectContentType(binary_file)
	w.Header().Set("Content-Type", contentType)

	w.Header().Set("Content-Disposition", "attachment; filename=file")
	w.Write(binary_file)
}
