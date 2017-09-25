// Copyright (c) 2017 Infoblox Inc. All Rights Reserved.

package cfg

import (
	"flag"

	// "github.com/Infoblox-CTO/janus-common/log"
	log "github.com/Sirupsen/logrus"
)

const (
	defaultLogLevel      = "debug"
	defaultTestData      = "test.data"
	defaultTestResult    = "test.output"
	defaultNumOfWorkers  = 16
	defaultTestDuration  = 10
	defaultServerAddress = "localhost:50051"
)

var (

	// log and healthcheck cofig
	logLevel     = flag.String("log-level", defaultLogLevel, "[debug | info | warning | error | fatal | panic] set log level")
	testData     = flag.String("test-data", defaultTestData, "test data file contains list of domains")
	testResult   = flag.String("test-result", defaultTestResult, "test result file contains test results")
	numOfWokers  = flag.Int("num-of-workers", defaultNumOfWorkers, "Concurrent threads which are making queries")
	testDuration = flag.Int("test-duration", defaultTestDuration, "seconds to run the test")

	serverAddress = flag.String("server-addr", defaultServerAddress, "pip server address")
)

// AppConfig is a global config object used
type AppConfig struct {
	LogLevel      string
	TestData      string
	TestResult    string
	NumOfWorkers  int
	TestDuration  int
	ServerAddress string
}

var globalConfig *AppConfig

// Load the application configurworkDirations
func Load() error {
	flag.Parse()

	globalConfig = &AppConfig{}
	globalConfig.LogLevel = *logLevel
	globalConfig.TestData = *testData
	globalConfig.TestResult = *testResult
	globalConfig.NumOfWorkers = *numOfWokers
	globalConfig.TestDuration = *testDuration
	globalConfig.ServerAddress = *serverAddress
	return nil

}

// Config returns configuration info instance
func Config() *AppConfig {
	if globalConfig == nil {
		log.Fatal("Configuration is not loaded")
	}
	return globalConfig
}
