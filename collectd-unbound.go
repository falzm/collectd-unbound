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

var (
	metrics = map[string]string{
		"num_queries":              "derive",
		"num_cachehits":            "derive",
		"num_cachemiss":            "derive",
		"num_prefetch":             "derive",
		"num_zero_ttl":             "derive",
		"num_recursivereplies":     "derive",
		"requestlist_avg":          "gauge",
		"requestlist_max":          "derive",
		"requestlist_overwritten":  "derive",
		"requestlist_exceeded":     "derive",
		"requestlist_current_all":  "derive",
		"requestlist_current_user": "derive",
		"recursion_time_avg":       "gauge",
		"recursion_time_median":    "gauge",
	}
)

func main() {
	e := exec.NewExecutor()
	e.VoidCallback(unboundStats, exec.Interval())
	e.Run(context.Background())
}

func unboundStats(ctx context.Context, interval time.Duration) {
	var vl *api.ValueList

	buf := &bytes.Buffer{}
	cmd := osExec.Command("/bin/sh", "-c", "unbound-control stats_noreset")
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

		if _, ok := metrics[metric]; !ok {
			continue
		}

		value, err := strconv.ParseFloat(fields[1], 64)
		if err != nil {
			log.Printf("error: unable to parse metric value: %v", err)
			continue
		}

		switch metrics[metric] {
		case "derive":
			vl = &api.ValueList{
				Identifier: api.Identifier{
					Host:         exec.Hostname(),
					Plugin:       "unbound",
					Type:         "derive",
					TypeInstance: metric,
				},
				Time:     now,
				Interval: interval,
				Values:   []api.Value{api.Derive(value)},
			}

		case "gauge":
			vl = &api.ValueList{
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
		}

		exec.Putval.Write(ctx, vl)
	}
}
