package ws

import (
	"fmt"
	"net/http"

	"github.com/olahol/melody"
	"github.com/rs/zerolog/log"
)

var m = melody.New()

func InitWS() {

	m.Upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	m.HandleConnect(func(s *melody.Session) {
		log.Info().Str("remote_addr", s.Request.RemoteAddr).Msg("New User Connected")
		m.Broadcast([]byte("New user have connected"))
	})

	m.HandleDisconnect(func(s *melody.Session) {
		log.Info().Str("remote_addr", s.Request.RemoteAddr).Msg("User Disconnected")
		m.Broadcast([]byte("1 have left the room"))
	})

	m.HandleMessage(func(s *melody.Session, b []byte) {
		if _, ok := s.Get("username"); !ok {
			s.Set("username", string(b))
			s.Write([]byte("Welcome " + string(b)))
			return
		}
		if name, ok := s.Get("username"); ok {
			broadcastMsg := fmt.Sprintf("%s: %s", name, string(b))
			m.Broadcast([]byte(broadcastMsg))
		}
	})

}
func HandleWs(w http.ResponseWriter, r *http.Request) {
	if err := m.HandleRequest(w, r); err != nil {
		log.Error().Err(err).Msg("WebSocket upgrade failed")
		http.Error(w, "WebSocket upgrade failed", http.StatusInternalServerError)
	}
}
