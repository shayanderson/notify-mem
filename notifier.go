package notifymem

import (
	"os/exec"
)

type Notifier interface {
	Notify(title, message string) error
}

type notifier struct {
}

func NewNotifier() *notifier {
	return &notifier{}
}

// Notify sends a notification using notify-send
func (n *notifier) Notify(title, message string) error {
	cmd := exec.Command("notify-send", title, message)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
