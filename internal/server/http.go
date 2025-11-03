package server

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"
)

var DefaultShutdownTimeout = 10 * time.Second

type Server struct {
	Logger *slog.Logger
}

func RegisterHandlers(mux *http.ServeMux, logger *slog.Logger) {
	s := &Server{Logger: logger}
	mux.HandleFunc("POST /rest/{module}/depot/", s.handleDepotAppend)   // /append
	mux.HandleFunc("POST /rest/{module}/pstate/", s.handlePStateSelect) // /select|selectOne
	mux.HandleFunc("POST /rest/{module}/query/", s.handleQueryInvoke)   // /invoke
}

type DepotAppendRequest struct {
	Data     any    `json:"data"`
	AckLevel string `json:"ackLevel,omitempty"`
}

func (s *Server) handleDepotAppend(w http.ResponseWriter, r *http.Request) {
	var req DepotAppendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid JSON body", err)
		return
	}
	if req.AckLevel == "" { req.AckLevel = "ack" }
	switch req.AckLevel {
	case "ack", "appendAck", "none":
	default:
		writeErr(w, http.StatusBadRequest, "invalid ackLevel", errors.New("ack|appendAck|none"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if req.AckLevel == "ack" {
		w.Write([]byte(`{"topology":"success"}`)) // placeholder non-empty map
		return
	}
	w.Write([]byte(`{}`))
}

func (s *Server) handlePStateSelect(w http.ResponseWriter, r *http.Request) {
	writeErr(w, http.StatusNotImplemented, "pstate select not implemented", nil)
}

func (s *Server) handleQueryInvoke(w http.ResponseWriter, r *http.Request) {
	writeErr(w, http.StatusNotImplemented, "query invoke not implemented", nil)
}

func writeErr(w http.ResponseWriter, code int, msg string, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"error": msg,
		"detail": func() any { if err!=nil { return err.Error() }; return nil }(),
	})
}
