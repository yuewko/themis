package main

// #cgo LDFLAGS: libts.a -lpthread -ldl -lrt -lssl -lcrypto
// #include <stdlib.h>
// #include "ts.h"
// #include "test.h"
import "C"

import (
	"encoding/json"
//	"fmt"
	"net"
	"os"
	"unsafe"

	log "github.com/Sirupsen/logrus"

	pb "github.com/infobloxopen/themis/pip-service"

	"github.com/infobloxopen/themis/pdp"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// func C.Init() int

func init() {
    _ = C.Init()
}

// server is used to implement helloworld.GreeterServer.
type server struct {
	CategoryMap map[string]string
}

const supportedQueryType = "domain-category"

func getCategoryMap(categoryFile string) (map[string]string, error) {
	categoryMap := make(map[string]string)

	f, err := os.Open(categoryFile)
	if err != nil {
		log.Fatalf("cannot open category file '%s': '%s'\n", categoryFile, err)
	}
	defer f.Close()

	dec := json.NewDecoder(f)

	if err = dec.Decode(&categoryMap); err != nil {
		log.Fatalf("cannot decode category file '%s': '%s\n", categoryFile, err)
	}

	return categoryMap, err
}

func (s *server) GetAttribute(ctx context.Context, in *pb.Request) (*pb.Response, error) {
	responseStatus := pb.Response_OK

	queryType := in.GetQueryType()
	if queryType != supportedQueryType {
		log.Errorf("Query type '%s' is not supported", queryType)
		responseStatus = pb.Response_SERVICEERROR
	}
	inAttrs := in.GetAttributes()
	// fmt.Printf("inAttrs[0] is '%v'\n", inAttrs[0])
	domainStr := inAttrs[0].GetValue()
/*
	category, ok := s.CategoryMap[domainStr]
	if !ok {
		category = "unknown category"
		responseStatus = pb.Response_NOTFOUND
	}
*/
	c_category := C.RateUrl(C.CString(domainStr))
	var category string
	if c_category != nil {
		category = C.GoString(c_category)
		C.free(unsafe.Pointer(c_category))
	} else {
		category = ""
	}
	respAttr := &pb.Attribute{Type: pdp.TypeListOfStrings, Value: category}
	// fmt.Printf("respAttr is '%v'\n", respAttr)
	values := []*pb.Attribute{respAttr}
	return &pb.Response{Status: responseStatus, Values: values}, nil
}

func main() {
	// _ = C.RateUrl(C.CString("www.abcnews.com"))
	pipServiceName := "mcafee-ts"
	conn, err := pb.GetPIPConnection(pipServiceName)
	if err != nil {
		log.Fatalf("Cannot get PIP connection for PIP service '%s': %s", pipServiceName, err)
	}

	lis, err := net.Listen("tcp", conn)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	categoryMap, err := getCategoryMap("category-map.json")
	if err != nil {
		log.Fatalf("getCategoryMap error: %s\n", err)
		os.Exit(1)
	}

	s := grpc.NewServer()
	pb.RegisterPIPServer(s, &server{CategoryMap: categoryMap})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
