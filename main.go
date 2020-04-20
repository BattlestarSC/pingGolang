package main

import (
	"flag"
	"fmt"
	"os"
	"pingGolang/base"
	"strconv"
	"time"
)

func main() {
	flag.Parse()
	address := flag.Arg(0)
	result, err := base.CreateTarget(address)
	//DEBUG
	fmt.Println("DEBUG! Create target result and error=> ", result, err)
	//DEBUG
	if err != nil {
		fmt.Println("Unable to find an address to ping, please double check address or usage")
		fmt.Println("Usage: sudo ping <address or hostname> <time between pings, in seconds, optional>")
		os.Exit(1)
	}
	delay := flag.Arg(1)
	//DEBUG
	fmt.Println("DEBUG! Delay flag ", delay)
	//DEBUG
	if delay == "" {
		delay = "1s"
	} else {
		t, err := strconv.Atoi(delay)
		if err != nil {
			fmt.Println("Invalid delay specification, need delay time in seconds, ex: 2")
			fmt.Println("Usage: sudo ping <address or hostname> <time between pings, in seconds, optional>")
			os.Exit(1)
		}
		delay = strconv.Itoa(t) + "s"
	}
	configDelay, err := time.ParseDuration(delay)
	//DEBUG
	fmt.Println("DEBUG! Delay parse => ", configDelay, err)
	//DEBUG
	if err != nil {
		fmt.Println("Invalid delay specification, need delay time in seconds, ex: 2")
		fmt.Println("Usage: sudo ping <address or hostname> <time between pings, in seconds, optional>")
		os.Exit(1)
	}

	//DEBUG TIMEOUT
	tmout, _ := time.ParseDuration("5s")

	config := base.Configuration{
		Target:  result,
		Delay:   configDelay,
		Timeout: tmout,
		//DEBUG count
		Count : 5,
		Inf: false,
	}
	//DEBUG
	fmt.Println("DEBUG! Config creation => ", config)
	//DEBUG
	output := make(chan base.Response)
	//DEBUG
	fmt.Println("DEBUG! Output channel created")
	//DEBUG
	go base.Ping(config, output)
	//DEBUG
	fmt.Println("DEBUG! After started base.Ping")
	//DEBUG

	var total int = 0
	var avgTi int64 = 0
	var recv int = 0

	for {
		//DEBUG
		fmt.Println("DEBUG! Begin infinite loop")
		//DEBUG
		resp, open := <-output

		if !open {
			fmt.Println("Encountered terminating error")
			agStats(total, avgTi, recv)
			return
		}

		seq := resp.Seq
		total++
		//DEBUG
		fmt.Println("DEBUG! Got response: ", resp)
		//DEBUG
		if resp.Received {
			tim := resp.Latency.Nanoseconds() / 1000000
			avgTi += tim
			recv++
			fmt.Println("Response number " + strconv.Itoa(seq) + " from " + config.Target.Host + " in " + strconv.FormatInt(tim, 10) + "ms")
		} else {
			if resp.Err == nil {
				//this should never happen, but apparently it does somehow
				fmt.Println("No response for ping number " + strconv.Itoa(seq) + " ICMP timed out")
			} else {
				fmt.Println("No response for ping number " + strconv.Itoa(seq) + " with error " + resp.Err.Error())
			}
		}
	}
}

func agStats(total int, averageTime int64, numberRecieved int) {
	var percRecv float64
	var avgTimeResult int64
	if numberRecieved > 0 {
		percRecv = (1 - (float64(numberRecieved) / float64(total))) * 100
		avgTimeMs := averageTime / 1000000
		avgTimeResult = avgTimeMs / int64(numberRecieved)
	} else {
		percRecv = 0
		avgTimeResult = 0
	}
	fmt.Println("Aggregate stats: ")
	fmt.Println(total, " pings sent ", numberRecieved, " received for ", percRecv, "percent loss")
	fmt.Println(avgTimeResult, "ms average latency")
}
