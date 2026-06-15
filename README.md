# memorabilia

## What is memorabilia?

memorabilia is striving to be an open-source, in-memory data structure store. Someday!

## How to Run

memorabilia can run in two modes:

- **Single-node mode** — the original behaviour. One process, one in-memory
  store, no replication. Data is lost if the process restarts.
- **Raft mode** — multiple processes form a cluster. Writes are replicated
  to a quorum of nodes before being acknowledged, and the cluster survives
  the loss of a minority of nodes. State also persists across restarts via
  on-disk Raft log and snapshot files.

The mode is selected automatically: if `--node-id` (or `MEMORABILIA_NODE_ID`)
is left unset, the server starts in single-node mode. Setting it enables Raft
mode.

### Prerequisites

- Go 1.22+
- (Optional, for testing) [`grpcurl`](https://github.com/fullstorydev/grpcurl)

### Build

```bash
go build -o bin/server ./cmd
```

Rebuild whenever the source changes — the binary is static, all per-node
configuration comes from flags or environment variables at runtime.

---

### Single-Node Mode (No Replication)

This is the simplest way to run memorabilia — useful for local development,
testing, or anywhere you don't need fault tolerance.

```bash
./bin/server --port=50051
```

That's it. No data directory is created, no extra ports are opened, and the
TTL cleanup job runs locally against this node's own store.

Test it:

```bash
grpcurl -plaintext -d '{"id":"foo","value":"bar","ttl":60000}' \
  127.0.0.1:50051 commands.Commands/Set

grpcurl -plaintext -d '{"id":"foo"}' \
  127.0.0.1:50051 commands.Commands/Get
```

---

### Raft Mode (Replicated Cluster)

In Raft mode, every node opens **three** ports:

| Port | Purpose |
|---|---|
| `--port` | gRPC — client reads and writes (`Set`, `Get`, `Delete`, ...) |
| `--raft-addr` | Raft peer-to-peer traffic (log replication, elections) |
| `--http-mgmt-addr` | HTTP cluster management (`/raft/join`, `/raft/leader`, `/raft/peers`) |

A cluster is formed by **bootstrapping exactly one node**, then having every
other node **join** through that node's HTTP management address.

#### 1. Bootstrap the first node

```bash
./bin/server \
  --node-id=n1 \
  --port=50051 \
  --raft-addr=127.0.0.1:7001 \
  --http-mgmt-addr=127.0.0.1:8081 \
  --data-dir=./data \
  --bootstrap
```

`--bootstrap` tells this node to form a brand-new single-node cluster and
elect itself leader. **Only ever pass `--bootstrap` on the very first run of
the very first node.** It is safe to leave it on a config file for that one
node — `replication.NewNode` checks for existing on-disk state and skips
re-bootstrapping on subsequent restarts — but it should never be passed when
joining a node to an existing cluster.

Verify it came up before continuing:

```bash
curl http://127.0.0.1:8081/raft/leader
# → 127.0.0.1:7001
```

#### 2. Join the second node

```bash
./bin/server \
  --node-id=n2 \
  --port=50052 \
  --raft-addr=127.0.0.1:7002 \
  --http-mgmt-addr=127.0.0.1:8082 \
  --data-dir=./data \
  --leader-http=127.0.0.1:8081
```

At startup, this node POSTs to `n1`'s `/raft/join` endpoint, which calls
`raft.AddVoter`. From this point on, `n1` won't commit a write unless `n2`
also acknowledges it.

#### 3. Join the third node

```bash
./bin/server \
  --node-id=n3 \
  --port=50053 \
  --raft-addr=127.0.0.1:7003 \
  --http-mgmt-addr=127.0.0.1:8083 \
  --data-dir=./data \
  --leader-http=127.0.0.1:8081
```

With 3 nodes, the cluster tolerates the loss of **1** node and still has a
quorum (2 of 3) to keep committing writes.

#### 4. Verify the cluster

```bash
curl http://127.0.0.1:8081/raft/peers
```

Returns a JSON array with all three servers and their Raft addresses.

#### 5. Test replication

```bash
# Write through the leader (n1)
grpcurl -plaintext -d '{"id":"foo","value":"bar","ttl":60000}' \
  127.0.0.1:50051 commands.Commands/Set

# Read from a follower (n2) — should return the replicated value
grpcurl -plaintext -d '{"id":"foo"}' \
  127.0.0.1:50052 commands.Commands/Get
```

If you `Set` against a follower's gRPC port, you'll get a
`FailedPrecondition` error containing the current leader's Raft address —
this is the expected "not the leader, retry here" behaviour.

---

### Restarting a Node

Point it at the **same** `--data-dir` and **omit** `--bootstrap`. The node
will replay its Raft log and snapshots from disk and rejoin the cluster at
the term and index it left off at.

### Resetting a Cluster

```bash
pkill -f bin/server
rm -rf ./data
```

Then bootstrap node 1 again as in step 1.

### A Note on Startup Logs

You may see lines like `Rollback failed: tx closed` from `raft-boltdb` on
startup. This is a known cosmetic issue in the BoltDB store initialization —
it logs a no-op rollback after a transaction has already committed
successfully. It does not indicate data loss or corruption and can be
ignored.

---

### Configuration Reference

Every flag can also be set via an environment variable. **Flags take
precedence over environment variables**, and environment variables take
precedence over the built-in defaults shown below.

| Flag | Env Var | Default | Mode | Description |
|---|---|---|---|---|
| `--port` | `MEMORABILIA_PORT` | `50051` | Both | gRPC server port for client traffic |
| `--ttl-cleanup-ms` | `MEMORABILIA_TTL_CLEANUP_MS` | `4000` | Both | Interval (ms) for the background TTL expiry cleanup job |
| `--node-id` | `MEMORABILIA_NODE_ID` | `""` | — | Unique node identifier (e.g. `n1`). **Setting this enables Raft mode.** Leave unset for single-node mode. |
| `--raft-addr` | `MEMORABILIA_RAFT_ADDR` | `0.0.0.0:7000` | Raft only | TCP address this node's Raft transport binds to |
| `--advertise-addr` | `MEMORABILIA_ADVERTISE_ADDR` | *(same as raft-addr)* | Raft only | Address other nodes dial to reach this one. Set when behind NAT, a load balancer, or in Docker where the bind address (`0.0.0.0`) isn't reachable from other containers |
| `--http-mgmt-addr` | `MEMORABILIA_HTTP_MGMT_ADDR` | `0.0.0.0:8081` | Raft only | Address for `/raft/join`, `/raft/leader`, `/raft/peers` |
| `--data-dir` | `MEMORABILIA_DATA_DIR` | `./data` | Raft only | Base directory for Raft log, stable store, and snapshots. A subdirectory named after `--node-id` is created automatically (e.g. `./data/n1`) |
| `--bootstrap` | `MEMORABILIA_BOOTSTRAP` | `false` | Raft only | Form a brand-new single-node cluster and self-elect as leader. Set only on the first run of the first node — never on join |
| `--leader-http` | `MEMORABILIA_LEADER_HTTP` | `""` | Raft only | HTTP management address of the cluster leader. Set on every node **except** the bootstrap node, so it can register via `/raft/join` at startup |

#### Example: configuring via environment variables

```bash
export MEMORABILIA_NODE_ID=n2
export MEMORABILIA_PORT=50052
export MEMORABILIA_RAFT_ADDR=127.0.0.1:7002
export MEMORABILIA_HTTP_MGMT_ADDR=127.0.0.1:8082
export MEMORABILIA_DATA_DIR=./data
export MEMORABILIA_LEADER_HTTP=127.0.0.1:8081

./bin/server
```

This is equivalent to the `--node-id=n2 --port=50052 ...` invocation in step
2 above — useful for Docker Compose, Kubernetes manifests, or systemd unit
files where environment variables are more natural than long argument lists.

#### Required combinations

| Scenario | Required flags/env |
|---|---|
| Single-node, default settings | *(none — just run the binary)* |
| Single-node, custom port | `--port` |
| Raft, first node | `--node-id`, `--raft-addr`, `--http-mgmt-addr`, `--data-dir`, `--bootstrap` |
| Raft, joining node | `--node-id`, `--raft-addr`, `--http-mgmt-addr`, `--data-dir`, `--leader-http` |
| Raft, node behind NAT/Docker | add `--advertise-addr` to either of the above |
