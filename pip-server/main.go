package main

// #cgo LDFLAGS: libts.a -lpthread -ldl -lrt -lssl -lcrypto
// #include <stdlib.h>
// #include "ts.h"
// #include "test.h"
import "C"

import (
	"net"
	"unsafe"

	log "github.com/Sirupsen/logrus"
	pb "github.com/infobloxopen/themis/pipservice"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

// server is used to implement helloworld.GreeterServer.
type server struct{}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Debug("enter SayHello()")
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func (s *server) GetCategories(ctx context.Context, in *pb.Request) (*pb.Response, error) {
	//log.Debug("enter GetCategories()")
	url := in.GetQueryURL()
	cURL := C.CString(url)
	defer C.free(unsafe.Pointer(cURL))
	// Get the categories in C string
	cCategories := C.RateUrl(cURL)
	defer C.free(unsafe.Pointer(cCategories))

  goCategories := C.GoString(cCategories)
	resp := pb.Response{}
	if goCategories != "uncategorized" {
		resp.Status = pb.Response_OK
	} else {
		resp.Status = pb.Response_NOTFOUND
	}

	resp.Categories = goCategories

	return &resp, nil
}

func main() {
	log.SetLevel(log.DebugLevel)
	log.Infoln("Iniitializing McAfee TS SDK")
	_ = C.InitSDK()
	defer C.DestroySDK()

	log.Infoln("Starting server on port: %s", port)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterPIPServer(s, &server{})
	//pb.RegisterGreeterServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	log.Infoln("Serving requests....")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
