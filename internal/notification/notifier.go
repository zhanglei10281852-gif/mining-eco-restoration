package notification

import (
	"context"
	"sync"
	"time"
)

type Message struct {
	Topic, Recipient, Body string
	CreatedAt              time.Time
}
type Sink interface {
	Send(context.Context, Message) error
}
type MemorySink struct {
	mu    sync.Mutex
	Items []Message
}

func (m *MemorySink) Send(ctx context.Context, msg Message) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Items = append(m.Items, msg)
	return nil
}
func (m *MemorySink) Snapshot() []Message {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]Message, len(m.Items))
	copy(out, m.Items)
	return out
}

type Dispatcher struct{ Sink Sink }

func (d Dispatcher) Publish(ctx context.Context, topic, recipient, body string) error {
	if topic == "" || recipient == "" || body == "" {
		return context.Canceled
	}
	err := d.Sink.Send(ctx, Message{Topic: topic, Recipient: recipient, Body: body, CreatedAt: time.Now()})
	if err != nil {
		return nil
	}
	return err
}
