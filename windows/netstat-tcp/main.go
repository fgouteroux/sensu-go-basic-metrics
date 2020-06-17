package main

import (
	"fmt"
	"os"
	"flag"
	"time"
	"strings"

	"github.com/cakturk/go-netstat/netstat"
)

func main() {
	scheme := flag.String("scheme", "", "Metric naming scheme, text to prepend to metric.")
	flag.Parse()

	var metrics []string

	states := map[string]int{
		"UNKNOWN": 0,
		"ESTABLISHED": 0,
		"SYN_SENT": 0,
		"SYN_RECV": 0,
		"FIN_WAIT1": 0,
		"FIN_WAIT2": 0,
		"TIME_WAIT": 0,
		"CLOSE": 0, 
		"CLOSE_WAIT": 0,
		"LAST_ACK": 0,
		"LISTEN": 0,
		"CLOSING": 0,
	}

	now := time.Now()
	timestamp := now.Unix()

	// list all the TCP sockets
	socks, err := netstat.TCPSocks(netstat.NoopFilter)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	// count TCP sockets states
	for _, sock := range socks {
		state := strings.Fields(fmt.Sprintf("%v", sock))[3]
		if state == "" {
			states["CLOSE"] = states["CLOSE"]+1
		} else {
			states[state] = states[state]+1
		}
	}

	for k, v := range states {
		metrics = append(
			metrics,
			fmt.Sprintf("net.%s %d %d\n", k, v, timestamp),
		)
	}

	if *scheme != "" {
		for _, metric := range metrics {
			fmt.Printf("%s.%s", *scheme, metric)
		}
	} else {
		for _, metric := range metrics {
			fmt.Print(metric)
		}
	}
}
