package base

import (
	"fmt"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
	"net"
	"time"
)

func recieve(conn net.Conn, seq int, timeSend time.Time, configuration Configuration, out chan Response) {

	//DEBUG
	fmt.Println("DEBUG! recieve func in recieve.go started")
	//DEBUG

	//get a response
	response := make([]byte,1500)
	_, err := conn.Read(response)

	//DEBUG
	fmt.Println("DEBUG! recieve func in recieve.go connection read error: ", err)
	//DEBUG

	//if err, assume it timed out
	if err != nil {
		//DEBUG
		fmt.Println("DEBUG! recieve func in recieve.go connection read error terminated")
		//DEBUG
		out <- Response{
			Seq:      seq,
			Latency:  -1,
			Received: false,
			Err:      err,
		}
		return
	}

	//DEBUG
	fmt.Println("DEBUG! recieve func in recieve.go response was valid: ", response)
	//DEBUG

	//otherwise its probably valid without an error
	//this is according to the ONE example in the docs
	//and the x extensions seem to offer no way to read
	//the data properly and check the seq number
	//so we will just take the docs lack of a word for it
	if configuration.Target.V4 {
		_, err = icmp.ParseMessage(ipv4.ICMPTypeEchoReply.Protocol(), response)
	} else {
		_, err = icmp.ParseMessage(ipv6.ICMPTypeEchoReply.Protocol(), response)
	}

	//DEBUG
	fmt.Println("DEBUG! recieve func in recieve.go Parse message error result: ", err)
	//DEBUG

	if err != nil {

		//DEBUG
		fmt.Println("DEBUG! recieve func in recieve.go dying on parsing recieved message")
		//DEBUG

		out <- Response{
			Seq:      seq,
			Latency:  -1,
			Received: false,
			Err:      err,
		}
		return
	}

	//compute latency
	lat := time.Since(timeSend)

	//DEBUG
	fmt.Println("DEBUG! recieve func in recieve.go completed and calculated latency of ", timeSend.String())
	//DEBUG

	out <- Response{
		Seq:      seq,
		Latency:  lat,
		Received: true,
		Err:      nil,
	}

	//DEBUG
	fmt.Println("DEBUG! recieve func in recieve.go returned successfully")
	//DEBUG

}