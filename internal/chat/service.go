package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	streamPerUSer  = "dm:stream:"
	pubsubPrefix   = "dm:to:"
	ackKeyPrefix   = "dm:last_ack:"
	defaultMsgType = "chat"
)

type Service struct {
	Pool *redis.Client
}

func NewService(pool *redis.Client) *Service {
	return &Service{Pool: pool}
}

func (s *Service) Send(ctx context.Context, from string, req SendRequest) (string, error) {
	if req.Type == "" {
		req.Type = defaultMsgType
	}
	msg := Message{
		ID:     "",
		From:   from,
		To:     req.To,
		Body:   req.Body,
		Type:   req.Type,
		SentAt: time.Now().UTC(),
	}
	id, err := s.Pool.XAdd(ctx, &redis.XAddArgs{
		Stream: streamPerUSer,
		ID:     "*",
		Values: map[string]interface{}{
			"from":    msg.From,
			"to":      msg.To,
			"body":    msg.Body,
			"type":    msg.Type,
			"sent_at": msg.SentAt,
			"cid":     req.CID,
		},
	}).Result()
	if err != nil {
		return "", err
	}
	msg.ID = id
	msgByte, _ := json.Marshal(msg)
	if err := s.Pool.Publish(ctx, pubsubPrefix+msg.ID, msgByte); err != nil {
		println("Message publish to channel failed")
	}
	return id, nil

}

func (s *Service) Subscribe(ctx context.Context, userID string, write func([]byte)) (func(), error) {
	ch := pubsubPrefix + userID
	pubsub := s.Pool.Subscribe(ctx, ch)
	_, err := pubsub.Receive(ctx)
	if err != nil {
		return nil, err
	}
	stop := make(chan struct{})
	go func() {
		ch := pubsub.Channel()
		defer pubsub.Close()
		for {
			select {
			case msg, ok := <-ch:
				if !ok {
					return
				}
				write([]byte(msg.Payload))
			case <-ctx.Done():
				return

			case <-stop:
				return
			}
		}
	}()

	cancel := func() { close(stop) }
	return cancel, nil
}

func (s *Service) Replay(ctx context.Context, userID string, write func([]byte) error) error {
	lastkey := ackKeyPrefix + userID
	lastID, err := s.Pool.Get(ctx, lastkey).Result()
	if err == redis.Nil {
		lastID = "0-0"
	} else if err != nil {
		return err
	}
	next := lastID
	for {
		select {
		default:
			entries, err := s.Pool.XRange(ctx, streamPerUSer+userID, "("+next, "+").Result()
			if err != nil {
				return err
			}
			if len(entries) == 0 {
				return nil
			}
			for _, e := range entries {
				kv := e.Values
				to, _ := kv["to"].(string)
				if to == userID {
					wire := []byte(fmt.Sprintf(
						`{"id":"%s","from":"%s","to":"%s","body":"%s","type":"%s","sent_at":%s}`,
						e.ID,
						kv["from"],
						kv["to"],
						escapeJSON(fmt.Sprintf("%v", kv["body"])),
						kv["type"],
						kv["sent_at"],
					))
					_ = write(wire)
					next = e.ID
					_ = s.Pool.Set(ctx, lastkey, e.ID, 0).Err()
				}
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func escapeJSON(s string) string {
	b, _ := json.Marshal(s)
	if len(b) >= 2 {
		return string(b[1 : len(b)-1])
	}
	return s
}
