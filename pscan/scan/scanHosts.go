package scan

import (
	"fmt"
	"net"
	"time"
)

// define PortState to represent the state of a single TCP port
type PortState struct {
	Port int
	Open state
}

// define the scan results for a single host
type Results struct {
	Host       string
	NotFound   bool
	PortStates []PortState
}

type state bool

// convert the boolean value of state to a human readable string
func (s state) String() string {
	if s {
		return "open"
	}

	return "closed"
}

// scanPort performs a port scan on a single TCP port
func scanPort(host string, port int) PortState {
	p := PortState{
		Port: port,
	}

	// define the address of the target
	address := net.JoinHostPort(host, fmt.Sprintf("%d", port))

	// perform a connection attempt
	scanConn, err := net.DialTimeout("tcp", address, 5*time.Second)
	if err != nil {
		return p
	}

	scanConn.Close()
	p.Open = true
	return p
}

// define the Run function
func Run(h *HostsList, ports []int) []Results {
	res := make([]Results, 0, len(h.Hosts))

	for _, host := range h.Hosts {
		r := Results{
			Host: host,
		}

		if _, err := net.LookupHost(host); err != nil {
			r.NotFound = true
			res = append(res, r)
			continue
		}

		for _, port := range ports {
			r.PortStates = append(r.PortStates, scanPort(host, port))
		}

		res = append(res, r)
	}

	return res
}
