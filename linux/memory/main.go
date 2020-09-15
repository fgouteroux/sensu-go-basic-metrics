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

	file, err := os.Open("/proc/meminfo")
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	mem_metrics := map[string]string{
		"MemTotal": "total",
		"MemFree": "free",
		"Buffers": "buffers",
		"Cached": "cached",
		"SwapTotal": "swapTotal",
		"SwapFree": "swapFree",
		"Dirty": "dirty",
		"MemAvailable": "available",
	}

	mem_val := map[string]int{}

	for scanner.Scan() {
		kv := strings.SplitN(scanner.Text(), ":", 2)
		if len(kv) != 2 {
			continue
		}
		name := kv[0]
		fields := strings.Fields(kv[1])

		if mname, ok := mem_metrics[name]; ok {
			name = mname

			val, err := strconv.ParseUint(fields[0], 10, 64)
			if err != nil {
				fmt.Errorf("failed to parse %s", name)
			}
			mem_val[name] = int(val) * 1024
			metrics = append(
				metrics,
				fmt.Sprintf("memory.%s %d %d\n", name, val * 1024, timestamp),
			)
		}

		if err := scanner.Err(); err != nil {
			fmt.Errorf("scan error for /proc/meminfo: %s", err)
		}
	}

	mem_val["swapUsed"] = mem_val["swapTotal"] - mem_val["swapFree"]
	mem_val["used"] = mem_val["total"] - mem_val["free"]

	if val, ok := mem_val["available"]; ok {
		mem_val["usedWOBuffersCaches"] = mem_val["total"] - val
		mem_val["freeWOBuffersCaches"] = val
	} else {
		mem_val["usedWOBuffersCaches"] = mem_val["used"] - (mem_val["buffers"] + mem_val["cached"])
		mem_val["freeWOBuffersCaches"] = mem_val["free"] + (mem_val["buffers"] + mem_val["cached"])
	}

	metrics = append(
		metrics,
		fmt.Sprintf("memory.swapUsed %d %d\n", mem_val["swapUsed"], timestamp),
		fmt.Sprintf("memory.used %d %d\n", mem_val["used"], timestamp),
		fmt.Sprintf("memory.usedWOBuffersCaches %d %d\n", mem_val["usedWOBuffersCaches"], timestamp),
		fmt.Sprintf("memory.freeWOBuffersCaches %d %d\n", mem_val["freeWOBuffersCaches"], timestamp),
	)

	if mem_val["swapTotal"] > 0 {
		mem_val["swapUsedPercentage"] = 100 * mem_val["swapUsed"] / mem_val["swapTotal"]
		metrics = append(
			metrics,
			fmt.Sprintf("memory.swapUsedPercentage %d %d\n", mem_val["swapUsedPercentage"], timestamp),
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
