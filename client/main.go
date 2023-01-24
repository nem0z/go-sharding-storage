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

	// Upload the file
	file, err := os.Open(config.filePath)
	Handle(err)

	hash, err := client.Upload(file)
	if err != nil {
		log.Fatal("Error when uploading the file :", err)
	}

	fmt.Println("File successfully uploaded :", hash)

	// Retrive the file
	binary_file, err := client.Get(hash)
	Handle(err)

	file, err = os.Create("./files/example_retrive.png")
	Handle(err)

	_, err = file.Write(binary_file)
	Handle(err)

	fmt.Println("File successfully retrived :", hash)

}
