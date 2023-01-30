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
