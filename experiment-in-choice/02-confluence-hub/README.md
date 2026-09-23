# The Confluence Hub

A multi-room chat hub for human–human, human–AI and AI–AI conversation, which was the one direction the user explicitly asked for. It has a Go WebSocket router with rooms, a web 'collaboration desk' with discipline channels, a React Native room component, and a Python AI-agent plugin that joins a room and replies to messages.

![screenshot](screenshot.png)

## Run it

Without Docker (Go 1.21+), from this folder, in two terminals:

```sh
(cd backend && go run .)                       # WebSocket on :8080/ws?room=…&type=human|ai&id=…
(cd frontend && python3 -m http.server 3000)   # then open http://localhost:3000
```

With Docker: `docker compose up --build` from this folder (not tested here; see the top-level README).

AI agent: `pip install websockets && python plugins/agent_plugin.py` (reads `HUB_WS_URL`, default localhost). `mobile/Room.js` is a component to drop into a React Native app, not a standalone app.

## What was checked

Round trip tested: a message typed in the browser appeared in the chat, the Python agent received it and posted a reply into the room. Room.js parses as JSX. It was not run inside a mobile app.

## Repairs to the transcript's code

- `backend/main.go`: gorilla/websocket import restored; comment/declaration splits.
- `plugins/agent_plugin.py`: import and comment/code line splits. It now reads `HUB_WS_URL` from the environment, which docker-compose already sets for each agent. Previously it ignored that and connected to localhost, so it could never reach the hub from inside Docker.
- Dockerfiles: instruction splits; `COPY go.mod go.sum ./`.

`go.mod` / `go.sum` were generated here, following the transcript's own `go mod init` / `go get` steps. Every other change is visible with `git diff` against the verbatim-extraction commit.

## Where each file came from

| File | Transcript lines |
|---|---|
| `backend/main.go` | 2378–2535 |
| `mobile/Room.js` | 2541–2618 |
| `plugins/agent_plugin.py` | 2653–2710 |
| `frontend/index.html` | 2754–2941 |
| `docker-compose.yml` | 2977–3032 |
| `backend/Dockerfile` | 3119–3130 |
| `plugins/Dockerfile` | 3134–3138 |
