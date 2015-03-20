package main

import (
	"bufio"
	"bytes"
	"fmt"
	osExec "os/exec"
	"strconv"
	"time"

	"collectd.org/api"
	"collectd.org/exec"
)

func main() {
	e := exec.NewExecutor()
	e.VoidCallback(unboundStats, exec.Interval())
	e.Run()
}

func unboundStats(interval time.Duration) {
	var (
		err          error
		metric       []byte
		value        float64
		pos, advance int
		cmdStdOut    bytes.Buffer
	)

	cmd := osExec.Command("/bin/sh", "-c", "unbound-control stats")
	cmd.Stdout = &cmdStdOut

	if err := cmd.Run(); err != nil {
		fmt.Printf("error: unable to execute unbound-control: %s\n", err)
		return
	}

	now := time.Now()

	line := []byte{}
	for pos < cmdStdOut.Len() {
		if advance, line, err = bufio.ScanLines(cmdStdOut.Bytes()[pos:], true); err != nil {
			fmt.Printf("error: unable to parse line: %s\n", err)
			return
		}

		if bytes.HasPrefix(line, []byte("total.")) {
			splitLine := bytes.Split(line, []byte("="))
			metric = bytes.Replace(bytes.TrimPrefix(splitLine[0], []byte("total.")), []byte("."), []byte("_"), -1)

			if value, err = strconv.ParseFloat(string(splitLine[1]), 32); err != nil {
				fmt.Printf("error: unable to parse metric value: %s\n", err)
				continue
			}

			vl := api.ValueList{
				Identifier: api.Identifier{
					Host:         exec.Hostname(),
					Plugin:       "unbound",
					Type:         "gauge",
					TypeInstance: string(metric),
				},
				Time:     now,
				Interval: exec.Interval(),
				Values:   []api.Value{api.Gauge(value)},
			}

			exec.Putval.Dispatch(vl)
		}

		// Set scanner position to the beginning of the next line
		pos += advance
	}
}
