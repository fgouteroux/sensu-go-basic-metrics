package main

import (
	"bufio"
	"fmt"
	"os"
	"flag"
	"time"
	"strings"
	"strconv"
)

func main() {
	scheme := flag.String("scheme", "", "Metric naming scheme, text to prepend to metric.")
	flag.Parse()

	var metrics []string

	now := time.Now()
	timestamp := now.Unix()

	file, err := os.Open("/proc/vmstat")
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		metric := strings.Split(line, " ")

		val, err := strconv.ParseUint(metric[1], 10, 64)
		if err != nil {
			fmt.Errorf("failed to scan %s from /proc/vmstat", metric[1])
		}

		metrics = append(
		metrics,
		fmt.Sprintf("vmstat.%s %d %d\n", metric[0], val, timestamp),
		)
			
		if err := scanner.Err(); err != nil {
			fmt.Errorf("scan error for /proc/vmstat: %s", err)
		}
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
