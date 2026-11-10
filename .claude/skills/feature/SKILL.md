---
name: feature
description: Plan, build and check a new game feature together with a kid, step by step - plan, OK, implement, screenshot-test, fix, show. Use whenever someone asks for a new feature or change in Paket-Blitz ("mach dass...", "ich will ...", "/feature ...").
---

# Feature: plan → build → test with screenshots

The person is probably a kid seeing agentic engineering for the first time. The point is that they
see the loop an engineer runs with an agent: **understand → plan → build → check → fix → show**.
Say out loud which step you're in ("Schritt 2: Plan"). Keep each message short, in their language.

## 1. Understand

Repeat the idea back in one sentence. If something decides how it plays (how big, how fast, what
happens when…), ask at most two questions, each with a suggested answer they can just accept.
Don't ask about things you can pick sensibly yourself.

## 2. Plan

Write a plan of 3–6 short points:

- what the player will see and do differently,
- which files change (one phrase each, e.g. "`game.go`: Turbo-Zeit zählen"),
- **how we'll check it works**: which screenshots you'll take and what should be in them.

Then stop and ask: "Passt der Plan?" Wait for an OK. Let them change the plan.

## 3. Build

Make the smallest change that does the plan. Then:

```bash
go vet -tags x11 . && go build -tags x11 -o paket-blitz .
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
run `go test -tags x11 .`. No tests for things that already worked.

## 6. Show

Show the final screenshots and tell them to try it: `./fhs.sh go run .`. Ask what they'd
change next. Then, if they're happy, commit with a short message saying what the feature does.
