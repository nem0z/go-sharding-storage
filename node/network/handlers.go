package network

import (
	"crypto/sha256"
	"encoding/hex"
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

	hash_table := make(map[int]string, len(chunks))

	for i, chunk := range chunks {
		hash := sha256.Sum256(chunk)
		hash_table[i] = fmt.Sprintf("%x", hash)

		err := s.Put(hash[:], chunk)
		if err != nil {
			http.Error(w, "Error storing file part", http.StatusBadRequest)
			return
		}
	}

	file_hash := sha256.Sum256(binary_file)
	json_hash_table, err := json.Marshal(hash_table)
	if err != nil {
		http.Error(w, "Error marshalling JSON", http.StatusInternalServerError)
		return
	}

	s.Put(file_hash[:], json_hash_table)

	w.Header().Set("Content-Type", "application/json")

	fmt.Fprintf(w, "%x", file_hash)
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

	hash, err := hex.DecodeString(parts[2])
	if err != nil {
		http.Error(w, "Can't convert the hash to buffer representation", http.StatusBadRequest)
		return
	}

	binary_file, err := s.Get(hash)

	if err != nil {
		http.Error(w, "Error fetching the file for given hash", http.StatusBadRequest)
		return
	}

	var hash_table map[int]string
	err = json.Unmarshal(binary_file, &hash_table)
	if err != nil {
		http.Error(w, "Error parsing the table hash", http.StatusInternalServerError)
		return
	}

	file_data := []byte{}
	for i := 0; i < len(hash_table); i++ {
		hash, err := hex.DecodeString(hash_table[i])
		if err != nil {
			return
		}

		file_part, err := s.Get(hash)
		if err != nil {
			return
		}

		file_data = append(file_data, file_part...)
	}

	w.Header().Set("Content-Type", "application/octet-stream")

	_, err = w.Write(file_data)

	if err != nil {
		http.Error(w, "Error writing response body", http.StatusInternalServerError)
		return
	}
}
