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

func (node *Node) Init(storage_path string, http_port string) {
	storage := &s.Storage{}
	storage.Init(storage_path)
	node.Storage = storage

	network := &n.Network{}
	network.Init(http_port, node.Storage)
	node.Network = network
}
