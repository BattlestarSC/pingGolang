package base

import (
	"fmt"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
	"net"
	"time"
)

func sendPing(conn net.Conn, seq int, configuration Configuration, errChan chan error) time.Time {

	//DEBUG
	fmt.Println("DEBUG! sendPing func in send.go start")
	//DEBUG

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

	//DEBUG
	fmt.Println("DEBUG! sendPing func in send.go configuration type and config data: ", typ, configuration)
	//DEBUG

	//create packet
	packet, err = (&icmp.Message{
		Type: typ,
		Code: 0,
		Body: &icmp.Echo{
			Seq: seq,
		},
	}).Marshal(nil)

	//DEBUG
	fmt.Println("DEBUG! sendPing func in send.go packet created: ", packet, err)
	//DEBUG

	//if error, give up
	if err != nil {

		//DEBUG
		fmt.Println("DEBUG! sendPing func in send.go packet creation failed, exiting with error, ", err.Error())
		//DEBUG

		errChan <- err
		close(errChan)
		return time.Now()
	}

	_, err = conn.Write(packet)

	//DEBUG
	fmt.Println("DEBUG! sendPing func in send.go wrote packet to connection, error was: ", err)
	//DEBUG

	//if error, give up
	if err != nil {

		//DEBUG
		fmt.Println("DEBUG! sendPing func in send.go send error encountered, terminating with error ", err)
		//DEBUG

		errChan <- err
		close(errChan)
		return time.Now()
	}

	//DEBUG
	fmt.Println("DEBUG! sendPing func in send.go sent correctly, exiting")
	//DEBUG

	errChan <- nil
	return time.Now()
}