# Code recovered from "experiment in choice"

These folders hold all the code from `experiment_in_choice.txt`, a transcript of
an earlier undirected-generation session. In that session the user mostly
replied "➡️" and the model picked each next step. The transcript marks 13
project cycles, and each folder here is one cycle, in order.

Nearly every cycle has the same shape: a Go program that simulates something
and streams it as JSON over a WebSocket, an HTML page that visualizes (and
sometimes sonifies) the stream, a Dockerfile, and a docker-compose file.

| # | Project | Simulates | Runs? |
|---|---|---|---|
| 01 | [unbound-net](01-unbound-net/) | "resonance" heartbeat → audio-visual page, plus Python/Node/Prometheus/K8s extras | ✅ |
| 02 | [confluence-hub](02-confluence-hub/) | multi-room chat for humans and AI agents | ✅ round trip tested |
| 03 | [biomimetic-automata](03-biomimetic-automata/) | cellular automaton on a torus, with audio | ✅ |
| 04 | [quantum-circuit](04-quantum-circuit/) | 2-qubit state vector (H + CNOT → Bell state) | ✅ |
| 05 | [neuro-evolution](05-neuro-evolution/) | small neural graph with random weight drift | ✅ |
| 06 | [market-graph](06-market-graph/) | three AMM pools and the arbitrage gap between them | ✅ |
| 07 | [consensus-cluster](07-consensus-cluster/) | 5-node ledger with faults and quorum | ✅ |
| 08 | [ray-marcher](08-ray-marcher/) | SDF torus rendered to ASCII | ✅ |
| 09 | [living-codex](09-living-codex/) | "seven planes" symbolic state machine | ✅ |
| 10 | [kernel-lineage](10-kernel-lineage/) | state engine + SQLite ledger + analyzer + alert manager | ✅ locally, ⚠️ Docker volume issue |
| 11 | [adaptive-graph](11-adaptive-graph/) | graph that rewires itself by node potential, with audio | ✅ |
| 12 | [work-stealing-scheduler](12-work-stealing-scheduler/) | 4-worker work-stealing pool | ✅ |
| 13 | [dht-chord-ring](13-dht-chord-ring/) | consistent-hashing ring (routing hops are faked) | ✅ |

Each project folder has its own README with how to run it, a screenshot, what
was checked, the repairs made, and the transcript line range of every file.

## What "runs" means here

For every project I:

1. built the Go backend and ran `go vet`;
2. started it, connected over WebSocket, and checked that it streams valid JSON
   containing every field the front-end reads;
3. opened the front-end in headless Chromium against the live backend for
   about 6 seconds and recorded page errors, console errors, and frames
   received. None of the 13 threw a JavaScript error. The only console
   message was a single 404 on project 1 that didn't recur on a rerun. 12
   pages received a live stream; the chat hub only sends when someone posts,
   so it was tested separately by sending a message;
4. replayed each backend's Dockerfile build steps in a clean folder, copying
   in only the files the Dockerfile copies.

Also run directly: the Python terminal scripts, the Prometheus exporter, the
Node.js stream, and project 10's analyzer and alert manager against a real
database the backend wrote.

**Not run:** Docker itself (no Docker daemon was available), docker-compose,
the Kubernetes manifests, the Grafana dashboard import, and the React Native
component inside an app. Those were checked only for syntax and consistency.
"Runs" also doesn't mean the simulations are faithful. Several are
decorative. For example, the DHT's "routing hops" are random, and the
"neuro-evolution" never selects. The project READMEs say where.

## What was wrong with the transcript's code

The transcript is a copy of rendered chat, and the copying damaged the code in
consistent ways:

- **Merged lines.** A line starting at column 0 was often glued onto the end
  of the line before it: `import osimport sys`, `// …viewportstype Foo struct {`,
  `# …PythonFROM python:3.11-slim`, `apiVersion: v1kind: ConfigMap`. The
  dangerous cases are where code ends up *inside a comment*. For example, in
  `run.sh` the `while true; do` line was absorbed into a comment, so the script
  still parsed but ran once instead of looping.
- **Truncated links.** Every `github.com/gorilla/websocket` import became
  `"://github.com"`; the SQLite driver import did too.
- **A deleted bracket.** `let lastQueueSizes = [0, 0, 0, 0];` became
  `let lastQueueSizes =;`, which broke the whole page.

Separately, the model that wrote the code made ordinary bugs: `Math.Cos` in Go,
`time.sleep` in Go, `rand.Int64()` (not in `math/rand`), two missing `math`
imports, a garbled `scheduler CentralScheduler.Mu.Lock()`, two `main()`
functions in one package, a Python file mixing tabs and spaces, a plugin that
ignored the address docker-compose gave it, and Dockerfiles that could never
build (they copy a `go.mod` that didn't exist and never copy a `go.sum`).

## How the history is laid out

- **Commit 1 — verbatim extraction.** Every file exactly as it appears in the
  transcript, with only surrounding blank lines trimmed. `manifest.json` records
  each file's source line range (its paths are from that commit, so one file,
  `10-kernel-lineage/backend/alertmanager.go`, has since moved).
- **Commit 2 — repairs**, plus the generated `go.mod`/`go.sum` files, READMEs
  and screenshots.

So `git diff <commit 1> <commit 2> -- '*.go' '*.py' '*.js' '*.html' '*.yml' '*.yaml' '*Dockerfile*' '*.sh'`
shows every change made to the original code, and nothing else.

Earlier versions that the transcript itself replaced are kept, not deleted:
five earlier front-ends in `01-unbound-net/frontend/history/`, one each in 03
and 11, and project 10's first backend and patch snippets in
`10-kernel-lineage/history/`.
