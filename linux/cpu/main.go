package main

import (
	"fmt"
	"os"
	"flag"
	"time"

	"github.com/mackerelio/go-osstat/cpu"
)

func main() {
	scheme := flag.String("scheme", "", "Metric naming scheme, text to prepend to metric.")
	flag.Parse()

	var metrics []string

	now := time.Now()
	timestamp := now.Unix()

	before, err := cpu.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	time.Sleep(time.Duration(1) * time.Second)
	after, err := cpu.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	total := float64(after.Total - before.Total)

	metrics = append(
		metrics,
		fmt.Sprintf("cpu.user %.2f %d\n", float64(after.User-before.User)/total*100, timestamp),
		fmt.Sprintf("cpu.nice %.2f %d\n", float64(after.Nice-before.Nice)/total*100, timestamp),
		fmt.Sprintf("cpu.system %.2f %d\n", float64(after.System-before.System)/total*100, timestamp),
		fmt.Sprintf("cpu.idle %.2f %d\n", float64(after.Idle-before.Idle)/total*100, timestamp),
		fmt.Sprintf("cpu.iowait %.2f %d\n", float64(after.Iowait-before.Iowait)/total*100, timestamp),
		fmt.Sprintf("cpu.irq %.2f %d\n", float64(after.Irq-before.Irq)/total*100, timestamp),
		fmt.Sprintf("cpu.softirq %.2f %d\n", float64(after.Softirq-before.Softirq)/total*100, timestamp),
		fmt.Sprintf("cpu.steal %.2f %d\n", float64(after.Steal-before.Steal)/total*100, timestamp),
		fmt.Sprintf("cpu.guest %.2f %d\n", float64(after.Guest-before.Guest)/total*100, timestamp),
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
