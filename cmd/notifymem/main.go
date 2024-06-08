package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	notifymem "github.com/shayanderson/notify-mem"
)

var version string // do not remove

func errorExit(err error) {
	fmt.Fprintf(os.Stderr, "[error] %s\n", err)
	os.Exit(1)
}

func main() {
	fmt.Printf("running notifymem (version %s)\n", version)

	type flags struct {
		debug       bool
		interval    int
		resendDelay int
		threshold   int
	}

	f := flags{}
	flag.BoolVar(&f.debug, "debug", false, "enable debug mode")
	flag.IntVar(&f.interval, "interval", 2, "interval between memory checks in seconds")
	flag.IntVar(&f.resendDelay, "delay", 30, "delay between notifications being sent in seconds")
	flag.IntVar(&f.threshold, "threshold", 80, "memory threshold as a percentage")
	flag.Parse()

	n := notifymem.NewNotifier()
	mOpts := notifymem.MonitorOptions{
		Interval:    time.Duration(f.interval) * time.Second,
		ResendDelay: time.Duration(f.resendDelay) * time.Second,
		Threshold:   f.threshold,
	}

	if f.debug {
		mOpts.DebugFunc = func(s string) { fmt.Println(s) }
	}

	m, err := notifymem.NewMonitor(n, mOpts)
	if err != nil {
		errorExit(err)
	}

	ctx := context.Background()
	if err := m.Run(ctx); err != nil && err != context.Canceled {
		errorExit(err)
	}
}
