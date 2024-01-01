package main

import (
	"context"
	pb "github.com/panjiawan/learningExample/grpc/protocol/pb/ecommerce"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"io"
	"log"
	"time"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:8009", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := pb.NewOrderManagementClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	retrievedOrder, err := client.GetOrder(ctx, &wrapperspb.StringValue{Value: "101"})
	if err != nil {
		panic(err)
	}
	log.Print("GetOrder Response -> : ", retrievedOrder)

	stream, err := client.SearchOrders(ctx, &wrapperspb.StringValue{Value: "e"})
	if err != nil {
		panic(err)
	}

	for {
		order, err := stream.Recv()
		if err == io.EOF {
			break
		}

		log.Println("SearchOrders Result: ", order)
	}

	stream2, err := client.UpdateOrders(ctx)
	if err != nil {
		panic(err)
	}
	if err := stream2.Send(&pb.Order{
		Id:          "00",
		Items:       []string{"A", "B"},
		Description: "A with B",
		Price:       0.11,
		Destination: "ABC",
	}); err != nil {
		panic(err)
	}
	if err := stream2.Send(&pb.Order{
		Id:          "01",
		Items:       []string{"C", "D"},
		Description: "C with D",
		Price:       1.11,
		Destination: "ABCDEFG",
	}); err != nil {
		panic(err)
	}

	res, err := stream2.CloseAndRecv()
	if err != nil {
		panic(err)
	}
	log.Printf("Update Orders Res : %s", res)

	stream3, err := client.ProcessOrders(ctx)
	if err != nil {
		panic(err)
	}
	//发送
	go func() {
		if err := stream3.Send(&wrapperspb.StringValue{Value: "101"}); err != nil {
			panic(err)
		}

		if err := stream3.Send(&wrapperspb.StringValue{Value: "102"}); err != nil {
			panic(err)
		}

		if err := stream3.CloseSend(); err != nil {
			panic(err)
		}
	}()
	//读取
	for {
		combinedShipment, err := stream3.Recv()
		if err == io.EOF {
			break
		}
		log.Println("Combined shipment : ", combinedShipment.OrderList)
	}
}
