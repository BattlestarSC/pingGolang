package base

import (
	"fmt"
	"time"
)

func Ping(configuration Configuration, output chan Response) {

	//DEBUG
	fmt.Println("DEBUG! Ping function in ping.go starting")
	//DEBUG

	//setup
	//sequence number channel
	seq := make(chan int, 1)
	seq <- 1
	//error channel
	errChan := make(chan error, 1)

	//DEBUG
	fmt.Println("DEBUG! Ping function in ping.go init channels created")
	//DEBUG

	for !configuration.Inf && configuration.Count > 0 {

		//DEBUG
		fmt.Println("DEBUG! Ping function in ping.go started inf loop")
		//DEBUG

		//first create a channel
		conn, err := Connection(configuration)

		//DEBUG
		fmt.Println("DEBUG! Ping function in ping.go created connection, results: ", conn, err)
		//DEBUG

		//on error, close the response channel and die
		if err != nil {

			//DEBUG
			fmt.Println("DEBUG! Ping function in ping.go closing channel and dying on error from connection creation: ", err.Error())
			//DEBUG

			close(output)
			return
		}

		seqId := <-seq

		//DEBUG
		fmt.Println("DEBUG! Ping function in ping.go got seq ", seqId)
		//DEBUG

		//now send and receive
		sentTime := sendPing(conn, seqId, configuration, errChan)

		//DEBUG
		fmt.Println("DEBUG! Ping function in ping.go sent message, got time of ", sentTime.String())
		//DEBUG

		//if no error, continue to read
		err = <-errChan
		if nil != err {

			//DEBUG
			fmt.Println("DEBUG! Ping function in ping.go got error from send function, dying with error: ", err.Error())
			//DEBUG

			//otherwise die
			output <- Response{
				Seq:      seqId,
				Latency:  0,
				Received: false,
				Err:      err,
			}
			close(output)
			return
		}

		seq <- seqId + 1

		recieve(conn, seqId, sentTime, configuration, output)

		//DEBUG
		fmt.Println("DEBUG! Ping function in ping.go recieved message, pushed to main thread")
		//DEBUG

		time.Sleep(configuration.Delay)

		//DEBUG
		fmt.Println("DEBUG! Ping function in ping.go passed sleep of " + configuration.Delay.String())
		//DEBUG

		configuration.Count--

		//DEBUG
		fmt.Println("DEBUG! Ping function in ping.go reduced count to ", configuration.Count)
		//DEBUG

		//if we are done, kill the main function by closing the channel
		if configuration.Count == 0 {
			close(output)
		}

	}
}
