package base

import (
	"errors"
	"net"
)

//create a target structure and check the address
func CreateTarget(address string) (Target, error) {
	var addr   net.IP
	var output Target
	var addresses []net.IP
	//first, parse the address
	addr = net.ParseIP(address)
	//if that fails
	if addr == nil {
		//attempt to resolve the name
		addresses, _ = net.LookupIP(address)
		//check length
		if len(addresses) < 1 {
			//this failed to find a suitable address
			err := "Unable to find a suitable target address for hostname: " + address
			return output, errors.New(err)
		} else {
			addr = addresses[0]
		}
	}

	//now that we have a valid ip, lets check for v4 or 6
	if addr != nil {
		if addr.To4() == nil {
			output.ConnType = "ip6:ipv6-icmp"
			output.V4 = false
		} else {
			output.ConnType = "ip4:icmp"
			output.V4 = true
		}
	} else {
		//This should never happen, but just in case
		//assume ipv6
		output.ConnType = "ip6:ipv6-icmp"
		output.V4 = false
	}
	//and load the original address into the struct
	//now that it has been proven to be accurate
	output.Host = address
	return output, nil
}
