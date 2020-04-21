package base

import (
	"time"
)

func Ping(configuration Configuration, output chan Response) {

	//setup
	//sequence number channel
	seq := make(chan int, 1)
	seq <- 1
	//error channel
	errChan := make(chan error, 1)

	for !configuration.Inf && configuration.Count > 0 {

		//first create a channel
		conn, err := Connection(configuration)

		//on error, close the response channel and die
		if err != nil {
			close(output)
			return
		}

		seqId := <-seq

		//now send and receive
		sentTime := sendPing(conn, seqId, configuration, errChan)

		//if no error, continue to read
		err = <-errChan
		if nil != err {

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

		time.Sleep(configuration.Delay)

		configuration.Count--

		//if we are done, kill the main function by closing the channel
		if configuration.Count == 0 {
			close(output)
		}

	}
}
