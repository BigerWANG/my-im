package server

import (
	"fmt"
	"log"
	"net"
)

type Server struct {
	IP string
	Port int
}

type OnlineMap map[string]string


func NewServer(ip string, port int) *Server {
	return  &Server{
		IP:   ip,
		Port: port,
	}
}


func(s *Server) Start(){
	// Socket listen
	listener, err := net.Listen("tcp4", fmt.Sprintf("%s:%d", s.IP, s.Port))
	if  err != nil{
		log.Fatal()
	}

	defer func() {
		err = listener.Close()
		log.Println(err)
	}()
	// accept
	for {
		conn, err := listener.Accept()
		if err !=nil{
			if err := fmt.Errorf(err.Error()); err!=nil{
			}
			continue
		}
		go Handler(conn)
	}

}

func Handler(conn net.Conn) {
	fmt.Println("当前连接建立成功")
	content, err := conn.Read([]byte{100})
	if err != nil{
		fmt.Println(err)
	}
	fmt.Println(content)
}