package network

import (
	"log"
	"net/http"

	s "github.com/nem0z/go-sharding-storage/node/storage"
)

type Network struct {
	Peers []string
}

func (n *Network) Init(port string, s *s.Storage) {
	http.HandleFunc("/file/", HandleFile(s))

	log.Fatal(http.ListenAndServe(port, nil))
}
