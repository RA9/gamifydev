package server

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/RA9/gamifydev/platform/internal/auth"
	"github.com/coder/websocket"
)

// hub fans out realtime notifications to a user's open connections. Messages are
// tiny topic strings (e.g. "submissions"); the browser reacts by asking htmx to
// refresh the relevant region, so the server stays the single source of truth.
type hub struct {
	mu      sync.RWMutex
	clients map[int64]map[*wsClient]struct{} // userID -> connections
}

func newHub() *hub { return &hub{clients: map[int64]map[*wsClient]struct{}{}} }

type wsClient struct {
	conn    *websocket.Conn
	writeMu sync.Mutex // serialize writes to this connection
}

func (h *hub) add(userID int64, c *wsClient) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.clients[userID] == nil {
		h.clients[userID] = map[*wsClient]struct{}{}
	}
	h.clients[userID][c] = struct{}{}
}

func (h *hub) remove(userID int64, c *wsClient) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if set := h.clients[userID]; set != nil {
		delete(set, c)
		if len(set) == 0 {
			delete(h.clients, userID)
		}
	}
}

// Notify sends a topic message to every live connection for a user. It is safe
// to call from any goroutine and never blocks the caller meaningfully.
func (h *hub) Notify(userID int64, topic string) {
	h.mu.RLock()
	clients := make([]*wsClient, 0, len(h.clients[userID]))
	for c := range h.clients[userID] {
		clients = append(clients, c)
	}
	h.mu.RUnlock()

	for _, c := range clients {
		go func(c *wsClient) {
			c.writeMu.Lock()
			defer c.writeMu.Unlock()
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = c.conn.Write(ctx, websocket.MessageText, []byte(topic))
		}(c)
	}
}

// handleWS upgrades to a WebSocket and keeps it open for server-pushed refresh
// signals. Auth is enforced by the requireAuth middleware; the socket is
// write-only from the server's side (reads are drained to detect disconnect).
func (s *Server) handleWS(w http.ResponseWriter, r *http.Request) {
	u := auth.CurrentUser(r.Context())
	if u == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	conn, err := websocket.Accept(w, r, nil) // default: same-origin only
	if err != nil {
		return
	}
	client := &wsClient{conn: conn}
	s.hub.add(u.ID, client)
	defer s.hub.remove(u.ID, client)

	// CloseRead drains incoming frames and returns a context cancelled when the
	// client disconnects; we hold the handler open until then.
	ctx := conn.CloseRead(r.Context())
	<-ctx.Done()
	conn.Close(websocket.StatusNormalClosure, "")
}
