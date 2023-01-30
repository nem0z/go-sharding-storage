package main

import (
	"encoding/hex"
	"fmt"
	"log"

	"github.com/nem0z/go-sharding-storage/node/network"
	"github.com/nem0z/go-sharding-storage/node/node"
	"github.com/nem0z/go-sharding-storage/node/utils"
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

	hash_string := "809494ec22ddda2edfd1ff67b8b9125ab39cd4550dc2b90a666b8d44021f56fc"
	hash, err := hex.DecodeString(hash_string)
	Handle(err)

	binary_file, err := main_node.GetFile(hash)
	Handle(err)

	if utils.VerifyFile(hash_string, binary_file) {
		err := utils.ExportFile("./files/example_retrive.png", binary_file)
		Handle(err)
		log.Printf("File successfully retrived : %x\n", hash)
	} else {
		log.Println("Error retrived data doesn't match expected file")
	}

}
