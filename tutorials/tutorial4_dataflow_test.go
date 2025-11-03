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

func TestT4_Dataflow_BranchingAndTransform(t *testing.T) {
	mux := http.NewServeMux()
	server.RegisterHandlers(mux, slog.Default())

	b,_ := json.Marshal(map[string]any{"data": map[string]any{"user":"u1","type":"click"}, "ackLevel":"ack"})
	req := httptest.NewRequest(http.MethodPost, "/rest/tutorial/depot/%2Aevents/append", bytes.NewReader(b))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	pathJSON := []any{"u1","click"}
	pb,_ := json.Marshal(pathJSON)
	req2 := httptest.NewRequest(http.MethodPost, "/rest/tutorial/pstate/%24%24counts/selectOne", bytes.NewReader(pb))
	rec2 := httptest.NewRecorder()
	mux.ServeHTTP(rec2, req2)
	_ = rec2
}
