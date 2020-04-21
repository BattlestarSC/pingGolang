package base

import (
	"net"
	"time"
)

func Connection(data Configuration) (net.Conn, error) {
	//create a connection
	connection, err := net.Dial(data.Target.ConnType, data.Target.Host)

	//terminate on fail
	if err != nil {
		return nil, err
	}

	//set its deadline
	deadline := time.Now().Add(data.Timeout)
	connection.SetDeadline(deadline)

	//return the connection
	return connection, nil

	/*
	//create a one time connection with a timeout value
	var dialer net.Dialer
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	//create the connection
	connection, err := dialer.DialContext(ctx,
		data.target.ConnType,
		data.target.Host)

	//failed to connect
	if err != nil {
		return nil, err
	}

	//return the new connection
	return connection, nil

	 */
}
