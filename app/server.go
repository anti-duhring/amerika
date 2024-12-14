package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
)

type Request struct {
	Size          uint32
	CorrelationID uint32
}

type Response struct {
	Size          uint32
	CorrelationID uint32
}

func handleRequest(conn net.Conn) (*Request, error) {
	buff := make([]byte, 1024)
	_, err := conn.Read(buff)
	if err != nil {
		return nil, fmt.Errorf("Error reading from header: %w", err)
	}

	msgSize := binary.BigEndian.Uint32(buff[0:4])
	cID := binary.BigEndian.Uint32(buff[8:12])

	return &Request{
		Size:          msgSize,
		CorrelationID: cID,
	}, nil
}

func generateResponse(conn net.Conn, res Response) {

	resBytes := []byte{}

	res.Size = 0
	msgSizeBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(msgSizeBytes, res.Size)

	cIDbytes := make([]byte, 4)
	binary.BigEndian.PutUint32(cIDbytes, res.CorrelationID)

	resBytes = append(resBytes, msgSizeBytes...)
	resBytes = append(resBytes, cIDbytes...)

	conn.Write(resBytes)
}

func handleConnection(conn net.Conn) error {
	defer conn.Close()

	req, err := handleRequest(conn)
	if err != nil {
		return fmt.Errorf("Error on handleRequest: ", err)
	}

	generateResponse(conn, Response{
		CorrelationID: req.CorrelationID,
	})

	return nil
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

	err = handleConnection(conn)
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
}
