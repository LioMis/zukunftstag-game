# Paket-Blitz

A small 3D delivery game for Zukunftstag (kids visiting for a day). Kids drive a truck through a
town, pick up parcels, deliver them, collect coins, and race for the highscore. Non-violent.

The people prompting you here are often kids aged about 10–14 who are seeing agentic engineering
for the first time. Talk in their language (usually German), in short plain sentences, and explain
what you're doing and why at each step. No jargon without a one-line explanation.

## Stack

- Go + raylib-go (cgo, OpenGL 3.3). Always build with `-tags x11` (no Wayland headers needed).
- Only target: NixOS, through the bubblewrap FHS shell in `fhs.nix`: `nix-shell`, then `go run .`
  (or `nix-shell --run 'go run .'`). Flakes are disabled on the target machine; `flake.nix` is an
  optional extra. The shell sets `GOFLAGS=-tags=x11`.

## Code map

| File | What |
| --- | --- |
| `main.go` | flags, window, main loop, screenshot capture |
| `game.go` | game states (title, play, game over, scores), jobs, coins, camera |
| `draw.go` | all rendering: 3D scene, markers, HUD, menus |
| `city.go` | town grid generation, road tile rotation, collision |
| `car.go` | truck driving physics |
| `input.go` | keyboard/gamepad input and the replay script format |
| `assets.go` | model, font, sound loading; the lit shader |
| `scores.go` | highscore table (`scores.json`) |

## Seeing the game: replay scripts and screenshots

You can't see the window, so test by replaying scripted input and looking at screenshots:

```bash
scripts/shoot.sh scripts/delivery.txt
```

This builds the game, replays the script with a fixed seed and fixed 60 fps timestep, and writes
PNGs to `shots/`. Read them with the Read tool and show them to the user. The script format is
documented on `Script` in `input.go` (`hold`, `release`, `enter`, `type`, `warp`, `shot`, `quit`).
`warp pickup|dropoff` teleports the truck next to the current target, so you don't have to steer.

## Rules

- Assets: CC0 only (Kenney kits from github.com/KenneyNL), fonts under OFL. No "free for
  non-commercial", no CC-BY-SA, nothing that forces a license on the game. Record new assets in
  `README.md` under Credits.
- Tests: don't write them upfront. Write a test only for a bug you found while screenshotting, and
  name the bug in a comment on the test.
- To add a feature, use the `feature` skill (plan → build → screenshot → fix).
