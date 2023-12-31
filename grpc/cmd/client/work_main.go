package main

import (
	"context"
	"github.com/panjiawan/learningExample/grpc/protocol/pb/work"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

func main() {
	grpcConn, err := grpc.Dial("127.0.0.1:8001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer grpcConn.Close()

	grpcClient := work.NewWorkWindowClient(grpcConn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	resChan := make(chan *work.Response)
	resErr := make(chan error)
	go func() {
		res, err := grpcClient.GetWork(ctx, &work.Request{Name: "潘家湾", Age: 18})
		if err != nil {
			resErr <- err
		} else {
			resChan <- res
		}
	}()
	select {
	case res := <-resChan:
		log.Println("收到指示：", res)
	case res := <-resErr:
		log.Println("出错：", res)
	case <-ctx.Done():
		log.Println("超时了。。。。。")
	}
	time.Sleep(time.Second * 1)
}
