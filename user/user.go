package user

import (
	"fmt"
	"net"
)

type User struct {
	Name  string
	Addr string
	UserChan chan string
	conn net.Conn
}


func NewUser(conn net.Conn)*User{

	return &User{
		Name:     conn.RemoteAddr().String(),
		Addr:     conn.RemoteAddr().String(),
		UserChan: make(chan string),
		conn:     conn,
	}

}

// 监听当前User channel的方法，一旦有消息就直接发送给对端客户端
func(u *User)ListenMessage(){

	for {
		msg := <-u.UserChan
		_, err := u.conn.Write([]byte(msg+"\n"))
		if err != nil{
			fmt.Println(err)
		}
	}

}
