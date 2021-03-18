package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"context"
	// accesslogconfig "github.com/envoyproxy/go-control-plane/envoy/config/accesslog/v2"
	accesslog "github.com/envoyproxy/go-control-plane/envoy/service/accesslog/v3"
	
	"google.golang.org/grpc"
	"github.com/golang/protobuf/jsonpb"
	// "github.com/golang/protobuf/ptypes"
)

type server struct {
	marshaler jsonpb.Marshaler
}

func main() {
	ctx := context.Background()
	port:= 18090
	runAccessLogServer(ctx, port)
}

func runAccessLogServer(ctx context.Context, port int) {
	grpcServer := grpc.NewServer()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Println("failed to listen")
		log.Panicf("%+v", err)
	}
	// var als accesslog.AccessLogServiceServer
	als := &server{}
	accesslog.RegisterAccessLogServiceServer(grpcServer, als)
	log.Println("Access Server listening")

	go func() {
		log.Printf("Inside Go routine")
		if err = grpcServer.Serve(lis); err != nil {
			log.Printf("%+v", err)
		}
	}()
	<-ctx.Done()

	grpcServer.GracefulStop()
}

func (s *server) StreamAccessLogs(stream accesslog.AccessLogService_StreamAccessLogsServer) error {
	log.Println("Started stream")
	for {
		in, err := stream.Recv()
		log.Println("Received value")
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		str, _ := s.marshaler.MarshalToString(in)
		log.Println(str)
	}
}