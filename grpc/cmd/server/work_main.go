package main

import (
	"github.com/panjiawan/note/grpc/protocol/pb/work"
	workService "github.com/panjiawan/note/grpc/service/work"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:8001")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("服务端开始监听8001接口")
	grpcServer := grpc.NewServer()
	work.RegisterWorkWindowServer(grpcServer, &workService.Server{})

	defer func() {
		grpcServer.Stop()
		listener.Close()
	}()

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
