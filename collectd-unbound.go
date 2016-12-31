package main

import (
	"bufio"
	"bytes"
	"context"
	"log"
	osExec "os/exec"
	"strconv"
	"strings"
	"time"

	"collectd.org/api"
	"collectd.org/exec"
)

func main() {
	e := exec.NewExecutor()
	e.VoidCallback(unboundStats, exec.Interval())
	e.Run(context.Background())
}

func unboundStats(ctx context.Context, interval time.Duration) {
	buf := &bytes.Buffer{}
	cmd := osExec.Command("/bin/sh", "-c", "unbound-control stats")
	cmd.Stdout = buf

	if err := cmd.Run(); err != nil {
		log.Fatalf("unable to execute unbound-control: %v", err)
	}

	now := time.Now()

	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "total.") {
			continue
		}

		fields := strings.SplitN(line, "=", 2)
		if len(fields) != 2 {
			continue
		}

		metric := fields[0]
		metric = strings.TrimPrefix(metric, "total.")
		metric = strings.Replace(metric, ".", "_", -1)

		value, err := strconv.ParseFloat(fields[1], 64)
		if err != nil {
			log.Printf("error: unable to parse metric value: %v", err)
			continue
		}

		vl := &api.ValueList{
			Identifier: api.Identifier{
				Host:         exec.Hostname(),
				Plugin:       "unbound",
				Type:         "gauge",
				TypeInstance: metric,
			},
			Time:     now,
			Interval: interval,
			Values:   []api.Value{api.Gauge(value)},
		}

		exec.Putval.Write(ctx, vl)
	}
}
