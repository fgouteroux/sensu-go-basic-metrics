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

func contains(arr []string, str string) bool {
   for _, a := range arr {
      if (a != "" && strings.HasPrefix(str, a)) || a == str {
         return true
      }
   }
   return false
}

func main() {
	scheme := flag.String("scheme", "", "Metric naming scheme, text to prepend to metric.")
	convert := flag.Bool("convert", false, "Convert devicemapper to logical volume name")
	ignore_device := flag.String("ignore-device", "", "Ignore devices matching pattern(s)")
	flag.Parse()

	var metrics []string

	now := time.Now()
	timestamp := now.Unix()

	file, err := os.Open("/proc/diskstats")
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	desc := []string{
		"reads",
		"readsMerged",
		"sectorsRead",
		"readTime",
		"writes",
		"writesMerged",
		"sectorsWritten",
		"writeTime",
		"ioInProgress",
		"ioTime",
		"ioTimeWeighted",
	}

	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 14 {
			continue
		}

		name := fields[2]

		ignore_devices := strings.Split(*ignore_device, ",")
		if contains(ignore_devices, name) {
			continue
		}

		if *convert {
			// Convert devicemapper to logical volume name
			if strings.HasPrefix(name, "dm-") {
				file, err := os.Open(fmt.Sprintf("/sys/block/%s/dm/name", name))
				if err != nil {
					fmt.Print(err)
					os.Exit(1)
				}
				defer file.Close()
				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					name = strings.Fields(scanner.Text())[0]
				}
			}
		}

		for i := 3; i < 14; i++ {
			field := desc[i-3]
			val, err := strconv.ParseUint(fields[i], 10, 64)
			if err != nil {
				fmt.Errorf("failed to parse %s of %s", field, name)
			}
			metrics = append(
				metrics,
				fmt.Sprintf("disk_metrics.%s.%s %d %d\n", name, field, val, timestamp),
			)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Errorf("scan error for /proc/diskstats: %s", err)
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
