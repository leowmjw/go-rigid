package server

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

var DefaultShutdownTimeout = 10 * time.Second

type Server struct {
	Logger *slog.Logger
	// Tutorial-specific in-memory stores for TDD
	tutorialWordCounts map[string]int
	tutorialCounts     map[string]map[string]int // user -> type -> count
	tutorialAdj        map[string][]string       // src -> []dst
}

func RegisterHandlers(mux *http.ServeMux, logger *slog.Logger) {
	s := &Server{
		Logger:             logger,
		tutorialWordCounts: make(map[string]int),
		tutorialCounts:     make(map[string]map[string]int),
		tutorialAdj:        make(map[string][]string),
	}
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

	// Parse module and depot
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 6 || parts[1] != "rest" || parts[3] != "depot" {
		writeErr(w, http.StatusBadRequest, "invalid path", nil)
		return
	}
	module := parts[2]
	depot := parts[4]

	// Module-specific processing
	if module == "tutorial" || module == "social" {
		s.processTutorialDepot(depot, req.Data)
	}

	w.Header().Set("Content-Type", "application/json")
	if req.AckLevel == "ack" {
		w.Write([]byte(`{"topology":"success"}`)) // placeholder non-empty map
		return
	}
	w.Write([]byte(`{}`))
}

func (s *Server) processTutorialDepot(depot string, data any) {
	switch depot {
	case "*text": // *text
		if str, ok := data.(string); ok {
			words := strings.Fields(str)
			for _, word := range words {
				s.tutorialWordCounts[word]++
			}
		}
	case "*events": // *events
		if m, ok := data.(map[string]any); ok {
			user, _ := m["user"].(string)
			typ, _ := m["type"].(string)
			if user != "" && typ != "" {
				if s.tutorialCounts[user] == nil {
					s.tutorialCounts[user] = make(map[string]int)
				}
				s.tutorialCounts[user][typ]++
			}
			// For social
			if typ == "Follow" || typ == "Unfollow" {
				src, _ := m["src"].(string)
				dst, _ := m["dst"].(string)
				if typ == "Follow" {
					s.tutorialAdj[src] = append(s.tutorialAdj[src], dst)
				} else if typ == "Unfollow" {
					// Remove dst from src's list
					list := s.tutorialAdj[src]
					for i, d := range list {
						if d == dst {
							s.tutorialAdj[src] = append(list[:i], list[i+1:]...)
							break
						}
					}
				}
			}
		}
	}
}

func (s *Server) handlePStateSelect(w http.ResponseWriter, r *http.Request) {
	// Parse path: /rest/{module}/pstate/{pstate}/{op}
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 6 || parts[1] != "rest" || parts[3] != "pstate" {
		writeErr(w, http.StatusBadRequest, "invalid path", nil)
		return
	}
	module := parts[2]
	pstate := parts[4]
	// op := parts[5] // select or selectOne, assume selectOne

	var path []any
	if err := json.NewDecoder(r.Body).Decode(&path); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid path JSON", err)
		return
	}

	if module == "tutorial" || module == "social" {
		result := s.selectTutorialPState(pstate, path)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
		return
	}

	writeErr(w, http.StatusNotImplemented, "pstate select not implemented for module", nil)
}

func (s *Server) selectTutorialPState(pstate string, path []any) any {
	switch pstate {
	case "$$wordCounts": // $$wordCounts
		if len(path) == 1 {
			if word, ok := path[0].(string); ok {
				return s.tutorialWordCounts[word]
			}
		}
	case "$$counts": // $$counts
		if len(path) == 2 {
			user, ok1 := path[0].(string)
			typ, ok2 := path[1].(string)
			if ok1 && ok2 {
				if userCounts, exists := s.tutorialCounts[user]; exists {
					return userCounts[typ]
				}
			}
		}
	case "$$adj": // $$adj
		if len(path) == 1 {
			if src, ok := path[0].(string); ok {
				return s.tutorialAdj[src]
			}
		}
	}
	return nil
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
