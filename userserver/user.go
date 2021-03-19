package userserver

import (
	"fmt"
	"io"
	"net"
)

type User struct {
	Name  string
	Addr string
	UserChan chan string
	conn net.Conn
	server *Server
}


func NewUser(conn net.Conn, serv *Server)*User{

	return &User{
		Name:     conn.RemoteAddr().String(),
		Addr:     conn.RemoteAddr().String(),
		UserChan: make(chan string),
		conn:     conn,
		server:     serv,
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

// 用户上线通知
func(u *User)OnlineNotice(){
	u.server.Mu.Lock()
	// 将当前用户保存在map中
	u.server.OnlineMap[u.Name]=u
	u.server.Mu.Unlock()
	// 广播当前用户上线消息
	u.server.BoradCast(u, "已上线")
}
// 用户下线通
func(u *User)OfflineNotice(){
	u.server.Mu.Lock()
	// 将当前用户从map中删除
	delete(u.server.OnlineMap, u.Name)
	u.server.Mu.Unlock()
	// 广播当前用户上线消息
	u.server.BoradCast(u, "已下线上线")
}
// 用户消息处理
func(u *User)DoMessage() error{
	buf := make([]byte, 1024)
	for {
		n, err := u.conn.Read(buf)
		if err != nil && err != io.EOF{
			return err
		}
		if n > 0 {
			msg := "userserver" + u.Name + "发送了消息: " + string(buf[:n-1])
			u.server.BoradCast(u, msg)
		}
		if n == 0{ // 当read是0 的时候就说明已经下线
			u.OfflineNotice()
			return nil
		}
	}


}
