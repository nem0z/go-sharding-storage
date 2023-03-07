package main

import (
	"fmt"
	"log"

	"github.com/nem0z/go-sharding-storage/network"
	"github.com/nem0z/go-sharding-storage/node"
)

func Handle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	main_node := &node.Node{}
	main_node.Init("./data", ":8888", ":9100", ":9200")
	main_node.Start()

	network_nodes := make([]*node.Node, 5)
	for i := range network_nodes {
		var node node.Node
		port_http := fmt.Sprintf(":%v", 9001+i)
		port_udp := fmt.Sprintf(":%v", 9101+i)
		port_tcp := fmt.Sprintf(":%v", 9201+i)
		data_dir := fmt.Sprintf("./data/%v", i)

		network_nodes[i] = &node

		node.Init(data_dir, port_http, port_udp, port_tcp)
		node.Start()

		main_node.Network.Peers = append(
			main_node.Network.Peers,
			&network.Peer{Address: fmt.Sprintf("localhost%v", port_udp)},
		)
	}

	select {}

}
