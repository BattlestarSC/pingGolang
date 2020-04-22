package base

import (
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
	"net"
	"time"
)

func recieve(conn net.Conn, seq int, timeSend time.Time, configuration Configuration, out chan Response) {

	//get a response
	response := make([]byte,1500)
	_, err := conn.Read(response)

	//if err, assume it timed out
	if err != nil {
		
		out <- Response{
			Seq:      seq,
			Latency:  -1,
			Received: false,
			Err:      err,
		}
		return
	}

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

	if err != nil {

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

	out <- Response{
		Seq:      seq,
		Latency:  lat,
		Received: true,
		Err:      nil,
	}

}
