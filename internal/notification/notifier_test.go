package notification_test

import (
	"context"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/notification"
	"testing"
)

func TestPublishAndSnapshotIsolation(t *testing.T) {
	sink := &notification.MemorySink{}
	d := notification.Dispatcher{Sink: sink}
	if e := d.Publish(context.Background(), "reminder", "user", "body"); e != nil {
		t.Fatal(e)
	}
	items := sink.Snapshot()
	if len(items) != 1 || items[0].Topic != "reminder" {
		t.Fatal(items)
	}
	items[0].Body = "changed"
	if sink.Snapshot()[0].Body != "body" {
		t.Fatal("snapshot aliases storage")
	}
}
func TestPublishValidationAndCancel(t *testing.T) {
	d := notification.Dispatcher{Sink: &notification.MemorySink{}}
	if e := d.Publish(context.Background(), "", "u", "b"); e == nil {
		t.Fatal("empty topic")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if e := d.Publish(ctx, "t", "u", "b"); e == nil {
		t.Fatal("cancel ignored")
	}
}
