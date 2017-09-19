package main

// #cgo LDFLAGS: -L/usr/local/lib64 -lcfi -lm -ldl
// #include <stdlib.h>
// #include "cfi.h"
// #include "test.h"
import "C"

import (
	"bufio"
	"fmt"
	"os"
	"unsafe"

	"github.com/Jeffail/tunny"
)


var Id int


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
	if c_category != nil {
		//category := C.GoString(c_category)
		//fmt.Printf("Category is '%s'\n\n", category)
	}
	fmt.Println("")

	return "Dummy"
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


func main() {
	// _ = C.RateUrl(C.CString("www.abcnews.com"))
	numWorkers := 4
	workers := make([]tunny.TunnyWorker, numWorkers)
	for i, _ := range workers {
		workers[i] = &CategoryWorker{}
	}
	pool, _ := tunny.CreateCustomPool(workers).Open()
	defer pool.Close()


	url_lst := "urls.lst"
	f, err := os.Open(url_lst)
	if err != nil {
		fmt.Printf("Cannot open '%s'\n", url_lst)
		os.Exit(1)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		url := scanner.Text()

		fmt.Printf("Getting category for '%s'...\n", url)
		res, _ := pool.SendWork(url)
		fmt.Printf("res is '%s'\n\n", res);
	}
}
