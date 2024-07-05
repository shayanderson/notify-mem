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
	debugFunc         func(string)
	interval          time.Duration
	lastSentAt        time.Time
	notificationsSent int
	notifier          Notifier
	resendDelay       time.Duration
	threshold         int
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

	if m.interval < 1*time.Second {
		return nil, errors.New("interval must be at least 1s")
	}

	if m.notifier == nil {
		return nil, errors.New("notifier is required")
	}

	if m.resendDelay < 5*time.Second {
		return nil, errors.New("resend delay must be at least 5s")
	}

	if m.threshold < 1 || m.threshold > 100 {
		return nil, errors.New("threshold must be between 1 and 100")
	}

	return m, nil
}

// debug prints a debug message when debug mode is enabled
func (m *Monitor) debug(msg string) {
	if m.debugFunc != nil {
		m.debugFunc(msg)
	}
}

// isThresholdReached checks if the memory usage threshold has been reached
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

// notify sends a notification if the resend delay has been reached
func (m *Monitor) notify(usage int) error {
	if time.Since(m.lastSentAt) < m.resendDelay {
		m.debug("resend delay not reached")
		return nil
	}

	m.debug("sending notification")
	fmt.Printf("threshold reached: %d%%\n", usage)
	nSentStr := ""
	if m.notificationsSent > 0 {
		nSentStr = fmt.Sprintf(" (%d)", m.notificationsSent)
	}
	m.lastSentAt = time.Now()
	err := m.notifier.Notify(
		"Memory Usage Threshold Reached [notifymem]"+nSentStr,
		fmt.Sprintf("Memory usage at %d%%", usage),
	)
	m.notificationsSent++
	if err != nil {
		return err
	}
	return nil
}

// Run starts the monitor
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
				if err := m.notify(usage); err != nil {
					return err
				}
			}

		}
	}
}

// Test executes a test to check memory usage and send a notification
func (m *Monitor) Test() error {
	m.debug("running test")
	_, usage, err := m.isThresholdReached()
	if err != nil {
		return err
	}
	m.notify(usage)
	return nil
}

// memoryUsage calculates memory usage as a percentage
func memoryUsage(m memory) (int, error) {
	if m.total == 0 {
		return 0, errors.New("total memory is 0")
	}
	if m.avail == 0 {
		return 0, errors.New("available memory is 0")
	}

	return int((float64(m.total-m.avail) / float64(m.total)) * 100), nil
}

// readMemory reads memory information from /proc/meminfo
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
