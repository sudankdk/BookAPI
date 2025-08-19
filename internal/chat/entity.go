package chat

import "time"

type Message struct {
	ID     string    `json:"id"` // stream ID
	From   string    `json:"from"`
	To     string    `json:"to"`
	Body   string    `json:"body"`
	Type   string    `json:"type"` // chat|delivered|read
	SentAt time.Time `json:"sent_at"`
}

type SendRequest struct {
	To   string `json:"to"`
	Body string `json:"body"`
	CID  string `json:"cid"`  // idempotency key (optional)
	Type string `json:"type"` // default "chat"
}
