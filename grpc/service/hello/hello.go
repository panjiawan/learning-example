package hello

import (
	"context"
	"github.com/panjiawan/note/grpc/protocol/pb/hello"
	"log"
)

type Server struct {
	hello.UnimplementedGreeterServer
}

func (s *Server) SayHello(ctx context.Context, in *hello.HelloRequest) (*hello.HelloResponse, error) {
	log.Println("SayHello run .....")
	return &hello.HelloResponse{Replay: "你好，" + in.GetName() + " 服务端返回"}, nil
}
