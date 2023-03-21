package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "loadbalance_test/addr"
	"log"
	"time"
)

const (
	TIMEOUT       = 5 * time.Second
	WAIT_DURATION = 1 * time.Second
)

func main() {
	initGRPCRoute()
}

func initGRPCRoute() {
	router := gin.New()

	// 根目錄回傳200
	router.GET("/", func(c *gin.Context) {
		log.Println("123132")
		connectGRPCAndShowResponseNoClose()
		time.Sleep(WAIT_DURATION)
	})
	err := router.Run(":8080")
	if err != nil {
		log.Println("Router dead ☠️")
	}
}

func connectGRPCAndShowResponse() {
	fmt.Printf("Making rpc...\n")
	for {
		conn, err := grpc.Dial("api-server-grpc-services:30051", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("failed to dial: %v", err)
		}

		client := pb.NewAddrClient(conn)
		callAndShowResponse(client)
		time.Sleep(WAIT_DURATION)
		conn.Close()
	}

}

var grpcCount int

func callAndShowResponse(client pb.AddrClient) {
	now := time.Now()
	reply, err := client.GetAddr(context.Background(), &empty.Empty{})
	if err != nil {
		panic(err)
	}

	grpcCount++
	fmt.Printf("#%d: %s exec:%v \n", grpcCount, reply.Addr, time.Since(now).Seconds())
}

func connectGRPCAndShowResponseNoClose() {
	conn, err := grpc.Dial("api-server-grpc-services:30051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewAddrClient(conn)
	fmt.Printf("Making rpc...\n")
	for {
		callAndShowResponse(client)
		time.Sleep(WAIT_DURATION)

	}
}
