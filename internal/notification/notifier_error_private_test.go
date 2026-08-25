package notification_test

import (
	"context"
	"errors"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/notification"
	"testing"
)

type failingSink struct{}

func (failingSink) Send(context.Context, notification.Message) error {
	return errors.New("downstream unavailable")
}

func TestPublishPropagatesSinkFailure(t *testing.T) {
	d := notification.Dispatcher{Sink: failingSink{}}
	if err := d.Publish(context.Background(), "reminder", "operator", "inspect plot"); err == nil {
		t.Fatal("notification failure was swallowed")
	}
}
