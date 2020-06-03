package main

import (
	"encoding/json"
	"fmt"
	"github.com/sitetester/tcp_service/src/common"
	"io"
	"log"
	"net"
	"strconv"
)

func handleConnection(conn net.Conn, onlineClientsChan chan common.TCPClient, offlineClientsChan chan common.Payload) {

	buf := make([]byte, 1024)
	var payload common.Payload

	for {
		n, conErr := conn.Read(buf)

		if n > 0 {
			err := json.Unmarshal(buf[:n], &payload)
			if err != nil {
				log.Println(err)
			}

			client := common.TCPClient{
				Payload: payload,
				Conn:    conn,
			}

			onlineClientsChan <- client
		}

		if conErr != nil && conErr == io.EOF {
			offlineClientsChan <- payload
			return
		}
	}
}

func isFriend(userId int, friends []int) bool {

	for _, friendId := range friends {
		if userId == friendId {
			return true
		}
	}

	return false
}

func trackOnlineStatus(onlineClients map[int]common.TCPClient, onlineClientsChan chan common.TCPClient, offlineClientsChan chan common.Payload) {

	for {
		select {
		case client := <-onlineClientsChan:
			onlineClients[client.Payload.UserId] = client

			fmt.Println()
			log.Printf("Client joined with address: %s & payload: %v", client.Conn.RemoteAddr().String(), client.Payload)
			fmt.Println("Total online clients: " + strconv.Itoa(len(onlineClients)))

			for _, friend := range onlineClients {
				if isFriend(client.Payload.UserId, friend.Payload.Friends) {
					log.Printf("Sending online notification to friend (userId): %d\n", friend.Payload.UserId)
					encoder := json.NewEncoder(friend.Conn)
					onlineStatus := common.OnlineStatus{
						UserId: client.Payload.UserId,
						Online: true,
					}

					if err := encoder.Encode(&onlineStatus); err != nil {
						log.Println(err)
					}
				}
			}

		case offlineUser := <-offlineClientsChan:
			log.Printf("userId:%d is now offline\n", offlineUser.UserId)
			delete(onlineClients, offlineUser.UserId)

			log.Printf("Total online clients: %d", len(onlineClients))

			for _, client := range onlineClients {
				if isFriend(client.Payload.UserId, offlineUser.Friends) {
					log.Printf("Sending offline notification to friend (userId):%d\n", client.Payload.UserId)
					encoder := json.NewEncoder(client.Conn)

					onlineStatus := common.OnlineStatus{
						UserId: offlineUser.UserId,
						Online: false,
					}

					if err := encoder.Encode(&onlineStatus); err != nil {
						log.Println(err)
					}

				} else {
					log.Printf("Skipping, userId:%d %s%d\n", client.Payload.UserId, "is not friend of userId:", offlineUser.UserId)
				}
			}
		}
	}
}

func main() {
	port := ":8081"
	listener, _ := net.Listen("tcp", port)

	defer listener.Close()
	log.Printf("Server started on port %s\n", port)

	onlineClients := make(map[int]common.TCPClient)

	var onlineClientsChan = make(chan common.TCPClient)
	var offlineClientsChan = make(chan common.Payload)

	for {
		conn, _ := listener.Accept()
		go handleConnection(conn, onlineClientsChan, offlineClientsChan)
		go trackOnlineStatus(onlineClients, onlineClientsChan, offlineClientsChan)
	}
}
