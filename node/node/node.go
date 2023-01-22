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

func (n *Node) Init(storage_path string) {
	storage := &s.Storage{}
	storage.Init(storage_path)
}
