package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc"

	log "github.com/Sirupsen/logrus"
	"github.com/infobloxopen/themis/pip-client/cfg"
	pb "github.com/infobloxopen/themis/pipservice"
)

type QueryResult struct {
	total uint64
	found uint64
}

var wg = new(sync.WaitGroup)

func main() {
	var (
		testData   = []string{}
		stop       = make(chan bool, 1)
		testResult chan QueryResult
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
	testResult = make(chan QueryResult, conf.NumOfWorkers)
	log.Debugf("inital test results: %+v", testResult)

	// create the gRPC client
	conn, err := grpc.Dial(conf.ServerAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewPIPClient(conn)

	runConcurrentQuery(conf.NumOfWorkers, conf.TestDuration, c, &testData, testResult, stop)
	qps, hitRate := aggregateResults(testResult, 5)
	//log.Infof("QPS: %d, Hit Rate: %v", qps, hitRate)
	writeResultToFile(conf.TestResult, conf.NumOfWorkers, qps, hitRate)
	log.Infoln("# Test is finished")
	log.Infof("# Summary: Number of threads: %d, QPS: %d, HitRate: %v",
		conf.NumOfWorkers, qps, hitRate)

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

// example func to make a single query at a time
func queryURL(url string) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(cfg.Config().ServerAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewPIPClient(conn)

	// Contact the server and print out its response
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
func queryFunc(client pb.PIPClient, testList *[]string, results chan QueryResult, stop chan bool) {
	log.Debug("Enter queryFunc goroutines...")
	var result QueryResult
	result.total = 0
	result.found = 0
	wg.Add(1)
	defer wg.Done()
	rSource := rand.NewSource(time.Now().UnixNano())
	r := rand.New(rSource)
	indexRange := len(*testList)

	for {
		select {
		case <-stop:
			log.Infoln("stopped gorouting gracefully.")
			log.Infof("send test result: %+v", result)
			results <- result
			return

		default:
			randomIndex := r.Intn(indexRange)
			log.Debugf("Get a random index: %d", randomIndex)
			url := (*testList)[randomIndex]
			log.Debugf("make query: %s ", url)

			req := pb.Request{QueryURL: url}

			resp, err := client.GetCategories(context.Background(), &req)
			if err != nil {
				log.Errorf("gRPC request GetCategories() error: %v", err)
			}
			if status := resp.GetStatus(); status == pb.Response_OK {
				result.found++
			}

			result.total++

			log.Debugf("%s is categorized as: \t%s", url, resp.Categories)
		}
	}

}

func runConcurrentQuery(numOfWorker int, duration int, client pb.PIPClient,
	testList *[]string, results chan QueryResult, stop chan bool) {
	log.Debugf("Start %d of query goroutines", numOfWorker)
	startTime := time.Now()
	for i := 0; i < numOfWorker; i++ {
		go queryFunc(client, testList, results, stop)
	}
	// run the test for duration seconds
	time.Sleep(time.Duration(duration) * time.Second)
	close(stop)
	endTime := time.Now()
	wg.Wait()
	log.Infof("Start Time: \t%s", startTime)
	log.Infof("End Time: \t%s", endTime)
}

func aggregateResults(results chan QueryResult, duration int) (qps uint32, rate float32) {
	num := len(results)
	log.Debugf("num of results: %d", num)
	var total uint64
	var found uint64
	total = 0
	found = 0

	for i := 0; i < num; i++ {
		val := <-results
		log.Debugf("The %d result: %+v", i, val)
		total += val.total
		found += val.found
	}
	log.Infof("Total queries for all threads: %d", total)
	log.Infof("Total categorized queries for all threads: %d", found)
	qps = uint32(total / uint64(duration))
	var hitRate float32
	hitRate = float32(float64(found) / float64(total))
	close(results)

	return qps, hitRate
}

func writeResultToFile(outFile string, numThread int, qps uint32, hitRate float32) {

	f, err := os.OpenFile(outFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	s := fmt.Sprintf("%d \t%d \t%v\n", numThread, qps, hitRate)
	_, err = f.WriteString(s)

	if err != nil {
		log.Fatal(err)
	}

}
