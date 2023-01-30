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
func (n *Network) Broadcast(message []byte) []byte {
	ch := make(chan []byte, 1)

	for _, peer := range n.Peers {
		go func(p *Peer, ch chan []byte) {
			resp := p.UDP(message)

			if resp != nil {
				ch <- resp
			}

		}(peer, ch)
	}

	resp := <-ch
	return resp
}
