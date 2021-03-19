package server

import (
	"fmt"
	"io"
	"log"
	"my-im/user"
	"net"
	"sync"
)

type Server struct {
	IP string
	Port int
	u *user.User
	OnlineMap map[string]*user.User
	Message chan string
	mu sync.RWMutex
}



func NewServer(ip string, port int) *Server {
	return  &Server{
		IP:   ip,
		Port: port,
		OnlineMap: make(map[string]*user.User),
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
	u := user.NewUser(conn)
	// 启动用户监听程序
	go u.ListenMessage()
	s.mu.Lock()
	// 将当前用户保存在map中
	s.OnlineMap[u.Name]=u
	s.mu.Unlock()
	// 广播当前用户上线消息
	s.BoradCast(u, "已上线")

	// 将用户输入消息也进行广播
	err := s.sendAll(u, conn)
	if err !=nil{
		fmt.Println("there is have an error!!!")
		fmt.Println(err)
		return
	}
}

// 广播消息
func(s *Server)BoradCast(u *user.User, msg string){
	sendMsg := fmt.Sprintf("[%s]%s: %s", u.Addr, u.Name, msg)
	s.Message <- sendMsg
}


// 向在线用户发送消息
func (s *Server)sendAll(currUser *user.User, conn net.Conn) (error error) {
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil && err != io.EOF{
			return err
		}
		msg := "user"+currUser.Name+ "发送了消息: " + string(buf[:n-1])
		s.BoradCast(currUser, msg)
	}

}

// 监听Message广播消息channnel的goroutine,一旦有消息就发送给全部在线的user
func (s *Server)ListenMessager()  {
	for {
		msg := <- s.Message

		s.mu.Lock()
		for _, cli := range s.OnlineMap{
			cli.UserChan <- msg
		}
		s.mu.Unlock()
	}
	
}