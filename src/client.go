package main

import (
	"encoding/json"
	"fmt"
	"github.com/sitetester/tcp_service/src/builder"
	"github.com/sitetester/tcp_service/src/common"
	"io"
	"log"
	"net"
)

func main() {

	port := ":8081"
	conn, err := net.Dial("tcp", port)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()
	log.Printf("Client connected to server %s\n", port)

	var payloadBuilder builder.PayloadBuilder
	payload := payloadBuilder.BuildPayload()

	log.Printf("Sending JSON payload: %v", payload)

	encoder := json.NewEncoder(conn)
	if err := encoder.Encode(&payload); err != nil {
		log.Println(err)
	}

	trackStatus(conn)
}

func trackStatus(conn net.Conn) {

	notify := make(chan error)

	go func() {
		bytes := make([]byte, 1024)
		for {
			n, err := conn.Read(bytes)
			if err != nil {
				notify <- err
				return
			}

			if n > 0 {
				var onlineStatus common.OnlineStatus
				err := json.Unmarshal(bytes[:n], &onlineStatus)
				if err != nil {
					log.Println(err)
				}

				fmt.Println(onlineStatus)
			}
		}
	}()

	for {
		select {
		case err := <-notify:
			if err == io.EOF {
				fmt.Println("Closed server connection detected.")
				return
			}
		}
	}
}
