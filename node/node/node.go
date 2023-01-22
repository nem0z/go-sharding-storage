package node

import (
	n "github.com/nem0z/go-sharding-storage/node/network"
	s "github.com/nem0z/go-sharding-storage/node/storage"
)

type Node struct {
	Adresse string
	Port    string
	Storage s.Storage
	Network n.Network
}
