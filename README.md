# GoAB
GoLang implementation of Apache Workbench.

Using channels and GoRoutines to be able to use concurrency.

## Part 1
First, let's use the original Apache Workbench to see the results when testing it with a local NGINX deployment.

1. Installed NGINX
2. Installed Apache benchmark
3. Run:     ```sh
            $ ab -n 5000 -c 500 http://localhost:80/
            ```
    The result is:
    ```sh
    Benchmarking localhost (be patient)
    Completed 500 requests
    Completed 1000 requests
    Completed 1500 requests
    Completed 2000 requests
    Completed 2500 requests
    Completed 3000 requests
    Completed 3500 requests
    Completed 4000 requests
    Completed 4500 requests
    Completed 5000 requests
    Finished 5000 requests
    
    
    Server Software:        nginx/1.18.0
    Server Hostname:        localhost
    Server Port:            80
    
    Document Path:          /
    Document Length:        612 bytes
    
    Concurrency Level:      500
    Time taken for tests:   0.240 seconds
    Complete requests:      5000
    Failed requests:        0
    Total transferred:      4270000 bytes
    HTML transferred:       3060000 bytes
    Requests per second:    20835.24 [#/sec] (mean)
    Time per request:       23.998 [ms] (mean)
    Time per request:       0.048 [ms] (mean, across all concurrent requests)
    Transfer rate:          17376.27 [Kbytes/sec] received
    
    Connection Times (ms)
                  min  mean[+/-sd] median   max
    Connect:        0   10   2.4     10      18
    Processing:     3   13   3.9     13      24
    Waiting:        0    9   3.2      8      18
    Total:         11   23   2.9     23      30
    
    Percentage of the requests served within a certain time (ms)
      50%     23
      66%     24
      75%     25
      80%     25
      90%     27
      95%     27
      98%     28
      99%     28
     100%     30 (longest request)
    ```
    After  testing different levels of concurrency, the only change between the runs is the field "Time per request (mean)".
    This field informs the user about the average amount of time it took for a concurrent group of requests to process, while the same field, but with the parenthesis "(mean, across all concurrent requests)"  tells you the average amount of time it took for a single request to process by itself.
    
    ### Test
    When I processed 500 requests concurrently, it took 24.286ms.
    When I processed them 5 requests concurrently, it would take 0.255ms * 100 = 25.5ms
    "Same number". There is no time savings to performing concurrent requests.
    
    When it comes to the CPU usage, its noticeable that when using concurrency I manage to get to the 100% of usage, while if only using individual requests, without the concurrency, the CPU usage is around 48%. None of the tests have ended with an error, all were succesful.
    
## Part 2
Now, let's try my GoLang implementation of Apache Workbench, and test it with a simple GoLang server.

To build the goAB implementation and the server, run:
```sh
    go build -o goserver server.go
    go build -o goab main.go
```
Then, just start the server with: 
```sh
    ./goserver
```
Start the goAB with the desired parameters (-n X to specify the X number of requests, -c X to specify the X of concurrent requests, and -k for keepAlive):
```sh
    ./goab [parameters] http://localhost:8080/
```

## Part 3
The comparison of speed with ApacheBench 2.3 has not resulted the way I thought it would. AB outperformed the GoAB version I implemented, even though the Go version should outperform it, specially when using more concurrency, as the GoRoutines of GoLang are more optimal than C version of threads.

Goroutines are multiplexed on a limited number of threads, avoiding a lot of context switching overhead, while when using C, there is a change of context that slows down the system, mainly if not used properly.

The results, would probably be better/quicker if my project was coded properly.