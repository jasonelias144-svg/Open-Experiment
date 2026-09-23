# Consensus Cluster

A 5-node replicated ledger with voting rounds and fault injection. When fewer than a quorum of nodes are live, the page shows the cluster halting.

![screenshot](screenshot.png)

## Run it

Without Docker (Go 1.21+), from this folder, in two terminals:

```sh
(cd backend && go run .)                       # WebSocket on :8080/consensus-stream
(cd frontend && python3 -m http.server 3000)   # then open http://localhost:3000
```

With Docker: `docker compose up --build` from this folder (not tested here; see the top-level README).

## What was checked

Ran and the page received 17 frames with no JS errors. (The transcript contains this compose file twice, identically, so it's kept once.)

## Repairs to the transcript's code

- gorilla/websocket import; comment/declaration splits; Dockerfile fixes.

`go.mod` / `go.sum` were generated here, following the transcript's own `go mod init` / `go get` steps. Every other change is visible with `git diff` against the verbatim-extraction commit.

## Where each file came from

| File | Transcript lines |
|---|---|
| `backend/main.go` | 5567–5705 |
| `frontend/index.html` | 5738–5899 |
| `backend/Dockerfile` | 5931–5941 |
| `docker-compose.yml` | 6064–6091 |
