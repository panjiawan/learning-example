package work

import (
	"context"
	"github.com/panjiawan/learningExample/grpc/protocol/pb/work"
	"log"
	"time"
)

type Server struct {
	work.UnimplementedWorkWindowServer
}

func (s *Server) GetWork(ctx context.Context, req *work.Request) (*work.Response, error) {
	log.Println("服务端获取到的姓名：", req.GetName())
	log.Println("服务端获取到的年龄：", req.GetAge())
	log.Println("sleep 3s....")
	time.Sleep(time.Second * 3)
	if req.GetAge() < 18 {
		return &work.Response{Work: req.GetName() + "你的工作是上学"}, nil
	} else {
		return &work.Response{Work: req.GetName() + "开始搬砖！！！"}, nil
	}
}
