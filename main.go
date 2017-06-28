package main

import (
	"testJsonRpc/dao"
	"net/rpc"
	"net"
	"log"
	"net/rpc/jsonrpc"
)



func startServer() {
	users := new(dao.Users)

	server := rpc.NewServer()
	server.Register(users)

	server.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)

	l, e := net.Listen("tcp", ":8222")
	if e != nil {
		log.Fatal("listen error:", e)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go server.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}

func main() {
	if err:=dao.SetupDb("users.db"); err!=nil {
		panic(err)
	}
	startServer()
	defer dao.GetDb().Close()
}
