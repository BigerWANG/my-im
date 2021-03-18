package main

import "my-im/server"

func main()  {


	serv := server.NewServer("127.0.0.1", 8888)

	serv.Start()

}
