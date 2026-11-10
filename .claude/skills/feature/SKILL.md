---
name: feature
description: Plan, build and check a new game feature together with a kid, step by step - plan, OK, implement, screenshot-test, fix, show. Use whenever someone asks for a new feature or change in Paket-Blitz ("mach dass...", "ich will ...", "/feature ...").
---

# Feature: plan → build → test with screenshots

An adult guides the session; kids watch and call out ideas. The point is that they see the loop an
engineer runs with an agent: **understand → plan → build → check → fix → show**.
Say which step you're in ("Schritt 2: Plan"). Every message is short and snappy, in their language:
a few lines, no lists of options.

**No suggestions.** Don't offer ideas, alternatives, examples or "what next?" options, and don't
propose answers to your own questions. The adult steers; you do what's asked and report.

## 1. Understand

Repeat the idea back in one sentence. If something decides how it plays (how big, how fast, what
happens when…), ask at most two short open questions. Pick everything else sensibly yourself.

## 2. Plan

Write a plan of 3–6 short points:

- what the player will see and do differently,
- which files change (one phrase each, e.g. "`game.go`: Turbo-Zeit zählen"),
- **how we'll check it works**: which screenshots you'll take and what should be in them.

Then stop and ask: "Passt der Plan?" Wait for an OK. Let them change the plan.

## 3. Build

Make the smallest change that does the plan. Then:

```bash
go vet . && go build -o bin/paket-blitz .
```

If it fails, read the error, explain it in one sentence ("Ich habe vergessen, …"), fix, repeat.

## 4. Check with screenshots

Write a replay script `scripts/<feature>.txt` (format: `Script` in `input.go`) that reaches the
moment the feature matters and takes `shot` frames there; use `warp` to skip driving. Run it:

```bash
scripts/shoot.sh scripts/<feature>.txt
```

Open every screenshot and compare it with what the plan said you'd see. Say what you see, honestly,
including what looks wrong.

## 5. Fix

For each thing that's wrong: say what's wrong and why, fix it, take the screenshots again. If the
bug can come back unnoticed, add a small Go test for it (`*_test.go`, comment names the bug) and
run `go test .`. No tests for things that already worked.

## 6. Show

Show the final screenshots and tell them to try it: `./run.sh`. Then, if they're happy, commit with a short message saying what the feature does.
