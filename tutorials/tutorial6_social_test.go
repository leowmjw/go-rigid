//go:build tdd

package tutorials

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/leowmjw/go-rigid/internal/server"
)

func TestT6_Social_FollowUnfollowAdjacency(t *testing.T) {
	mux := http.NewServeMux()
	server.RegisterHandlers(mux, slog.Default())

	post := func(ev map[string]any) {
		b,_ := json.Marshal(map[string]any{"data": ev, "ackLevel": "ack"})
		req := httptest.NewRequest(http.MethodPost, "/rest/social/depot/%2Aevents/append", bytes.NewReader(b))
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
	}
	post(map[string]any{"type":"Follow","src":"u2","dst":"u1"})
	post(map[string]any{"type":"Unfollow","src":"u2","dst":"u1"})

	b,_ := json.Marshal([]any{"u2"})
	req := httptest.NewRequest(http.MethodPost, "/rest/social/pstate/%24%24adj/selectOne", bytes.NewReader(b))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	_ = rec
}
