package main

import (
	"fmt"
	"os"
	"flag"
	"time"

	"github.com/mackerelio/go-osstat/loadavg"
)

func main() {
	scheme := flag.String("scheme", "", "Metric naming scheme, text to prepend to metric.")
	flag.Parse()

	var metrics []string

	now := time.Now()
	timestamp := now.Unix()

	loadavg, err := loadavg.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}

	metrics = append(
		metrics,
		fmt.Sprintf("loadavg.1m %.2f %d\n", loadavg.Loadavg1, timestamp),
		fmt.Sprintf("loadavg.5m %.2f %d\n", loadavg.Loadavg5, timestamp),
		fmt.Sprintf("loadavg.15m %.2f %d\n", loadavg.Loadavg15, timestamp),
	)

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
