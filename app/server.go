package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	header := []byte{}
	var msgSize uint32 = 0
	msgSizeBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(msgSizeBytes, msgSize)

	var cID uint32 = 7
	cIDbytes := make([]byte, 4)
	binary.BigEndian.PutUint32(cIDbytes, cID)

	header = append(header, msgSizeBytes...)
	header = append(header, cIDbytes...)

	conn.Write(header)
}

func main() {

	l, err := net.Listen("tcp", "0.0.0.0:9092")
	if err != nil {
		fmt.Println("Failed to bind to port 9092")
		os.Exit(1)
	}
	defer l.Close()

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	handleConnection(conn)
}
