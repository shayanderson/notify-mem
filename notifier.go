package notifymem

import "fmt"

type Notifier interface {
	Notify(title, message string) error
}

type notifier struct {
}

func NewNotifier() *notifier {
	return &notifier{}
}

func (n *notifier) Notify(title, message string) error {
	fmt.Println("notify:", title, message)
	return nil
}
