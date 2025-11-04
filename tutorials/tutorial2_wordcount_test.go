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
	if rec.Code != http.StatusOK {
		t.Fatalf("selectOne status=%d body=%s", rec.Code, rec.Body.String())
	}
	var count int
	if err := json.Unmarshal(rec.Body.Bytes(), &count); err != nil {
		t.Fatalf("failed to unmarshal count: %v", err)
	}
	if count != 3 {
		t.Fatalf("expected count 3 for 'to', got %d", count)
	}
}
