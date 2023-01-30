package network

import (
	"log"
	"net/http"

	s "github.com/nem0z/go-sharding-storage/node/storage"
)

type Network struct {
	PortHTTP string
	PortUDP  string
	PortTCP  string
	Peers    []*Peer
}

func (n *Network) HandleHTTP(s *s.Storage) {
	if n.PortHTTP == ":8888" {
	http.HandleFunc("/file/", HandleFile(s))
	}

	err := http.ListenAndServe(n.PortHTTP, nil)
	log.Fatal(err)
}
}
