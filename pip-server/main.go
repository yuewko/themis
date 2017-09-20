//go:generate protoc -I ../helloworld --go_out=plugins=grpc:../helloworld ../helloworld/helloworld.proto

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

// func C.Init() int

func init() {
	_ = C.Init()
}

// server is used to implement helloworld.GreeterServer.
type server struct{}

func (s *server) GetCategories(ctx context.Context, in *pb.Request) (*pb.Response, error) {
	//log.Debug("enter GetCategories()")
	url := in.GetQueryURL()

	c_category := C.RateUrl(C.CString(url))

	var category string
	resp := pb.Response{}

	if c_category != nil {
		resp.Status = pb.Response_OK
		category = C.GoString(c_category)
		C.free(unsafe.Pointer(c_category))

	} else {
		resp.Status = pb.Response_NOTFOUND
		category = "un-categorized"
	}
	resp.Categories = category

	return &resp, nil
}

func main() {
	log.SetLevel(log.DebugLevel)

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
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
