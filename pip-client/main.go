package main

import (
	//"google.golang.org/grpc"
	//"google.golang.org/grpc/reflection"
	//	"fmt"

	"bufio"
	"context"
	"encoding/json"
	"os"
	"strings"
	"sync"

	"google.golang.org/grpc"

	log "github.com/Sirupsen/logrus"
	"github.com/infobloxopen/themis/pip-client/cfg"
	ps "github.com/infobloxopen/themis/pip-service"
	pb "github.com/infobloxopen/themis/pipservice"
)

const (
	address = "localhost:50051"
	//address     = "10.82.16.198:50051"
	defaultName = "pip-client"
)

type QueryResult struct {
	total uint64
	found uint64
}

var wg = new(sync.WaitGroup)

func main() {
	var (
		testData = []string{}
		//testResult chan QueryResult
	)

	if err := cfg.Load(); err != nil {
		log.Fatalf("Failed to load config: %s", err)
	}

	conf := cfg.Config()
	block, err := json.MarshalIndent(conf, "", "  ")
	if err != nil {
		log.Errorf("error: %s", err)
	}
	log.Infof("Configuration: \n%s", string(block))

	log.SetLevel(getLogLevelFromStr(conf.LogLevel))
	log.Debug("tests....")

	log.Info("Read test data into list...")
	log.Debugf("before read: %+v", testData)

	err = loadTestData(conf.TestData, &testData)
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Debugf("after read: %+v", testData)
	log.Debugf("%d of domains are loaded.", len(testData))

	// init test resultsf chan to collect results from multiple threads
	//testResult = make(chan QueryResult, conf.NumOfWorkers)
	makeQuery(testData)
	// test()

}

func getLogLevelFromStr(level string) log.Level {
	var ret log.Level
	switch level {
	case "debug":
		ret = log.DebugLevel
	case "info":
		ret = log.InfoLevel
	case "warning":
		ret = log.WarnLevel
	case "error":
		ret = log.ErrorLevel
	case "fatal":
		ret = log.FatalLevel
	case "panic":
		ret = log.PanicLevel
	default:
		log.Errorf("%s could not be parsed", level)
	}
	return ret
}

func loadTestData(inputFile string, lst *[]string) error {
	fileHandle, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer fileHandle.Close()

	fileScanner := bufio.NewScanner(fileHandle)

	for fileScanner.Scan() {
		line := strings.TrimSpace(fileScanner.Text())
		//log.Debug(line)
		// get rid of empty lines
		if len(line) >= 1 {
			*lst = append(*lst, line)
		}

	}
	//log.Debugln(*lst)
	return nil
}

func createClient(address string, domain string) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := ps.NewPIPClient(conn)
	attrList := []*ps.Attribute{}
	attr := ps.Attribute{Id: "1", Type: 6, Value: domain}
	attrList = append(attrList, &attr)
	log.Debugf("attrList: %+v", attrList)
	request := &ps.Request{QueryType: "domain-category", Attributes: attrList}
	log.Debugf("request: %+v", request)

	r, err := c.GetAttribute(context.Background(), request)
	if err != nil {
		log.Errorln(err)
	}
	if r.Status != ps.Response_OK {
		log.Debugf("Return Status: %v", r.Status)
	}

	res := r.GetValues()
	log.Debugf("Return values: %+v", res)

}

func makeQuery(testLst []string) {
	log.Debugf("Input Test list: %+v", testLst)
	for _, domain := range testLst {
		log.Debugf("query: %s", domain)
		// createClient("127.0.0.1:5368", domain)
		go queryURL(domain)
	}
}

func queryURL(url string) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewPIPClient(conn)

	// Contact the server and print out its response.

	// r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: name})
	// if err != nil {
	// 	log.Fatalf("could not greet: %v", err)
	// }
	// log.Printf("Greeting: %s", r.Message)

	req := pb.Request{QueryURL: url}

	resp, err := c.GetCategories(context.Background(), &req)
	if err != nil {
		log.Fatalf("could not be categorized: %v", err)
	}
	log.Debugf("Categorized: %s", resp.Categories)

}

// used for go routing
func queryFunc(client pb.PIPClient, wg sync.WaitGroup, results chan QueryResult, stop chan bool) {
	log.Debug("Enter queryFunc go routing...")
	var result QueryResult
	result.total = 0
	result.found = 0
	wg.Add(1)
	defer wg.Done()
	for {
		select {
		case <-stop:
			log.Info("stopped gorouting gracefully.")
			log.Info("send test result: %+v", result)
			results <- result
			return

		default:

			log.Debug("make query ...")
			result.total++
		}
		log.Info("send test result: %+v", result)

	}

}
