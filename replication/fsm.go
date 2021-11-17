package replication

import (
	"SDCC/database"
	"SDCC/metadata"
	"SDCC/utils"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/raft"
	"io"
	"os"
	"strings"
)

type CommandPayload struct {
	Operation string
	Key       []byte
	Value     [][]byte
}

type ApplyResponse struct {
	Error error
}

type AppendResponse struct {
	Error error
	ToBeOffloaded bool
}

var keyDb = metadata.GetKeyDb()

// snapshotNoop handle noop snapshot
type snapshotNoop struct{}

// Persist persists to disk. Return nil on success, otherwise return error.
func (s snapshotNoop) Persist(_ raft.SnapshotSink) error { return nil }

// Release releases the lock after persist snapshot.
// Release is invoked when we are finished with the snapshot.
func (s snapshotNoop) Release() {}

// newSnapshotNoop is returned by an FSM in response to a snapshotNoop
// It must be safe to invoke FSMSnapshot methods with concurrent
// calls to Apply.
func newSnapshotNoop() (raft.FSMSnapshot, error) {
	return &snapshotNoop{}, nil
}

type FSM struct {
	db      *database.Database
	cluster *utils.ClusterRoutine
	dht     *DhtRoutine
}

// Apply log is invoked once a log entry is committed.
// It returns a value which will be made available in the
// ApplyFuture returned by Raft.Apply method if that
// method was called on the same Raft node as the FSM.
func (f FSM) Apply(log *raft.Log) interface{} {
	switch log.Type {
	case raft.LogCommand:
		var payload = CommandPayload{}
		if err := json.Unmarshal(log.Data, &payload); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "error marshalling store payload %s\n", err.Error())
			return nil
		}

		op := strings.ToUpper(strings.TrimSpace(payload.Operation))
		switch op {
		case "PUT":
			return &ApplyResponse{
				Error: (*f.db).Put(payload.Key, payload.Value),
			}
		case "APPEND":
			val, err := (*f.db).Append(payload.Key, payload.Value)
			if err == nil {
				if utils.GetSize(val) > utils.Threshold {
					err = (*f.db).Del(payload.Key)
					return &AppendResponse{
						Error: err,
						ToBeOffloaded: true,
					}
				} else {
					return &AppendResponse{
						Error: nil,
						ToBeOffloaded: false,
					}
				}
			} else {
				return &AppendResponse{
					Error: err,
				}
			}
		case "DELETE":
			keyDb.DelKey(string(payload.Key))
			return &ApplyResponse{
				Error: (*f.db).Del(payload.Key),
			}
		case "JOIN":
			f.cluster.Join(string(payload.Key))
			f.dht.UpdateCluster()
			return &ApplyResponse{
				Error: nil,
			}
		case "LEAVE":
			f.cluster.Leave(string(payload.Key))
			f.dht.UpdateCluster()
			return &ApplyResponse{
				Error: nil,
			}
		}
	}

	_, _ = fmt.Fprintf(os.Stderr, "not replication log command type\n")
	return nil
}

// Snapshot will be called during make snapshot.
// Snapshot is used to support log compaction.
// No need to call snapshot since it already persisted in disk (using BadgerDB) when replication calling Apply function.
func (f FSM) Snapshot() (raft.FSMSnapshot, error) {
	return newSnapshotNoop()
}

// Restore is used to restore an FSM from a Snapshot. It is not called
// concurrently with any other command. The FSM must discard all previous
// state.
// Restore will update all data in BadgerDB
func (f FSM) Restore(rClose io.ReadCloser) error {
	defer func() {
		if err := rClose.Close(); err != nil {
			_, _ = fmt.Fprintf(os.Stdout, "[FINALLY RESTORE] close error %s\n", err.Error())
		}
	}()

	_, _ = fmt.Fprintf(os.Stdout, "[START RESTORE] read all message from snapshot\n")
	var totalRestored int

	decoder := json.NewDecoder(rClose)
	for decoder.More() {
		var data = &CommandPayload{}
		err := decoder.Decode(data)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stdout, "[END RESTORE] error decode data %s\n", err.Error())
			return err
		}

		if err := (*f.db).Put(data.Key, data.Value); err != nil {
			_, _ = fmt.Fprintf(os.Stdout, "[END RESTORE] error persist data %s\n", err.Error())
			return err
		}

		totalRestored++
	}

	// read closing bracket
	_, err := decoder.Token()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stdout, "[END RESTORE] error %s\n", err.Error())
		return err
	}

	_, _ = fmt.Fprintf(os.Stdout, "[END RESTORE] success restore %d messages in snapshot\n", totalRestored)
	return nil
}

// NewFSM raft.FSM implementation using badgerDB
func NewFSM(database *database.Database, cluster *utils.ClusterRoutine, dht *DhtRoutine) raft.FSM {
	return &FSM{
		db:      database,
		cluster: cluster,
		dht:     dht,
	}
}
