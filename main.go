package main

import (
	"my-im/userserver"
)

func main()  {


	serv := userserver.NewServer("127.0.0.1", 8888)

	serv.Start()

}
