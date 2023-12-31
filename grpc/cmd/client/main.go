package main

import (
	"context"
	"github.com/panjiawan/learningExample/grpc/protocol/pb/hello"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

var (
	addr = "localhost:8000"
	name = "潘家湾"
)

func main() {
	//连接server端
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("grpc.Dial err")
	}
	defer conn.Close()

	c := hello.NewGreeterClient(conn)                                     //创建一个用于调用 gRPC 服务方法的客户端实例
	ctx, cancel := context.WithTimeout(context.Background(), time.Second) //设置超时
	defer cancel()
	r, err := c.SayHello(ctx, &hello.HelloRequest{Name: name}) //调用接口的方法（SayHello()）
	if err != nil {
		log.Fatalf("#{err}")
	}
	log.Println(r.Replay)
}
