package main

// #cgo LDFLAGS: -L/usr/local/lib64 -lcfi -lm -ldl
// #include <stdlib.h>
// #include "cfi.h"
// #include "test.h"
import "C"

import (
	"fmt"
	"net"
	"unsafe"

	log "github.com/Sirupsen/logrus"
	pb "github.com/infobloxopen/themis/pipservice"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"github.com/Jeffail/tunny"
)


var Id int

const (
	port = ":50051"
	numWorkers = 50
)


// Implements TunnyWorker and TunnyExtendedWorker interfaces
type CategoryWorker struct {
	cfi_handle unsafe.Pointer
	ready bool

	id int
}

func (cw *CategoryWorker)TunnyJob(obj interface{}) (res interface{}) {
	url, _ := obj.(string)
	fmt.Printf("Getting category for '%s'...%d\n", url, cw.id)
	c_category := C.RateUrl(cw.cfi_handle, C.CString(url))
	var category string
	if c_category != nil {
		category = C.GoString(c_category)
		//fmt.Printf("Category is '%s'\n\n", category)
	}
	//fmt.Println("")

	return category
}

func (cw *CategoryWorker) TunnyReady() bool {
	return cw.ready
}

func (cw *CategoryWorker) TunnyInitialize() {
	fmt.Println("Calling TunnyInitialize()")
	cw.cfi_handle = C.Init()
	cw.ready = true
	cw.id = Id
	Id++
}

func (cw *CategoryWorker) TunnyTerminate() {
}


type server struct{
	workPool *tunny.WorkPool
}

func (s *server) GetCategories(ctx context.Context, in *pb.Request) (*pb.Response, error) {
	//log.Debug("enter GetCategories()")
	url := in.GetQueryURL()

	resp := pb.Response{}
	workRes, err := s.workPool.SendWork(url)
	if err != nil {
		resp.Categories = workRes.(string)
	}
	return &resp, err
}


func main() {
	log.SetLevel(log.DebugLevel)

        C.StartUp()

	// Create worker thread pool
	workers := make([]tunny.TunnyWorker, numWorkers)
	for i, _ := range workers {
		workers[i] = &CategoryWorker{}
	}
	pool, _ := tunny.CreateCustomPool(workers).Open()
	defer pool.Close()


	log.Infoln("Starting server on port: %s", port)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterPIPServer(s, &server{workPool: pool})
	//pb.RegisterGreeterServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
