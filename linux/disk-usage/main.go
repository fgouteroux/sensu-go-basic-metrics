package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
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
	ignore_device := flag.String("ignore-device", "", "Ignore devices matching pattern(s)")
	local := flag.Bool("local", false, "Only check local filesystems (df -l option)")
	block_size := flag.String("block-size", "M", "Set block size for sizes printed")
	flag.Parse()

	var metrics []string

	now := time.Now()
	timestamp := now.Unix()
	sanitizedChars := strings.NewReplacer("/", "_")


	args := []string{"-P", fmt.Sprintf("-B%s", *block_size)}

	if *local {
		args = append(args, fmt.Sprintf("-l"))
	}
	out, err := exec.Command("df", args...).Output()

    if err != nil {
    	fmt.Print(err)
		os.Exit(1)
    }

    scanner := bufio.NewScanner(strings.NewReader(string(out)))
    // skip the first line
    // Filesystem           1024-blocks     Used Available Capacity Mounted on
    scanner.Scan()

	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())

		name := fields[5]

		if name == "/" {
			name = "root"
		} else {
			name = sanitizedChars.Replace(strings.TrimPrefix(name, "/"))
		}

		ignore_devices := strings.Split(*ignore_device, ",")
		if contains(ignore_devices, name) {
			continue
		}

		used, err := strconv.ParseUint(strings.TrimSuffix(fields[2], *block_size), 10, 64)
		if err != nil {
			fmt.Errorf("failed to parse used of %s", name)
		}

		avail, err := strconv.ParseUint(strings.TrimSuffix(fields[3], *block_size), 10, 64)
		if err != nil {
			fmt.Errorf("failed to parse avail of %s", name)
		}

		used_percentage, err := strconv.ParseUint(strings.TrimSuffix(fields[4], "%"), 10, 64)
		if err != nil {
			fmt.Errorf("failed to parse used_percentage of %s", name)
		}

		metrics = append(
			metrics,
			fmt.Sprintf("disk_usage.%s.used %d %d\n", name, used, timestamp),
			fmt.Sprintf("disk_usage.%s.avail %d %d\n", name, avail, timestamp),
			fmt.Sprintf("disk_usage.%s.used_percentage %d %d\n", name, used_percentage, timestamp),
		)
	}

	if err := scanner.Err(); err != nil {
		fmt.Errorf("scan error for df: %s", err)
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
