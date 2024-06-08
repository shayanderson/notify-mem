package notifymem

import (
	"testing"
	"time"
)

func TestNewMonitor(t *testing.T) {
	// test NewMonitor
	n := NewNotifier()
	mOpts := MonitorOptions{
		Interval:    2 * time.Second,
		ResendDelay: 30 * time.Second,
		Threshold:   80,
	}
	_, err := NewMonitor(n, mOpts)
	if err != nil {
		t.Fatalf("expected nil, got: %s", err)
	}

	mOpts.Interval = 0 * time.Millisecond
	_, err = NewMonitor(n, mOpts)
	if err == nil {
		t.Fatalf("expected bad interval error, got nil")
	}

	mOpts.Interval = 2 * time.Second
	mOpts.ResendDelay = 0 * time.Second
	_, err = NewMonitor(n, mOpts)
	if err == nil {
		t.Fatalf("expected bad resend delay error, got nil")
	}

	mOpts.ResendDelay = 30 * time.Second
	mOpts.Threshold = -1
	_, err = NewMonitor(n, mOpts)
	if err == nil {
		t.Fatalf("expected bad threshold error, got nil")
	}
}

type mockNotifier struct {
	notifyCalled bool
}

func (m *mockNotifier) Notify(title, message string) error {
	m.notifyCalled = true
	return nil
}

func TestMonitorNotify(t *testing.T) {
	// test Monitor.Notify
	m := &Monitor{}
	m.notifier = &mockNotifier{}
	m.notify(80)
	if !m.notifier.(*mockNotifier).notifyCalled {
		t.Fatalf("expected notify to be called")
	}
}
func TestMemoryUsage(t *testing.T) {
	mem := memory{
		total: 100,
		avail: 50,
	}

	usage, err := memoryUsage(mem)
	if err != nil {
		t.Fatalf("expected nil, got: %s", err)
	}

	if usage != 50 {
		t.Fatalf("expected 50, got: %d", usage)
	}
}
