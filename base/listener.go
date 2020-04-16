package base

import (
	"golang.org/x/net/icmp"
)

func CreateListener(configuration Configuration) (*icmp.PacketConn, error) {
	//select v4/v6
	var networkString string
	//listen address
	var listenAddress string
	if configuration.V6 {
		networkString = "udp6"
		listenAddress = "::"
	} else {
		networkString = "udp4"
		listenAddress = "0.0.0.0"
	}
	if configuration.Interface != "none" {
		networkString += "%" + configuration.Interface
	}
	//make listener
	listener, err := icmp.ListenPacket(networkString, listenAddress)

	//return error and listener
	return listener, err
}