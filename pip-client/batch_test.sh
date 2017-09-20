#!/bin/bash
#
# Usage of ./pip-client:
#   -log-level string
#     	[debug | info | warning | error | fatal | panic] set log level (default "debug")
#   -num-of-workers int
#     	Concurrent threads which are making queries (default 16)
#   -server-addr string
#     	pip server address (default "127.0.0.1:5368")
#   -test-data string
#     	test data file contains list of domains (default "test.data")
#   -test-duration int
#     	seconds to run the test (default 10)
#   -test-result string
#     	test result file contains test results (default "test.output")
# example:
# ./pip-client -log-level info -num-of-workers 4 -server-addr "127.0.0.1:50051" -test-data test.data -test-duration 10 -test-result test.ouput 2>&1 | tee log.txt
#
set -x
# test data
test_data_set="blacklist.data threat.data" 
# duration for each test
# TODO: please use 20 for official test
seconds=20
logLevel=info
serverAddr="127.0.0.1:50051"
# numOfThreads:
# TODO: please use 1032 for official test
#max_threads=1032
max_threads=1032
for test_data in $test_data_set
do
    tm=`date +%F-%S`
    output="$test_data".result."$tm".txt
    echo "test data: $test_data"
    echo "test output: $output"
    num_of_worker=1
    while [ $num_of_worker -lt $max_threads ]; do
    	echo "Number of Workers: $num_of_worker"
        ./pip-client -log-level $logLevel -num-of-workers $num_of_worker -server-addr $serverAddr -test-data $test_data -test-duration $seconds -test-result $output
        echo "Increase num_of_worker"      
        if [ $num_of_worker -lt 4 ]; then
            let num_of_worker+=1
        elif [ $num_of_worker -ge 4 ] && [ $num_of_worker -le 32 ]; then
            let num_of_worker+=4
        else
            let num_of_worker+=32
        fi
    done
done

