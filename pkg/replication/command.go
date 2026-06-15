package replication

import (
	"encoding/json"
	"time"
)

type OpType uint8

const (
	OpSet OpType = iota
	OpDelete
	OpBatchDelete
)

type RaftCommand struct {
	Op         OpType    `json:"op"`
	Expiration time.Time `json:"expiration,omitempty"` // default == no expiration
	Value      string    `json:"value,omitempty"`
	Key        string    `json:"key,omitempty"`
	Keys       []string  `json:"keys,omitempty"`
}

// Encode serializes a raft command mainly for raft.Apply()
func (rc *RaftCommand) Encode() ([]byte, error) {
	return json.Marshal(rc)
}

func DecodeCommand(command []byte) (*RaftCommand, error) {
	var decoded RaftCommand
	if err := json.Unmarshal(command, &decoded); err != nil {
		return nil, err
	}
	return &decoded, nil
}
