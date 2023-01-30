package node

import (
	n "github.com/nem0z/go-sharding-storage/node/network"
	s "github.com/nem0z/go-sharding-storage/node/storage"
)

type Node struct {
	Adresse string
	Port    string
	Storage *s.Storage
	Network *n.Network
}

func (node *Node) Init(storage_path string, http_port string, udp_port string, tcp_port string) {

	storage := &s.Storage{}
	storage.Init(storage_path)
	node.Storage = storage

	network := &n.Network{
		PortHTTP: http_port,
		PortUDP:  udp_port,
		PortTCP:  tcp_port,
	}

	node.Network = network
}

func (node *Node) Start() {
	ch_udp_req := make(chan *n.RequestUDP)
	ch_udp_resp := make(chan *n.RequestUDP)

	go node.Network.HandleHTTP(node.Storage)
	go node.Network.HandleUDP(ch_udp_req, ch_udp_resp)

	go func(ch_req chan *n.RequestUDP, ch_resp chan *n.RequestUDP) {
		for {
			req := <-ch_req
			splitted_data := bytes.SplitN(req.Data, []byte(" "), 2)
			_, hash := splitted_data[0], splitted_data[1]

			file_part, err := node.Storage.Get(hash)

			if err != nil {
				log.Println("Error retriving file part:", err)
				continue
			}

			resp := &n.RequestUDP{Addr: req.Addr, Data: file_part}
			ch_resp <- resp
		}
	}(ch_udp_req, ch_udp_resp)
}
