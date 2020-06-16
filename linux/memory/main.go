package main

import (
	"fmt"
	"os"
	"flag"
	"time"

	"github.com/mackerelio/go-osstat/memory"
)

func main() {
	scheme := flag.String("scheme", "", "Metric naming scheme, text to prepend to metric.")
	flag.Parse()

	var metrics []string

	now := time.Now()
	timestamp := now.Unix()

	memory, err := memory.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}

	metrics = append(
		metrics,
		fmt.Sprintf("memory.total %d %d\n", memory.Total, timestamp),
		fmt.Sprintf("memory.used %d %d\n", memory.Used, timestamp),
		fmt.Sprintf("memory.cached %d %d\n", memory.Cached, timestamp),
		fmt.Sprintf("memory.free %d %d\n", memory.Free, timestamp),
		fmt.Sprintf("memory.active %d %d\n", memory.Active, timestamp),
		fmt.Sprintf("memory.inactive %d %d\n", memory.Inactive, timestamp),
		fmt.Sprintf("memory.swaptotal %d %d\n", memory.SwapTotal, timestamp),
		fmt.Sprintf("memory.swapused %d %d\n", memory.SwapUsed, timestamp),
		fmt.Sprintf("memory.swapfree %d %d\n", memory.SwapFree, timestamp),
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
