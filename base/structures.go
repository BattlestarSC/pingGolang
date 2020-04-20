package base

import "time"

//Target for ping
//creation sanity checks addresses
//and determines the connection type
type Target struct {
	Host string
	ConnType string
	V4 bool
}

//Application configuration
type Configuration struct {
	Target  Target
	Delay   time.Duration
	Timeout time.Duration
}

type Response struct {
	Seq      int
	Latency  time.Duration
	Received bool
	Err      error
}