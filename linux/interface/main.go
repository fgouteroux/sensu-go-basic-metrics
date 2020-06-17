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

	file, err := os.Open("/proc/net/dev")
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	desc := []string{
		"bytes",
		"packets",
		"errs",
		"drop",
		"fifo",
		"frame",
		"compressed",
		"multicast",
		"bytes",
		"packets",
		"errs",
		"drop",
		"fifo",
		"colls",
		"carrier",
		"compressed",
	}

	for scanner.Scan() {
		// Reference: dev_seq_printf_stats in Linux source code
		kv := strings.SplitN(scanner.Text(), ":", 2)
		if len(kv) != 2 {
			continue
		}
		fields := strings.Fields(kv[1])
		if len(fields) < 16 {
			continue
		}
		iface := strings.TrimSpace(kv[0])

		// Parse receive
		for i := 0; i < 8; i++ {
			val, err := strconv.ParseUint(fields[i], 10, 64)
			if err != nil {
				fmt.Errorf("failed to parse rxFields of %s", iface)
			}
			metrics = append(
				metrics,
				fmt.Sprintf("interface.%s.rx%s %d %d\n", iface, strings.Title(desc[i]), val, timestamp),
			)
		}

		// Parse transmit
		for i := 8; i < 16; i++ {
			val, err := strconv.ParseUint(fields[i], 10, 64)
			if err != nil {
				fmt.Errorf("failed to parse txFields of %s", iface)
			}
			metrics = append(
				metrics,
				fmt.Sprintf("interface.%s.tx%s %d %d\n", iface, strings.Title(desc[i]), val, timestamp),
			)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Errorf("scan error for /proc/net/dev: %s", err)
		os.Exit(1)
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
