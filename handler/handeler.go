package handler

import "my-im/userserver"

// 抽象出来一层处理用户和server的对应关系
type Handler struct {
	serv *userserver.Server
	user *userserver.User
}


func(handle *Handler) handle(){

}