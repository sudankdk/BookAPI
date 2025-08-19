// package chat

// package ws

// import (
// 	"context"
// 	"encoding/json"
// 	"net/http"
// 	"strings"
// 	"time"

// 	"github.com/go-chi/chi/v5"
// 	"github.com/olahol/melody"
// 	"github.com/rs/zerolog"
// 	"github.com/gomodule/redigo/redis"

// 	"github.com/you/app/internal/auth"
// 	"github.com/you/app/internal/chat"
// )

// type Config struct {
// 	JWTSecret string
// 	Service   chat.ServicePort
// 	Logger    zerolog.Logger
// 	RedisPool *redis.Pool
// }

// type MelodyWS struct {
// 	cfg   Config
// 	mdy   *melody.Melody
// }

// func NewMelody(cfg Config) *MelodyWS {
// 	m := melody.New()
// 	m.Config.MaxMessageSize = 1 << 20 // 1MB
// 	m.Upgrader.CheckOrigin = func(r *http.Request) bool {
// 		// TODO: tighten in prod
// 		return true
// 	}

// 	ws := &MelodyWS{cfg: cfg, mdy: m}
// 	ws.bindHandlers()
// 	return ws
// }

// func (ws *MelodyWS) Router() chi.Router {
// 	r := chi.NewRouter()
// 	r.Get("/", ws.handle)
// 	return r
// }

// func (ws *MelodyWS) handle(w http.ResponseWriter, r *http.Request) {
// 	// Extract user from auth header or ?auth=
// 	uid, err := auth.ExtractUserIDFromAuthHeader(ws.cfg.JWTSecret, r.Header.Get("Authorization"), r.URL.Query().Get("auth"))
// 	if err != nil || uid == "" {
// 		http.Error(w, "unauthorized", http.StatusUnauthorized)
// 		return
// 	}
// 	ctx := context.WithValue(r.Context(), "userID", uid)
// 	_ = ws.mdy.HandleRequest(w, r.WithContext(ctx))
// }

// func (ws *MelodyWS) bindHandlers() {
// 	ws.mdy.HandleConnect(func(s *melody.Session) {
// 		uid, _ := s.Request.Context().Value("userID").(string)
// 		if uid == "" {
// 			_ = s.CloseWithMsg([]byte("unauthorized"))
// 			return
// 		}

// 		// Subscribe to realtime pubsub
// 		cancel, err := ws.cfg.Service.Subscribe(context.Background(), uid, func(b []byte) error {
// 			return ws.mdy.Write(s, b)
// 		})
// 		if err != nil {
// 			ws.cfg.Logger.Err(err).Msg("subscribe failed")
// 			_ = s.Close()
// 			return
// 		}
// 		s.Set("cancel", cancel)

// 		// Replay missed durable messages
// 		_ = ws.cfg.Service.Replay(context.Background(), uid, func(b []byte) error {
// 			return ws.mdy.Write(s, b)
// 		})
// 	})

// 	ws.mdy.HandleDisconnect(func(s *melody.Session) {
// 		if v, ok := s.Get("cancel"); ok {
// 			if cancel, ok := v.(func()); ok {
// 				cancel()
// 			}
// 		}
// 	})

// 	ws.mdy.HandleMessage(func(s *melody.Session, b []byte) {
// 		uid, _ := s.Request.Context().Value("userID").(string)
// 		if uid == "" {
// 			_ = s.CloseWithMsg([]byte("unauthorized"))
// 			return
// 		}
// 		// parse incoming
// 		var in chat.SendRequest
// 		if err := json.Unmarshal(b, &in); err != nil {
// 			_ = ws.mdy.Write(s, []byte(`{"error":"invalid_json"}`))
// 			return
// 		}
// 		in.Body = strings.TrimSpace(in.Body)
// 		if in.To == "" || in.Body == "" {
// 			_ = ws.mdy.Write(s, []byte(`{"error":"invalid_payload"}`))
// 			return
// 		}
// 		// send
// 		id, err := ws.cfg.Service.Send(s.Request.Context(), uid, in)
// 		if err != nil {
// 			_ = ws.mdy.Write(s, []byte(`{"error":"send_failed"}`))
// 			return
// 		}
// 		// ack to sender
// 		_ = ws.mdy.Write(s, []byte(`{"type":"ack","id":"`+id+`","at":"`+time.Now().UTC().Format(time.RFC3339Nano)+`"}`))
// 	})
// }
