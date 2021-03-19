package userserver

import (
	"fmt"
	"log"
	"net"
	"sync"
)

type Server struct {
	IP string
	Port int
	u *User
	OnlineMap map[string]*User
	Message chan string
	Mu sync.RWMutex
}



func NewServer(ip string, port int) *Server {
	return  &Server{
		IP:   ip,
		Port: port,
		OnlineMap: make(map[string]*User),
		Message: make(chan string),

	}
}


func(s *Server) Start(){
	// Socket listen
	listener, err := net.Listen("tcp4", fmt.Sprintf("%s:%d", s.IP, s.Port))
	if  err != nil{
		log.Fatal()
	}
	// 启动监听message的go
	go s.ListenMessager()

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
		go s.Handler(conn)
	}

}

func(s *Server) Handler(conn net.Conn) {
	u := NewUser(conn, s)
	// 启动用户监听程序
	go u.ListenMessage()
	// 将用户输入消息也进行广播
	u.OnlineNotice()
	err := u.DoMessage()
	if err !=nil{
		fmt.Println("there is have an error!!!")
		fmt.Println(err)
		return
	}
}

// 广播消息
func(s *Server)BoradCast(u *User, msg string){
	sendMsg := fmt.Sprintf("[%s]%s: %s", u.Addr, u.Name, msg)
	s.Message <- sendMsg
}


// 监听Message广播消息channnel的goroutine,一旦有消息就发送给全部在线的user
func (s *Server)ListenMessager()  {
	for {
		msg := <- s.Message
		s.Mu.Lock()
		for _, cli := range s.OnlineMap{
			fmt.Println("当前用户", cli.Name, msg)
			cli.UserChan <- msg
		}
		s.Mu.Unlock()
	}
}