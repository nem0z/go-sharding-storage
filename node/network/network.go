package network

import (
	"log"
	"net"
	"net/http"

	s "github.com/nem0z/go-sharding-storage/node/storage"
)

type RequestUDP struct {
	Addr net.Addr
	Data []byte
}

type Network struct {
	PortHTTP string
	PortUDP  string
	PortTCP  string
	Peers    []*Peer
}

func (n *Network) HandleHTTP(s *s.Storage, ch chan []byte) {
	if n.PortHTTP == ":8888" {
		http.HandleFunc("/file/", HandleFile(s, ch))
	}

	err := http.ListenAndServe(n.PortHTTP, nil)
	log.Fatal(err)
}

func (n *Network) HandleUDP(ch_send chan *RequestUDP, ch_recv chan *RequestUDP) {

	listener, err := net.ListenPacket("udp", n.PortUDP)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		defer listener.Close()

		for {
			buf := make([]byte, 65535)
			n, addr, err := listener.ReadFrom(buf)

			if err != nil {
				log.Println("Error handling UDP :", err)
				continue
			}

			req := &RequestUDP{Addr: addr, Data: buf[:n]}
			ch_send <- req
		}
	}()

	go func() {
		for {
			req := <-ch_recv
			_, err = listener.WriteTo(req.Data, req.Addr)
			if err != nil {
				log.Println("Erreur when sending udp resp:", err)
			}
		}
	}()
}

func (n *Network) Relay(message []byte) {
	for i := 0; i < 3; i++ {
		idx := (int(message[len(message)-1]) + i) % len(n.Peers)
		go n.Peers[idx].UDP(message)
	}
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
