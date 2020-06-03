package common

import (
	"net"
)

type TCPClient struct {
	Payload Payload
	Conn    net.Conn
}

type Payload struct {
	UserId  int
	Friends []int
}

type OnlineStatus struct {
	UserId int
	Online bool
}
