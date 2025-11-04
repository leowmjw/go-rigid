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

func TestT7_QueryInvoke_SimpleQuery(t *testing.T) {
	mux := http.NewServeMux()
	server.RegisterHandlers(mux, slog.Default())

	params := map[string]any{"param1": "value1"}
	b, _ := json.Marshal(params)
	req := httptest.NewRequest(http.MethodPost, "/rest/tutorial/query/simpleQuery/invoke", bytes.NewReader(b))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("invoke status=%d body=%s", rec.Code, rec.Body.String())
	}
	var result map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &result); err != nil {
		t.Fatalf("failed to unmarshal result: %v", err)
	}
	// TODO: assert expected result
}
