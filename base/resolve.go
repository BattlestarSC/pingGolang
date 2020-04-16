package base

import (
	"errors"
	"net"
)

//Take a string address and make an IP with it
//Return IP, if it is ipV6 or not, and an error (if any)
func ResolveAddress(address string) (net.IP, bool, error) {
	//convert to address
	ip := net.ParseIP(address)
	//if nil, its not an ip, try to resolve
	var err error
	if ip == nil {
		//attempt to resolve a name
		ips, erro := net.LookupIP(address)
		//if that failed
		if erro != nil {
			err = errors.New("input is not a valid IP address or hostname")
			return nil, false, err
		} else {
			//just attempt the first one
			ip = ips[0]
		}
	} else {
		err = nil
	}
	//check for if it is IPv4 or v6
	addr := ip.To4()
	//if addr is not nil, its ipv4, otherwise its 6
	var v6 bool
	if addr == nil {
		v6 = true
	} else {
		v6 = false
	}

	return ip, v6, err
}
