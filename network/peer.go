package network

import (
	"net"
	"time"
)

type Peer struct {
	Address string
}

func (p *Peer) Ping() bool {
	conn, err := net.DialTimeout("tcp", p.Address, time.Second*5)
	if err != nil {
		return false
	}

	defer conn.Close()

	conn.Write([]byte("ping"))

	res := make([]byte, 4)
	conn.Read(res)

	return string(res) == "pong"
}

func (p *Peer) UDP(message []byte) []byte {
	conn, err := net.Dial("udp", p.Address)
	if err != nil {
		return nil
	}

	defer conn.Close()

	conn.Write([]byte(message))

	response := make([]byte, 65535)
	conn.SetReadDeadline(time.Now().Add(time.Second * 5))

	n, err := conn.Read(response)
	if err != nil {
		return nil
	}

	return response[:n]
}
