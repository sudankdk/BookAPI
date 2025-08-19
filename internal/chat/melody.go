package chat

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/olahol/melody"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/sudankdk/bookstore/internal/chat"
	"github.com/sudankdk/bookstore/pkg/utils"
)

type Config struct {
	JWTSecret string
	service   Service
	logger    zerolog.Logger
	pool      *redis.Client
}

type MelodyWs struct {
	cfg Config
	mdy *melody.Melody
}

func NewMelody(cfg Config) *MelodyWs {
	m := melody.New()
	m.Config.MaxMessageSize = 1 << 20 //1 mb
	m.Upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	ws := &MelodyWs{
		cfg: cfg,
		mdy: m,
	}
	return ws
}

func (s *MelodyWs) Handle(w http.ResponseWriter, r *http.Request) {
	uid, err := utils.ExtractUserIDFromAuthHeader(s.cfg.JWTSecret, r.Header.Get("Authorization"), r.URL.Query().Get("auth"))
	if err != nil || uid == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	ctx := context.WithValue(r.Context(), "userID", uid)
	s.mdy.HandleRequest(w, r.WithContext(ctx))
}

func (ws *MelodyWs) bindHandlers() {
	ws.mdy.HandleConnect(func(s *melody.Session) {
		uid, _ := s.Request.Context().Value("userID").(string)
		if uid == "" {
			_ = s.CloseWithMsg([]byte("unauthorized"))
			return
		}

		//sub to realtime pubsub
		cancel, err := ws.cfg.service.Subscribe(context.Background(), uid, func(b []byte) {
			_ = s.Write(b)
		})
		if err != nil {
			ws.cfg.logger.Err(err).Msg("subscribe failed")
			_ = s.Close()
			return
		}
		s.Set("cancel", cancel)

		_ = ws.cfg.service.Replay(context.Background(), uid, func(b []byte) error {
			return s.Write(b)
		})
	})
	ws.mdy.HandleDisconnect(func(s *melody.Session) {
		if v, ok := s.Get("cancel"); ok {
			if cancel, ok := v.(func()); ok {
				cancel()
			}
		}
	})

	ws.mdy.HandleMessage(func(s *melody.Session, b []byte) {
		uid, _ := s.Request.Context().Value("userID").(string)
		if uid == "" {
			_ = s.CloseWithMsg([]byte("unauthorized"))
			return
		}
		// parse incoming
		var in chat.SendRequest
		if err := json.Unmarshal(b, &in); err != nil {
			_ = s.Write([]byte(`{"error":"invalid_json"}`))
			return
		}
		in.Body = strings.TrimSpace(in.Body)
		if in.To == "" || in.Body == "" {
			_ = s.Write([]byte(`{"error":"invalid_payload"}`))
			return
		}
		// send
		id, err := ws.cfg.service.Send(s.Request.Context(), uid, in)
		if err != nil {
			_ = s.Write([]byte(`{"error":"send_failed"}`))
			return
		}
		// ack to sender
		_ = s.Write([]byte(`{"type":"ack","id":"` + id + `","at":"` + time.Now().UTC().Format(time.RFC3339Nano) + `"}`))
	})
}
