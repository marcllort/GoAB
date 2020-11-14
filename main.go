package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

var nerr = 0
var nsucc = 0
var nfinished = 0
var channel = make(chan int)

func request(pool chan string, n int, keepalive bool, client *http.Client) {
	for true {
		url := <-pool
		req, err := http.NewRequest("GET", url, nil)
		if keepalive {
			req.Header.Set("Connection", "keep-alive")
		}
		resp, err := client.Do(req)
		defer resp.Body.Close()

		if err == nil {
			nsucc++
		} else {
			nerr++
		}
		nfinished++

		if nfinished >= n {
			channel <- 1
		}
	}
}

func initPool(pool chan string, n int, url string) {
	for i := 0; i < n; i++ {
		pool <- url
	}
}

func ab(url string, nreq int, concurrency int, keepalive bool) {
	tr := &http.Transport{
		DisableKeepAlives: !keepalive,
	}
	client := &http.Client{Transport: tr}

	start := time.Now()
	pool := make(chan string, concurrency)

	go initPool(pool, nreq, url)
	for i := 0; i < concurrency; i++ {
		go request(pool, nreq, keepalive, client)
	}

	//Receive channel data
	<-channel
	close(channel)

	time := float64(time.Since(start).Milliseconds())

	fmt.Println("Time taken for tests: ", time/1000, " seconds")
	fmt.Println("Complete requests: ", nreq-nerr)
	fmt.Println("Failed requests: ", nerr)
	fmt.Println("Time per request: ", time*float64(concurrency)/float64(nreq-nerr), "[ms]")
	fmt.Println("Time per request: ", time/float64(nreq-nerr), "[ms] (across all concurrent requests)")
}

func main() {
	// Default parameters
	keepalive := false
	nrequests := 1
	concurrency := 1

	for i := 1; i < len(os.Args)-1; i++ {
		if os.Args[i] == "-n" {
			nrequests, _ = strconv.Atoi(os.Args[i+1])
		} else if os.Args[i] == "-c" {
			concurrency, _ = strconv.Atoi(os.Args[i+1])
		} else if os.Args[i] == "-k" {
			keepalive = true
		}
	}

	url := os.Args[len(os.Args)-1]

	ab(url, nrequests, concurrency, keepalive)
}
