package network

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	s "github.com/nem0z/go-sharding-storage/node/storage"
	"github.com/nem0z/go-sharding-storage/node/utils"
)

func HandleFile(s *s.Storage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			PostFile(s, w, r)
		case http.MethodGet:
			GetFile(s, w, r)
		default:
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	}
}

func PostFile(s *s.Storage, w http.ResponseWriter, r *http.Request) {

	binary_file, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error when reading the binary file", http.StatusNotAcceptable)
		return
	}

	chunk_size := 100000
	chunks := utils.Chunk(binary_file, chunk_size)

	hashs := make([]string, len(chunks))

	for _, chunk := range chunks {
		hash := sha256.Sum256(chunk)
		hashs = append(hashs, fmt.Sprintf("%x", hash))

		err := s.Put(hash[:], chunk)
		if err != nil {
			http.Error(w, "Error storing file part", http.StatusBadRequest)
			return
		}
	}

	json_hashs, err := json.Marshal(hashs)
	if err != nil {
		http.Error(w, "Error marshalling JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(json_hashs)
	if err != nil {
		http.Error(w, "Error writing response body", http.StatusInternalServerError)
		return
	}

}

func GetFile(s *s.Storage, w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")

	if len(parts) < 3 {
		http.Error(w, "No hash found", http.StatusBadRequest)
		return
	}

	hash := parts[2]
	file_part, err := s.Get([]byte(hash))

	if err != nil {
		http.Error(w, "Error fetching the file for given hash", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")

	_, err = w.Write(file_part)

	if err != nil {
		http.Error(w, "Error writing response body", http.StatusInternalServerError)
		return
	}
}
