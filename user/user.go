package user

import "net"

type User struct {
	Name  string
	Addr string
	UserChan chan string
	conn net.Conn
}

