package main

import (
	"fmt"
	"os"
	"flag"
	"time"

	"github.com/mackerelio/go-osstat/network"
)

func main() {
	scheme := flag.String("scheme", "", "Metric naming scheme, text to prepend to metric.")
	flag.Parse()

	var metrics []string

	now := time.Now()
	timestamp := now.Unix()

	interfaces, err := network.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	for _, inet := range interfaces {
		metrics = append(
			metrics,
			fmt.Sprintf("interface.%s.rx_bytes %d %d\n", inet.Name, inet.RxBytes, timestamp),
			fmt.Sprintf("interface.%s.tx_bytes %d %d\n", inet.Name, inet.TxBytes, timestamp),
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
