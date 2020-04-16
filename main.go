package main

import (
	"flag"
	"fmt"
	"golang.org/x/net/icmp"
	"log"
	"net"
	"os"
	"ping/base"
)

func main() {
	//get flag arguments
	//count
	countPtr := flag.Int("count", -1, "How many packets to send before terminating (default: inf)")
	//byte count
	bytesPtr := flag.Int("bytes", 16, "How many bytes to send in each packet (default: 16, max allowed: 256)")
	//millisecond delay
	delayPtr := flag.Int("delay", 1000, "How many milliseconds between pings (default: 1000, minimum: 500)")
	//interface to use
	interfacePtr := flag.String("interface", "none", "What interface to use")
	//parse flags
	flag.Parse()
	//get the required non-flag argument
	address := flag.Arg(0)

	//check arguments
	if net.ParseIP(address) == nil {
		_, err := net.LookupIP(address)
		if err != nil {
			fmt.Println("Failed to get usable IP address")
			fmt.Println("Usage: ping <flags> <IP address>")
			os.Exit(1)
		}
	}

	if *countPtr == 0 || *countPtr < -1 {
		fmt.Println("Disallowed packet count value, valid range is 1...2147483647 or -1 for inf")
		os.Exit(1)
	}

	if *bytesPtr > 256 || *bytesPtr < 16 {
		fmt.Println("Disallowed packet size value, valid range is 16...256")
		os.Exit(1)
	}

	if *delayPtr < 500 {
		fmt.Println("Disallowed send speed, valid range is 500...2147483647")
		os.Exit(1)
	}

	//turn the address into an IP
	addr, v6, err := base.ResolveAddress(address)

	if err != nil {
		fmt.Println(err.Error())
		log.Fatal(err.Error())
	}

	//build configuration
	config := base.Configuration{
		Target: addr,
		Count:  *countPtr,
		Bytes:  *bytesPtr,
		Delay:  *delayPtr,
		V6:     v6,
		Interface: *interfacePtr,
	}

	//make the listener
	listener, err := base.CreateListener(config)
	//make sure it worked
	if err != nil {
		fmt.Println(err.Error())
		if config.Interface != "none" {
			fmt.Println("Please double check permissions on interface " + config.Interface)
		}
		log.Fatal(err.Error())
	}
	//remember to close on exit, errors are irrelevant
	defer listener.Close()

	//make a sequence channel
	sequenceChan := make(chan int)
	//init it
	sequenceChan <- 1
	//packet channel
	packetChannel := make(chan icmp.Message)
	//start the pinging go routine
	go base.Pinger(config, listener, sequenceChan, packetChannel)
}
