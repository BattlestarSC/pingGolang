package base

import (
	"fmt"
	"net"
	"time"
)

func Connection(data Configuration) (net.Conn, error) {
	//DEBUG
	fmt.Println("DEBUG! Connection func in dailer.go start ")
	//DEBUG

	//create a connection
	connection, err := net.Dial(data.Target.ConnType, data.Target.Host)

	//DEBUG
	fmt.Println("DEBUG! Connection fun in dialer.go connection creation, connection, error: ", connection, err)
	//DEBUG

	//terminate on fail
	if err != nil {
		return nil, err
	}

	//set its deadline
	deadline := time.Now().Add(data.Timeout)
	connection.SetDeadline(deadline)

	//DEBUG
	fmt.Println("DEBUG! Connection fun in dialer.go connection deadline set, ready to return. Connection: ", connection)
	//DEBUG

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
