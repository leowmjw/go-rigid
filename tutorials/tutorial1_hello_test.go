//go:build tdd

package tutorials

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/rig/internal/server"
)

func TestT1_DepotAppend_WithAck_ReturnsStreamingAck(t *testing.T) {
	mux := http.NewServeMux()
	server.RegisterHandlers(mux, slog.Default())

	b, _ := json.Marshal(map[string]any{"data": "hello world", "ackLevel": "ack"})
	req := httptest.NewRequest(http.MethodPost, "/rest/tutorial/depot/%2Ainput/append", bytes.NewReader(b))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
	var got map[string]any
	_ = json.Unmarshal(rec.Body.Bytes(), &got)
	if len(got) == 0 {
		t.Fatalf("expected non-empty streaming ack map; got %v", got)
	}
}
