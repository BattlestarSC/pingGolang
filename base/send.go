package base

import (
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
	"net"
	"time"
)

func sendPing(conn net.Conn, seq int, configuration Configuration, errChan chan error) time.Time {
	//create packet
	var packet []byte
	var typ icmp.Type
	var err error

	//figure out type
	if configuration.Target.V4 {
		typ = ipv4.ICMPTypeEcho
	} else {
		typ = ipv6.ICMPTypeEchoRequest
	}

	//create packet
	packet, err = (&icmp.Message{
		Type: typ,
		Code: 0,
		Body: &icmp.Echo{
			Seq: seq,
		},
	}).Marshal(nil)

	//if error, give up
	if err != nil {

		errChan <- err
		close(errChan)
		return time.Now()
	}

	_, err = conn.Write(packet)

	//if error, give up
	if err != nil {

		errChan <- err
		close(errChan)
		return time.Now()
	}

	errChan <- nil
	return time.Now()
}