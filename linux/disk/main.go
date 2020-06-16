package main

import (
	"fmt"
	"os"
	"flag"
	"time"

	"github.com/mackerelio/go-osstat/disk"
)

func main() {
	scheme := flag.String("scheme", "", "Metric naming scheme, text to prepend to metric.")
	flag.Parse()

	var metrics []string

	now := time.Now()
	timestamp := now.Unix()

	disks, err := disk.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	for _, disk := range disks {
		metrics = append(
			metrics,
			fmt.Sprintf("disk.%s.reads %d %d\n", disk.Name, disk.ReadsCompleted, timestamp),
			fmt.Sprintf("disk.%s.writes %d %d\n", disk.Name, disk.WritesCompleted, timestamp),
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
