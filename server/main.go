package main

import (
	"context"
	"fmt"
	"log"
	"net"

	kafkaApp "github.com/mannanmcc/bookshop/internal/kafkaapp"
	zap "go.uber.org/zap"
	"google.golang.org/grpc"

	protoOrder "github.com/mannanmcc/proto/order/pb"
)

type Server struct {
	//this method must be embedded to have forward compatible implementations which is mentioned in the proto buff
	protoOrder.UnimplementedOrderServiceServer
}

func getServer() *Server {
	return &Server{}
}

func (s *Server) PlaceOrder(ctx context.Context, order *protoOrder.OrderRequest) (*protoOrder.OrderResponse, error) {
	orderLogging, _ := zap.NewProduction()
	orderLogging.Info("new order has been placed")
	kafkaApp.Publish()
	return &protoOrder.OrderResponse{
		OrderNumber: "12345",
	}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	//register our struct to the grpc server
	fmt.Print("starting grpc server")
	protoOrder.RegisterOrderServiceServer(grpcServer, getServer())
	if err := grpcServer.Serve(listener); err != nil {
		fmt.Print("error to listen the grpc server")
		log.Fatalf("failed to serve: %v", err)
	}
}
