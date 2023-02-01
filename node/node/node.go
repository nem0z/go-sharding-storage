package node

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"

	n "github.com/nem0z/go-sharding-storage/node/network"
	s "github.com/nem0z/go-sharding-storage/node/storage"
	"github.com/nem0z/go-sharding-storage/node/utils"
	"github.com/syndtr/goleveldb/leveldb"
)

type Node struct {
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
	ch_http_parts := make(chan []byte)

	go node.Network.HandleHTTP(node.Storage, ch_http_parts)
	go node.Network.HandleUDP(ch_udp_req, ch_udp_resp)

	go func(ch_req chan *n.RequestUDP, ch_resp chan *n.RequestUDP) {
		for {
			req := <-ch_req
			splitted_data := bytes.SplitN(req.Data, []byte(" "), 2)
			method, data := splitted_data[0], splitted_data[1]

			switch string(method) {
			case "get_part":
				file_part, err := node.Storage.Get(data)

			if err != nil {
					if err != leveldb.ErrNotFound {
				log.Println("Error retriving file part:", err)
					}
				continue
			}

			resp := &n.RequestUDP{Addr: req.Addr, Data: file_part}
			ch_resp <- resp

			case "relay_part":
				hash := sha256.Sum256(data)
				node.Storage.Put(hash[:], data)
			}
		}
	}(ch_udp_req, ch_udp_resp)
func (node *Node) RelayPart(data []byte) {
	message := append([]byte("relay_part "), data...)
	node.Network.Relay(message)
}

func (node *Node) RelayFile(part_table [][]byte) {
	for _, part := range part_table {
		node.RelayPart(part)
	}
}

func (node *Node) GetPart(hash []byte) []byte {
	message := append([]byte("get_part "), hash...)
	return node.Network.Broadcast(message)
}

func (node *Node) GetFile(hash []byte) ([]byte, error) {
	data, err := node.Storage.Get(hash)

	if err != nil {
		return nil, err
	}

	var hash_table map[int]string
	err = json.Unmarshal(data, &hash_table)

	if err != nil {
		return nil, err
	}

	binary_file := []byte{}

	for i := 0; i < len(hash_table); i++ {
		hash, err := hex.DecodeString(hash_table[i])
		if err != nil {
			return nil, err
		}

		file_part, err := node.Storage.Get(hash)

		//if err != nil {
		file_part = node.GetPart(hash)
		//}

		binary_file = append(binary_file, file_part...)
	}

	return binary_file, nil
}
