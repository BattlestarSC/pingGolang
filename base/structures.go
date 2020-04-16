package base

import "net"

//Just hold configuration information
type Configuration struct {
	Target  net.IP//target address
	Count   int   //how many to send
	Bytes   int   //how many bytes to pack
	Delay   int   //how many milliseconds between sends
	V6      bool  //do we use ipv6 or not
	Interface string //what interface do we use
}