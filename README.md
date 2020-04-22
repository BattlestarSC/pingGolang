Golang ping program for Cloudflare internship challenge 2020.

## Build
On Debain/Ubuntu distros:

*sudo apt install golang git -y*

On Arch based distros:

*sudo pacman -Syy go git*

Then build with 
* *go get golang.org/x/net/icmp*
* *go build main.go*

Install:

*sudo cp main /bin/goping*

Possibly, if the command does not run, try:
    
    sudo chmod +x /bin/goping
    sudo chmod u+s /bin/goping

## Usage

After installing into your /bin folder:

Usage: *sudo goping flags target*

    note: flags must go before target spec

    note: without setuid, this command must be run as root via sudo

Target can be any of the following:

    -Hostname, like google.com

    -IPv4 address, like 8.8.8.8

    -IPv6 address, like ::1

Flags:

Help

        usage --help
        Display this usage menu
Count

        usage --count <number>
        The number of pings to send before stopping (default: inf)
Delay

        usage --delay "<time spec>"
        The amount of time between pings, specified in a time spec string, such as 1s (default: "1s")
        Allowed time spec endings are ms,s,m,h for milliseconds, seconds, minutes, and hours respectably
        Durations less than 100ms are prohibited
Delta

        usage --stats-delta <number>
        The number of pings between aggregate stats are printed (default: 1) Timeout
        usage --timeout "<time spec>"
        The amount of time allowed before a ping times out (default: "10s")
        Allowed time spec endings are ms,s,m,h for milliseconds, seconds, minutes, and hours respectably
        Durations less than 100ms are prohibited

## Design choices:
- A main function handling user input, metrics, and an interface to the ping go routine
- A ping go routine is an interface to all the backend ping functions
- A dialer/net.Conn is created for every ping sent for a few reasons, chief of which is that it allows me to set a precise timeout on the connection (which the docs do not cover well and is unclear if it is reusable or not) and secondly it allows me better debug information (each new conn is more information if it fails). 
- This was tested with Arch linux in Virtualbox with various example inputs (google.com, 127.0.0.1, ::1, etc)
- While deploying code, I would use more complete library options, such as fastping, since that is often a better way to deploy fast and reliable code, that felt like it was out of bounds for this challenge.
