package userserver

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

type Server struct {
	IP string
	Port int
	TimeOut time.Duration
	u *User
	OnlineMap map[string]*User
	Message chan string
	isalive chan struct{}

	Mu sync.RWMutex
}



func NewServer(ip string, port int) *Server {
	return  &Server{
		IP:   ip,
		Port: port,
		TimeOut: 10,
		OnlineMap: make(map[string]*User),
		Message: make(chan string),
		isalive: make(chan struct{}, 1),


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
	s.u = u
	// 将用户输入消息也进行广播
	// 启动用户监听程序
	go u.ListenMessage()
	go s.killer()
	u.OnlineNotice()
	buf := make([]byte, 1024)
	for {
		n, err := u.conn.Read(buf)
		if err != nil && err != io.EOF{
			fmt.Println(err)
			return
		}
		if n == 1{
			u.sendMsg(fmt.Sprintf("[%s]$:", u.Name))
			continue
		}

		if n > 1{ // 当有用户实际输入时
			msg := string(buf[:n-1])
			err = u.DoMessage(msg)
			if err !=nil{
				fmt.Println("there is have an error!!!", err)
				return
			}
			continue
		}
		if n == 0{ // 当read是0 的时候就说明已经下线
			u.OfflineNotice()
			return
		}
	}

}

// 监听用户是否活跃，负责超时强制踢人
func(s *Server)killer(){
	// 为每个用户创建一个计时器
	timer := time.After(s.TimeOut*time.Second)
	for{
		select {
		case <-s.isalive:
			fmt.Println("当前用户活跃")
			timer = time.After(s.TimeOut*time.Second)
		case <-timer:
			//s.u.sendMsg("你已超时")
			s.u.sendMsg("你已超时, 拜拜了您内")
			s.u.conn.Close()
		}
	}

}

// 广播消息
func(s *Server)BoradCast(u *User, msg string){
	sendMsg := fmt.Sprintf("[%s]: %s", u.Name, msg)
	s.Message <- sendMsg
}


// 监听Message广播消息channnel的goroutine,一旦有消息就发送给全部在线的user
func (s *Server)ListenMessager()  {
	for {
		msg := <- s.Message
		s.Mu.Lock()
		for _, cli := range s.OnlineMap{
			cli.UserChan <- msg
		}
		s.Mu.Unlock()
	}
}