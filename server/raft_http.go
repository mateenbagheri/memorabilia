package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/mateenbagheri/memorabilia/pkg/cluster"
	"github.com/mateenbagheri/memorabilia/pkg/replication"
)

// RaftHTTPHandler exposes Raft cluster management operations over plain HTTP:
//
//	POST /raft/join    add a new voting peer (leader only)
//	GET  /raft/leader  return the current leader's Raft address
//	GET  /raft/peers   return the full cluster configuration as JSON
//
// This is the HTTP-transport equivalent of CommandServer: CommandServer
// exposes data operations (Get/Set/Delete) over gRPC, RaftHTTPHandler
// exposes cluster operations (join/leader/peers) over HTTP.
//
// Its only dependency is *replication.Node — it knows nothing about gRPC,
// the scheduler, or the FSM's underlying repository. That makes it possible
// to unit test these three routes with httptest against a node fixture,
// without booting the rest of the server.
type RaftHTTPHandler struct {
	node   *replication.Node
	logger *slog.Logger
}

// NewRaftHTTPHandler constructs a handler for the given Raft node.
func NewRaftHTTPHandler(node *replication.Node, logger *slog.Logger) *RaftHTTPHandler {
	return &RaftHTTPHandler{node: node, logger: logger}
}

// RegisterRoutes registers all Raft management routes on mux.
// Call this once during server startup before starting the HTTP server.
func (h *RaftHTTPHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/raft/join", h.handleJoin)
	mux.HandleFunc("/raft/leader", h.handleLeader)
	mux.HandleFunc("/raft/peers", h.handlePeers)
}

// handleJoin accepts a JSON-encoded replication.JoinRequest body and adds the
// caller as a voting peer via raft.AddVoter.
//
// Only the leader can add voters. If this node is not the leader, it responds
// 421 Misdirected Request with the current leader's Raft address in the body,
// so the caller can retry against the correct node.
func (h *RaftHTTPHandler) handleJoin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req replication.JoinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	if !h.node.IsLeader() {
		leader := h.node.LeaderRaftAddr()
		if leader == "" {
			http.Error(w, "no leader elected yet", http.StatusServiceUnavailable)
			return
		}
		http.Error(w,
			fmt.Sprintf("not the leader; forward to leader raft addr %q", leader),
			http.StatusMisdirectedRequest,
		)
		return
	}

	if err := h.node.Join(req.NodeID, req.RaftAddr); err != nil {
		h.logger.Error("join failed", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Info("node joined cluster",
		slog.String("nodeID", req.NodeID),
		slog.String("raftAddr", req.RaftAddr),
	)
	w.WriteHeader(http.StatusOK)
}

// handleLeader returns the current leader's Raft transport address as plain text.
// Responds 503 if no leader has been elected yet.
func (h *RaftHTTPHandler) handleLeader(w http.ResponseWriter, r *http.Request) {
	leader := h.node.LeaderRaftAddr()
	if leader == "" {
		http.Error(w, "no leader elected", http.StatusServiceUnavailable)
		return
	}
	fmt.Fprintln(w, leader)
}

// handlePeers returns the full cluster configuration (every known server and
// its Raft address) as JSON.
func (h *RaftHTTPHandler) handlePeers(w http.ResponseWriter, r *http.Request) {
	m := cluster.NewMembership(h.node)
	servers, err := m.Servers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(servers)
}
