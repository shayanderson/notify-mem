package notifymem

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"runtime"
	"time"
)

type memory struct {
	total int
	avail int
}

type MonitorOptions struct {
	DebugFunc   func(string)
	Interval    time.Duration
	ResendDelay time.Duration
	Threshold   int
}

type Monitor struct {
	debugFunc   func(string)
	interval    time.Duration
	lastSentAt  time.Time
	notifier    Notifier
	resendDelay time.Duration
	threshold   int
}

func NewMonitor(notifier Notifier, opts MonitorOptions) (*Monitor, error) {
	if runtime.GOOS != "linux" {
		return nil, errors.New("notifymem only supports Linux")
	}

	m := &Monitor{
		debugFunc:   opts.DebugFunc,
		interval:    opts.Interval,
		notifier:    notifier,
		resendDelay: opts.ResendDelay,
		threshold:   opts.Threshold,
	}

	if m.interval < 100*time.Millisecond {
		return nil, errors.New("interval must be at least 100ms")
	}

	if m.notifier == nil {
		return nil, errors.New("notifier is required")
	}

	if m.resendDelay < 5*time.Second {
		return nil, errors.New("resend delay must be at least 5s")
	}

	if m.threshold < 0 || m.threshold > 100 {
		return nil, errors.New("threshold must be between 0 and 100")
	}

	return m, nil
}

func (m *Monitor) debug(msg string) {
	if m.debugFunc != nil {
		m.debugFunc(msg)
	}
}

func (m *Monitor) isThresholdReached() (bool, int, error) {
	mem, err := readMemory()
	if err != nil {
		return false, 0, err
	}

	usage, err := memoryUsage(mem)
	if err != nil {
		return false, 0, err
	}

	if usage >= m.threshold {
		return true, usage, nil
	}

	return false, usage, nil
}

func (m *Monitor) notify(usage int) error {
	if time.Since(m.lastSentAt) < m.resendDelay {
		m.debug("resend delay not reached")
		return nil
	}

	m.debug("sending notification")
	m.lastSentAt = time.Now()
	return m.notifier.Notify(
		"Memory Usage Threshold Reached [notifymem]",
		fmt.Sprintf("Memory usage at %d%%", usage),
	)
}

func (m *Monitor) Run(ctx context.Context) error {
	m.debug("staring monitor")
	for {
		select {
		case <-ctx.Done():
			m.debug("stopping monitor")
			return ctx.Err()

		case <-time.After(m.interval):
			reached, usage, err := m.isThresholdReached()
			if err != nil {
				return err
			}
			m.debug(fmt.Sprintf("memory usage: %d%%", usage))

			if reached {
				m.debug(fmt.Sprintf("threshold reached: %d%%", usage))
				if err := m.notify(usage); err != nil {
					return err
				}
			}

		}
	}
}

func memoryUsage(m memory) (int, error) {
	if m.total == 0 {
		return 0, errors.New("total memory is 0")
	}

	return int((float64(m.total-m.avail) / float64(m.total)) * 100), nil
}

func readMemory() (memory, error) {
	m := memory{}
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		return m, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	for s.Scan() {
		var k string
		var v int
		_, err := fmt.Sscanf(s.Text(), "%s %d", &k, &v)
		if err != nil {
			continue
		}
		switch k {
		case "MemTotal:":
			m.total = v
		case "MemAvailable:":
			m.avail = v
		}
	}

	return m, nil
}
