package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync/atomic"
	"time"
)

var delayTime int
var totalReceivedRequests atomic.Int64
var totalResponses atomic.Int64
var receivedRequestsInOneSecond atomic.Int64
var firstReceiveTime time.Time
var lastResponseTime time.Time
var isFirstRequest atomic.Bool

func handleRequest(w http.ResponseWriter, r *http.Request) {
	totalReceivedRequests.Add(1)
	if (isFirstRequest.CompareAndSwap(true, false)) {
		firstReceiveTime = time.Now()
		go func() {
			ticker := time.NewTicker(time.Second)
			for range ticker.C {
				num := receivedRequestsInOneSecond.Load()
				receivedRequestsInOneSecond.Store(0)
				fmt.Printf("%s, received about %d requests in one second\n", time.Now().Format(time.RFC3339), num)
			}
		}()
	}
	// go func(w http.ResponseWriter) {
		time.Sleep(time.Millisecond * time.Duration(delayTime))

		w.WriteHeader(http.StatusOK)
		totalResponses.Add(1)
		receivedRequestsInOneSecond.Add(1)
		lastResponseTime = time.Now()
	// }(w)
}


func getServerInformation() string {
	rsp := fmt.Sprintf("total received requests: %d\n", totalReceivedRequests.Load())
	rsp += fmt.Sprintf("total responses: %d\n", totalResponses.Load())
	rsp += fmt.Sprintf("first receive time: %s\n", firstReceiveTime.Format(time.RFC3339))
	rsp += fmt.Sprintf("last response time: %s\n", lastResponseTime.Format(time.RFC3339))

	return rsp
}

func handleAdmin(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handle admin")
	rsp := getServerInformation()
	fmt.Fprintf(w, "%s", rsp)
}

func main() {
	http.HandleFunc("/abc/123", handleRequest)
	http.HandleFunc("/admin", handleAdmin)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			if sig == os.Interrupt {
				info := getServerInformation()
				fmt.Println(info)
				os.Exit(0)
			}
		}
	}()

	totalReceivedRequests.Store(0)
	totalResponses.Store(0)
	receivedRequestsInOneSecond.Store(0)
	isFirstRequest.Store(true)
	port := ""
	if len(os.Args) > 2 {
		port = os.Args[1]
		var err error
		delayTime, err = strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Printf("convert delayTime %s failed\n", os.Args[2])
			os.Exit(1)
		}
	}else {
		fmt.Println("usage: ./mock_http_server port delayTime[Millisecond]")
		os.Exit(1)
	}


	host := ":" + port
	fmt.Printf("Listening on host %s\n", host)

	err := http.ListenAndServe(host, nil)
	if err != nil {
		panic(err)
	}
}