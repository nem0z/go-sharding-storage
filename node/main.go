package main

import (
	"github.com/nem0z/go-sharding-storage/node/node"
)

func main() {
	node := &node.Node{}
	node.Init("./data", ":8888")
}
