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

	//create flags
	countFlag := flag.Int("count", 0, "The number of pings to send before stopping (default: inf)")
	delayFlag := flag.String("delay", "1s", "The amount of time between pings, specified in a time spec string, such as 1s (default: 1s)")
	deltaFlag := flag.Int("stats-delta", 1, "The number of pings between aggregate stats are printed (default: 1)")
	timeoutFlag := flag.String("timeout", "10s", "The amount of time allowed before a ping times out (default: 10s)")
	helpFlag := flag.Bool("help", false, "Display help menu")

	flag.Parse()

	if *helpFlag {
		usage()
		os.Exit(0)
	}

	address := os.Args[len(os.Args) - 1]
	result, err := base.CreateTarget(address)

	if err != nil {
		fmt.Println("Unable to find a suitable target, got error: ", err.Error())
		usage()
		os.Exit(1)
	}

	configDelay, err := time.ParseDuration(*delayFlag)

	if err != nil {
		fmt.Println("Invalid delay specification")
		usage()
		os.Exit(1)
	}

	minimumDelay,_ := time.ParseDuration("100ms")

	if configDelay.Nanoseconds() < minimumDelay.Nanoseconds() {
		fmt.Println("Invalid delay specification, under 100ms")
		usage()
		os.Exit(1)
	}

	if *countFlag < 0 {
		fmt.Println("Negative count values are not permitted")
		usage()
		os.Exit(1)
	}

	var inf bool
	inf = *countFlag == 0

	tmout, err := time.ParseDuration(*timeoutFlag)

	if err != nil {
		fmt.Println("Invalid timeout time spec")
		usage()
		os.Exit(1)
	}

	if tmout.Nanoseconds() < minimumDelay.Nanoseconds() {
		fmt.Println("Invalid timeout specification, under 100ms")
		usage()
		os.Exit(1)
	}

	if *deltaFlag < 1 {
		fmt.Println("Invalid delta value")
		usage()
		os.Exit(1)
	}

	config := base.Configuration{
		Target:  result,
		Delay:   configDelay,
		Timeout: tmout,
		Count : *countFlag,
		Inf: inf,
	}

	output := make(chan base.Response, 1)

	go base.Ping(config, output)


	var total int = 0
	avgTi, _ := time.ParseDuration("0s")
	var recv int = 0

	for {

		resp, open := <-output

		if !open {
			fmt.Println("Encountered terminating error")
			agStats(total, avgTi, recv)
			return
		}

		seq := resp.Seq
		total++

		if resp.Received {
			avgTi += resp.Latency
			recv++
			fmt.Println("Response number " + strconv.Itoa(seq) + " from " + config.Target.Host + " in " + strconv.FormatInt(resp.Latency.Milliseconds(), 10) + "ms")
		} else {
			if resp.Err == nil {
				//this should never happen, but apparently it does somehow
				fmt.Println("No response for ping number " + strconv.Itoa(seq) + " unexpected unknown error")
			} else {
				fmt.Println("No response for ping number " + strconv.Itoa(seq) + " with error " + resp.Err.Error())
			}
		}
		if seq % *deltaFlag == 0 {
			agStats(total, avgTi, recv)
		}
	}
}

func agStats(total int, averageTime time.Duration, numberRecieved int) {
	var percRecv float64
	var avgTimeResult int64
	if numberRecieved > 0 {
		percRecv = (1 - (float64(numberRecieved) / float64(total))) * 100
		avgTimeResult = averageTime.Milliseconds() / int64(total)
	} else {
		percRecv = 0
		avgTimeResult = 0
	}
	fmt.Print("Aggregate stats: ")
	fmt.Print(total, " pings sent ", numberRecieved, " received for ", percRecv, " percent loss ")
	fmt.Println(avgTimeResult, "ms average latency")
}

func usage() {
	fmt.Println()
	fmt.Print("Usage: ")
	fmt.Println(os.Args[0], " <flags> target")
	fmt.Println("note: flags must go before target spec")
	fmt.Println()
	fmt.Println("Target can be any of the following:")
	fmt.Println("\t-Hostname, like google.com")
	fmt.Println("\t-IPv4 address, like 8.8.8.8")
	fmt.Println("\t-IPv6 address, like ::1")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println()
	fmt.Println("Help")
	fmt.Println("\tusage --help")
	fmt.Println("\tDisplay this usage menu")
	fmt.Println("Count")
	fmt.Println("\tusage --count <number>")
	fmt.Println("\tThe number of pings to send before stopping (default: inf)")
	fmt.Println("Delay")
	fmt.Println("\tusage --delay \"<time spec>\"")
	fmt.Println("\tThe amount of time between pings, specified in a time spec string, such as 1s (default: \"1s\")")
	fmt.Println("\tAllowed time spec endings are ms,s,m,h for milliseconds, seconds, minutes, and hours respectably")
	fmt.Println("\tDurations less than 100ms are prohibited")
	fmt.Println("Delta")
	fmt.Println("\tusage --stats-delta <number>")
	fmt.Println("The number of pings between aggregate stats are printed (default: 1)")
	fmt.Println("Timeout")
	fmt.Println("\tusage --timeout \"<time spec>\"")
	fmt.Println("\tThe amount of time allowed before a ping times out (default: \"10s\")")
	fmt.Println("\tAllowed time spec endings are ms,s,m,h for milliseconds, seconds, minutes, and hours respectably")
	fmt.Println("\tDurations less than 100ms are prohibited")
	fmt.Println()
}