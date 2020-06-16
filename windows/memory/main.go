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
		os.Exit(1)
	}

	metrics = append(
		metrics,
		fmt.Sprintf("memory.total %d %d\n", memory.Total, timestamp),
		fmt.Sprintf("memory.used %d %d\n", memory.Used, timestamp),
		fmt.Sprintf("memory.free %d %d\n", memory.Free, timestamp),
		fmt.Sprintf("memory.page_file_total %d %d\n", memory.PageFileTotal, timestamp),
		fmt.Sprintf("memory.page_file_free %d %d\n", memory.PageFileFree, timestamp),
		fmt.Sprintf("memory.virtual_total %d %d\n", memory.VirtualTotal, timestamp),
		fmt.Sprintf("memory.virtual_free %d %d\n", memory.VirtualFree, timestamp),
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
