package main

import (
	"bufio"
	"bytes"
	"context"
	"log"
	osExec "os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"collectd.org/api"
	"collectd.org/exec"
)

var (
	metrics = map[string]string{
		"*.num.queries":                 "derive",
		"*.num.queries_ip_ratelimited":  "derive",
		"*.num.cachehits":               "derive",
		"*.num.cachemiss":               "derive",
		"*.num.prefetch":                "derive",
		"*.num.zero_ttl":                "derive",
		"*.num.recursivereplies":        "derive",
		"*.num.request":                 "derive",
		"*.requestlist.avg":             "gauge",
		"*.requestlist.max":             "derive",
		"*.requestlist.overwritten":     "derive",
		"*.requestlist.exceeded":        "derive",
		"*.requestlist.current.all":     "derive",
		"*.requestlist.current.user":    "derive",
		"*.recursion.time.avg":          "gauge",
		"*.recursion.time.median":       "gauge",
		"*.tcpusage":                    "gauge",
		"mem.cache.rrset":               "gauge",
		"mem.cache.message":             "gauge",
		"mem.mod.iterator":              "gauge",
		"mem.mod.validator":             "gauge",
		"mem.mod.respip":                "gauge",
		"mem.streamwait":                "gauge",
		"num.query.type.A":              "derive",
		"num.query.type.SOA":            "derive",
		"num.query.type.PTR":            "derive",
		"num.query.type.TXT":            "derive",
		"num.query.type.AAAA":           "derive",
		"num.query.type.SRV":            "derive",
		"num.query.class.IN":            "derive",
		"num.query.opcode.QUERY":        "derive",
		"num.query.tcp":                 "derive",
		"num.query.tcpout":              "derive",
		"num.query.tls":                 "derive",
		"num.query.tls.resume":          "derive",
		"num.query.ipv6":                "derive",
		"num.query.flags.QR":            "derive",
		"num.query.flags.AA":            "derive",
		"num.query.flags.TC":            "derive",
		"num.query.flags.RD":            "derive",
		"num.query.flags.RA":            "derive",
		"num.query.flags.Z":             "derive",
		"num.query.flags.AD":            "derive",
		"num.query.flags.CD":            "derive",
		"num.query.edns.present":        "derive",
		"num.query.edns.DO":             "derive",
		"num.answer.rcode.NOERROR":      "derive",
		"num.answer.rcode.FORMERR":      "derive",
		"num.answer.rcode.SERVFAIL":     "derive",
		"num.answer.rcode.NXDOMAIN":     "derive",
		"num.answer.rcode.NOTIMPL":      "derive",
		"num.answer.rcode.REFUSED":      "derive",
		"num.answer.rcode.nodata":       "derive",
		"num.query.ratelimited":         "derive",
		"num.answer.secure":             "derive",
		"num.answer.bogus":              "derive",
		"num.rrset.bogus":               "derive",
		"num.query.aggressive.NOERROR":  "derive",
		"num.query.aggressive.NXDOMAIN": "derive",
		"unwanted.queries":              "derive",
		"unwanted.replies":              "derive",
		"msg.cache.count":               "gauge",
		"rrset.cache.count":             "gauge",
		"infra.cache.count":             "gauge",
		"key.cache.count":               "gauge",
		"num.query.authzone.up":         "derive",
		"num.query.authzone.down":       "derive",
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

		fields := strings.SplitN(line, "=", 2)
		if len(fields) != 2 {
			continue
		}

		metric := fields[0] //eg: thread0.num.queries

		pattern := find(metrics, metric)
		if pattern == "" {
			// unknown metric
			continue
		}
		if _, ok := metrics[pattern]; !ok {
			continue
		}

		// backwards compatibility with https://github.com/falzm/collectd-unbound by removing 'total' prefix
		metric = strings.TrimPrefix(metric, "total.")
		metric = strings.Replace(metric, ".", "_", -1)

		value, err := strconv.ParseFloat(fields[1], 64)
		if err != nil {
			log.Printf("error: unable to parse metric value: %v", err)
			continue
		}

		switch metrics[pattern] {
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

func find(metrics map[string]string, metric string) string {
	for pattern := range metrics {
		if match, _ := filepath.Match(pattern, metric); match {
			return pattern
		}
	}
	return ""
}
