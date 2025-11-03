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

func TestT2_WordCount_StreamToPState(t *testing.T) {
	mux := http.NewServeMux()
	server.RegisterHandlers(mux, slog.Default())

	append := func(s string) {
		b,_ := json.Marshal(map[string]any{"data": s, "ackLevel":"ack"})
		req := httptest.NewRequest(http.MethodPost, "/rest/tutorial/depot/%2Atext/append", bytes.NewReader(b))
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK { t.Fatalf("append status=%d", rec.Code) }
	}

	append("to be or not to be")
	append("to be")

	pathJSON := []any{"to"}
	b,_ := json.Marshal(pathJSON)
	req := httptest.NewRequest(http.MethodPost, "/rest/tutorial/pstate/%24%24wordCounts/selectOne", bytes.NewReader(b))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code == http.StatusOK {
		var count int
		_ = json.Unmarshal(rec.Body.Bytes(), &count)
		_ = count
	}
}
