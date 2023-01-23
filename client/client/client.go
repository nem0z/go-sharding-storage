package client

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type Node struct {
	adress string
	port   string
}

func (n *Node) url() string {
	return fmt.Sprintf("http://%v:%v/file/", n.adress, n.port)
}

type Client struct {
	Node *Node
}

func (c *Client) Init(node_adress string, node_port string) {
	node := &Node{adress: node_adress, port: node_port}
	c.Node = node
}

func (c *Client) Upload(data io.Reader) (string, error) {
	url := c.Node.url()
	req, err := http.NewRequest("POST", url, data)

	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/octet-stream")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	var hash string
	err = json.Unmarshal(body, &hash)

	return hash, err
}

func (c *Client) Get(hash string) ([]byte, error) {
	url := c.Node.url()
	resp, err := http.Get(url + hash)

	if err != nil {
		return []byte{}, err
	}

	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
