package main

import (
	"bufio"
	"fmt"
	"os"
	"flag"
	"time"
	"math"
	"strings"
	"strconv"
)

func main() {
	scheme := flag.String("scheme", "", "Metric naming scheme, text to prepend to metric.")
	flag.Parse()

	var metrics []string

	now := time.Now()
	timestamp := now.Unix()

	file, err := os.Open("/proc/stat")
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	cpu_metrics := []string{
		"user",
		"nice",
		"system",
		"idle",
		"iowait",
		"irq",
		"softirq",
		"steal",
		"guest",
		"guest_nice",
	}

    cpu_count := 0

	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		name := fields[0]

		if strings.HasPrefix(name, "cpu") {
			cpu_count += 1

			if name == "cpu" {
				name = "total"
			}

			for i := 1; i < 10; i++ {
				val, err := strconv.ParseUint(fields[i], 10, 64)
				if err != nil {
					fmt.Errorf("failed to parse %s of %s", cpu_metrics[i], name)
				}
				metrics = append(
					metrics,
					fmt.Sprintf("cpu.%s.%s %d %d\n", name, cpu_metrics[i], val, timestamp),
				)
			}
		} else {
			val, err := strconv.ParseUint(fields[1], 10, 64)
			if err != nil {
				fmt.Errorf("failed to parse %s", name)
			}
			metrics = append(
				metrics,
				fmt.Sprintf("cpu.%s %d %d\n", name, val, timestamp),
			)
		}

		if err := scanner.Err(); err != nil {
			fmt.Errorf("scan error for /proc/stat: %s", err)
		}
	}

	// false is number is positive
	if math.Signbit(float64(cpu_count)) == false {
		cpu_count = cpu_count - 1
		metrics = append(
			metrics,
			fmt.Sprintf("cpu.cpu_count %d %d\n", cpu_count, timestamp),
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
