package main

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	pb "loadbalance_test/addr"
	"log"
	"net"
	"net/http"
)

const (
	PLAIN_HTTP_PORT = ":80"
	GRPC_PORT       = ":30051"
)

var localIp string

func getOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP
}

func addr(resp http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(resp, localIp)
}

func health(resp http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(resp, "OK")
}

type serverGrpc struct{}

func (s *serverGrpc) GetAddr(ctx context.Context, in *empty.Empty) (*pb.AddrReply, error) {
	return &pb.AddrReply{Addr: localIp}, nil
}

func main() {
	localIp = getOutboundIP().String()

	// STEP 1-1：啟動HTTP服務
	go func() {
		log.Printf("Starting http server on %s...", PLAIN_HTTP_PORT)
		http.HandleFunc("/addr", addr)
		http.HandleFunc("/health", health)
		if err := http.ListenAndServe(PLAIN_HTTP_PORT, nil); err != nil {
			log.Fatalf("Fail to start HTTP server: %v", err)
		}
	}()

	// STEP 2-1：定義要監聽的 port 號
	log.Printf("Starting gRPC server on %s...", GRPC_PORT)
	lis, err := net.Listen("tcp", GRPC_PORT)
	if err != nil {
		log.Fatalf("Fail to start gRPC server: %v", err)
	}
	// STEP 2-2：使用 gRPC 的 NewServer 方法來建立 gRPC Server 的實例
	s := grpc.NewServer()

	// STEP 2-3：在 gRPC Server 中註冊 service 的實作
	// 使用 proto 提供的 RegisterRouteGuideServer 方法，並將 serverGrpc 作為參數傳入
	pb.RegisterAddrServer(s, &serverGrpc{})

	// STEP 2-4：啟動 grpcServer，並阻塞在這裡直到該程序被 kill 或 stop
	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
