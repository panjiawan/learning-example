package main

import (
	"github.com/panjiawan/learningExample/grpc/protocol/pb/ecommerce"
	ecommerceService "github.com/panjiawan/learningExample/grpc/service/ecommerce"
	"google.golang.org/grpc"
	"net"
)

func main() {
	s := grpc.NewServer()
	ecommerce.RegisterOrderManagementServer(s, &ecommerceService.OrderManagementImpl{})

	lis, err := net.Listen("tcp", ":8009")
	if err != nil {
		panic(err)
	}
	defer func() {
		s.Stop()
		lis.Close()
	}()

	if err := s.Serve(lis); err != nil {
		panic(err)
	}
}
