package userserver

import (
	"fmt"
	"net"
	"strings"
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

// 给你当前客户端发送消息
func (u *User)sendMsg(msg string){
	if _, err := u.conn.Write([]byte(msg)); err != nil{
		fmt.Println(err)
	}
}

func(u *User)rename(newname string){
	if u.server.OnlineMap[newname] != nil{
		u.sendMsg(fmt.Sprintf("[%s]此名称已经被占用了~，换个名字吧", newname))
		return
	}
	u.server.Mu.Lock()
	// 需要先删除原有的key
	delete(u.server.OnlineMap, u.Name)
	u.server.OnlineMap[newname] = u
	u.Name = newname
	u.sendMsg("昵称修改成功")
	u.server.Mu.Unlock()
}


func(u *User)who(){
	for _, cli := range u.server.OnlineMap{
		currMsg := fmt.Sprintf("当前[%s]在线", cli.Name)
		u.sendMsg(currMsg)
	}
}

// 用户消息处理
func(u *User)DoMessage(msg string) error {
	u.server.isalive <- struct{}{}
	switch msg {
	case "help":
		u.sendMsg("获取技能: \n who: 查看当前在线用户\n rename: 重命名你当前的用户\n\tusage rename|zhangsan")
	case "who":
		u.who()
	case "rename":
		u.sendMsg("rename usage: username|<your name>\n")
	default:
		if strings.Contains(msg, "rename|") {
			newname := strings.Split(msg, "|")[1]
			u.rename(newname)
			return nil
		}
		boradcastMsg := "发送了消息: " + msg
		u.server.BoradCast(u, boradcastMsg)
	}
	return nil
}
