//go:build tdd

package server_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"log/slog"

	"github.com/leowmjw/go-rigid/internal/server"
)

func TestServer_DepotAppend_AckLevels(t *testing.T) {
	mux := http.NewServeMux()
	server.RegisterHandlers(mux, slog.New(slog.NewTextHandler(&bytes.Buffer{}, nil)))
	cases := []struct{ ack string; expectNonEmpty bool }{
		{"ack", true}, {"appendAck", false}, {"none", false}, {"bad", false},
	}
	for _, tc := range cases {
		body := []byte(`{"data":123,"ackLevel":"` + tc.ack + `"}`)
		r := httptest.NewRequest(http.MethodPost, "/rest/mod/depot/foo", bytes.NewReader(body))
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)
		resp := w.Result()
		if tc.ack == "bad" {
			if resp.StatusCode != http.StatusBadRequest { t.Fatalf("expected 400 for bad ack, got %d", resp.StatusCode) }
			continue
		}
		if resp.StatusCode != http.StatusOK { t.Fatalf("expected 200 got %d (ack=%s)", resp.StatusCode, tc.ack) }
		buf := new(bytes.Buffer); buf.ReadFrom(resp.Body)
		if tc.expectNonEmpty && !bytes.Contains(buf.Bytes(), []byte("topology")) {
			t.Fatalf("expected topology field for ack level ack")
		}
	}
}
