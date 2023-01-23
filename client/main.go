package main

import (
	"fmt"
	"log"
	"os"

	c "github.com/nem0z/go-sharding-storage/client/client"
)

type Conf struct {
	filePath   string
	nodeAdress string
	nodePort   string
}

var config = &Conf{
	filePath:   "./files/example.png",
	nodeAdress: "localhost",
	nodePort:   "8888",
}

func Handle(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	client := &c.Client{}
	client.Init(config.nodeAdress, config.nodePort)

	file, err := os.Open(config.filePath)
	Handle(err)

	hash, err := client.Upload(file)
	if err != nil {
		log.Fatal("Error when uploading the file :", err)
	}

	fmt.Println("File successfully uploaded :", hash)
}
