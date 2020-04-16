package base

import (
	"crypto/rand"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
	"log"
	"net"
	"os"
	"time"
)

func makePacket(configuration Configuration, seq chan int) ([]byte,icmp.Message,error) {
	var bytes []byte
	var err error
	var message icmp.Message
	//get sequence number
	s := <- seq
	//inc
	seq <- s + 1
	//make random data
	data := make([]byte, configuration.Bytes)
	rand.Read(data)
	//if v4/v6
	if configuration.V6 {
		//build v6 message
		message = icmp.Message{
			Type:     ipv6.ICMPTypeEchoRequest,
			Code:     0,
			Checksum: 0,
			Body:     &icmp.Echo{
				ID: os.Getpid() & 0xffff,
				Seq: s,
				Data: data,
			},
		}
		//convert to bytes
		bytes, err = message.Marshal(nil)
	} else {
		//build v4 message
		message = icmp.Message{
			Type:     ipv4.ICMPTypeEcho,
			Code:     0,
			Checksum: 0,
			Body:     &icmp.Echo{
				ID: os.Getpid() & 0xffff,
				Seq: s,
				Data: data,
			},
		}
		//convert to bytes
		bytes, err = message.Marshal(nil)
	}

	//if error
	if err != nil {
		return nil,icmp.Message{},err
	}

	//otherwise make the message
	return bytes,message,nil
}

func Pinger(config Configuration, conn *icmp.PacketConn, sequenceChannel chan int, packetChannel chan icmp.Message) {
	for {
		packet,message,err := makePacket(config, sequenceChannel)
		if err != nil {
			close(sequenceChannel)
			log.Fatal(err.Error())
		}
		var sendError error
		if config.Interface == "none" {
			_, sendError = conn.WriteTo(packet, &net.UDPAddr{IP: config.Target})
		} else {
			_, sendError = conn.WriteTo(packet, &net.UDPAddr{
				IP:   config.Target,
				Zone: config.Interface,
			})
		}
		if sendError != nil {
			close(sequenceChannel)
			log.Fatal(sendError.Error())
		}
		packetChannel <- message
		//delay milliseconds to nanoseconds then sleep
		time.Sleep(time.Duration(config.Delay * 1000000))
	}
}