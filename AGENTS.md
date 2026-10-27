# Working in Open-Experiment

This repository follows the policy handbook at
https://github.com/jasonelias144-svg/policy-handbook. Read its orientation
(README.md) before your first session here. The rules below add to the
handbook; they never loosen it.

## Every session

1. Read `STATUS.md` and the open pull requests (handbook P1).
2. Work on a branch; never push to `main` (P3).
3. End by pushing, updating `STATUS.md` on your branch, and opening or updating
   the pull request (P1).

## Who decides

The founder decides what this project is for and who works on it. The technical
lead (Claude) decides everything about the code, and records an approval on
every pull request before it merges (P4). Intent questions go to the founder (P5).

## This repository

- **What it is:** experiments in prompt-minimal, open-ended generation: the
  code recovered from an earlier undirected session (`experiment-in-choice/`,
  14 small projects), and an autogram search (`autogram/`).
- **Run:** each project's own README. Most are a Go simulation that streams
  JSON over a WebSocket to an HTML page, with Docker files.
- **Test:** `cd experiment-in-choice/14-chord-for-real/backend && go test ./...`
  (the only project with unit tests). CI runs it.
- **Merge method:** squash
- **Release:** none. Nothing here is deployed.

## Rules specific to this repository

- **This is a record of experiments.** Keep each one as it was found or made,
  and add notes rather than rewriting history. A new experiment gets its own
  folder with a README that says what it is, why, and whether it runs.
- **Say honestly whether each thing runs,** and what was and wasn't tested,
  as the existing READMEs do.

## Always

- Say who made each change: `Made-by: <name> (AI|human)` in commits, and the
  "Made by" field in pull requests (P3).
- No secrets or personal data anywhere in the repository (P8).
- Report outcomes honestly: what you checked, and what you did not.
