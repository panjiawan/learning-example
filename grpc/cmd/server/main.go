package main

import (
	"github.com/panjiawan/learningExample/grpc/protocol/pb/hello"
	helloService "github.com/panjiawan/learningExample/grpc/service/hello"
	"google.golang.org/grpc"
	"log"
	"net"
)

const port = ":8000"

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()                                  // 创建一个grpc服务器，起到桥梁作用。
	hello.RegisterGreeterServer(s, &helloService.Server{}) //把grpc服务器和实现了接口的实例，注册到一起。

	defer func() {
		s.Stop()
		lis.Close()
	}()

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil { //把grpc服务器绑定到tcp的监听上，使grpc服务器可以处理来自监听的请求
		log.Fatalf("failed to serve: %v", err)
	}
}
