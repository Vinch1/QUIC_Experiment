package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/leo/quic-raft/internal/raft"
)

type proposalRequest struct {
	Command string `json:"command"`
	Key     string `json:"key"`
	Value   string `json:"value"`
}

type proposalResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message,omitempty"`
}

type getResponse struct {
	OK    bool   `json:"ok"`
	Key   string `json:"key"`
	Value string `json:"value,omitempty"`
	Found bool   `json:"found"`
}

func NewHandler(node *raft.Node) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeJSON(w, http.StatusMethodNotAllowed, proposalResponse{
				OK:      false,
				Message: "method not allowed",
			})
			return
		}

		writeJSON(w, http.StatusOK, node.Status())
	})

	mux.HandleFunc("/kv", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGet(w, r, node)
		case http.MethodPost:
			handlePost(w, r, node)
		default:
			writeJSON(w, http.StatusMethodNotAllowed, proposalResponse{
				OK:      false,
				Message: "method not allowed",
			})
		}
	})

	return mux
}

func handleGet(w http.ResponseWriter, r *http.Request, node *raft.Node) {
	key := r.URL.Query().Get("key")
	if key == "" {
		writeJSON(w, http.StatusBadRequest, proposalResponse{
			OK:      false,
			Message: "key is required",
		})
		return
	}

	value, found, err := node.Get(key)
	if err != nil {
		writeJSON(w, http.StatusServiceUnavailable, proposalResponse{
			OK:      false,
			Message: err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, getResponse{
		OK:    true,
		Key:   key,
		Value: value,
		Found: found,
	})
}

func handlePost(w http.ResponseWriter, r *http.Request, node *raft.Node) {
	var req proposalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, proposalResponse{
			OK:      false,
			Message: "invalid json body",
		})
		return
	}

	if req.Command == "" {
		req.Command = "put"
	}

	if req.Command != "put" {
		writeJSON(w, http.StatusBadRequest, proposalResponse{
			OK:      false,
			Message: "unsupported command",
		})
		return
	}

	if req.Key == "" {
		writeJSON(w, http.StatusBadRequest, proposalResponse{
			OK:      false,
			Message: "key is required",
		})
		return
	}

	if err := node.Propose(req.Key, req.Value); err != nil {
		writeJSON(w, http.StatusServiceUnavailable, proposalResponse{
			OK:      false,
			Message: err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, proposalResponse{
		OK:      true,
		Message: "proposal accepted",
	})
}

func writeJSON(w http.ResponseWriter, code int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(value)
}
