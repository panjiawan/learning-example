package ecommerce

import (
	"context"
	"fmt"
	pb "github.com/panjiawan/learningExample/grpc/protocol/pb/ecommerce"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"io"
	"log"
	"strings"
)

const (
	orderBatchSize = 3
)

var orders = make(map[string]pb.Order)

type OrderManagementImpl struct {
	pb.UnimplementedOrderManagementServer
}

func (s *OrderManagementImpl) GetOrder(ctx context.Context, orderId *wrapperspb.StringValue) (*pb.Order, error) {
	orders["101"] = pb.Order{
		Id:          "101",
		Items:       nil,
		Description: "我是101",
		Price:       0,
		Destination: "",
	}
	ord, exists := orders[orderId.Value]
	if exists {
		return &ord, status.New(codes.OK, "").Err()
	}

	return nil, status.Errorf(codes.NotFound, "Order does not exist. : ", orderId)
}

func (s *OrderManagementImpl) SearchOrders(query *wrapperspb.StringValue,
	stream pb.OrderManagement_SearchOrdersServer) error {
	orders["101"] = pb.Order{
		Id:          "101",
		Items:       []string{"a", "b", "c"},
		Description: "我是101",
		Price:       0,
		Destination: "",
	}
	orders["102"] = pb.Order{
		Id:          "102",
		Items:       []string{"e", "f", "g"},
		Description: "我是102",
		Price:       0,
		Destination: "",
	}
	orders["103"] = pb.Order{
		Id:          "103",
		Items:       []string{"eh", "i", "j"},
		Description: "我是103",
		Price:       0,
		Destination: "",
	}
	for _, order := range orders {
		for _, str := range order.Items {
			if strings.Contains(str, query.Value) {
				err := stream.Send(&order)
				if err != nil {
					return fmt.Errorf("error send: %v", err)
				}
			}
		}
	}

	return nil
}

// 在这段程序中，我们对每一个 Recv 都进行了处理
// 当发现 io.EOF (流关闭) 后，需要将最终的响应结果发送给客户端，同时关闭正在另外一侧等待的 Recv

func (s *OrderManagementImpl) UpdateOrders(stream pb.OrderManagement_UpdateOrdersServer) error {
	ordersStr := "Updated Order IDs : "
	for {
		order, err := stream.Recv()
		if err == io.EOF {
			// Finished reading the order stream.
			return stream.SendAndClose(
				&wrapperspb.StringValue{Value: "Orders processed " + ordersStr})
		}
		// Update order
		orders[order.Id] = *order

		log.Println("Order ID ", order.Id, ": Updated")
		ordersStr += order.Id + ", "
	}
}

func (s *OrderManagementImpl) ProcessOrders(stream pb.OrderManagement_ProcessOrdersServer) error {
	batchMarker := 1
	orders = map[string]pb.Order{
		"101": {
			Id: "101",
			Items: []string{
				"Google",
				"Baidu",
			},
			Description: "example",
			Price:       0,
			Destination: "example",
		},
	}

	var combinedShipmentMap = make(map[string]pb.CombinedShipment)

	for {
		orderId, err := stream.Recv() //读取消息
		log.Printf("Reading Proc order : %s", orderId)
		if err == io.EOF {
			log.Printf("EOF : %s", orderId)
			for _, shipment := range combinedShipmentMap {
				//写入响应
				if err := stream.Send(&shipment); err != nil {
					return err
				}
			}
			return nil
		}
		if err != nil {
			log.Println(err)
			return err
		}

		destination := orders[orderId.GetValue()].Destination
		shipment, found := combinedShipmentMap[destination]

		if found {
			ord := orders[orderId.GetValue()]
			shipment.OrderList = append(shipment.OrderList, &ord)
			combinedShipmentMap[destination] = shipment
		} else {
			comShip := pb.CombinedShipment{Id: "cmb - " + (orders[orderId.GetValue()].Destination), Status: "Processed!"}
			ord := orders[orderId.GetValue()]
			comShip.OrderList = append(shipment.OrderList, &ord)
			combinedShipmentMap[destination] = comShip
			log.Print(len(comShip.OrderList), comShip.GetId())
		}

		if batchMarker == orderBatchSize {
			for _, comb := range combinedShipmentMap {
				log.Printf("Shipping : %v -> %v", comb.Id, len(comb.OrderList))
				if err := stream.Send(&comb); err != nil {
					return err
				}
			}
			batchMarker = 0
			combinedShipmentMap = make(map[string]pb.CombinedShipment)
		} else {
			batchMarker++
		}
	}
}
