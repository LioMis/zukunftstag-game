# Paket-Blitz

A small 3D delivery game for Zukunftstag (kids visiting for a day). Kids drive a truck through a
town, pick up parcels, deliver them, collect coins, and race for the highscore. Non-violent.

An adult guides the session while kids aged about 10–14 watch, seeing agentic engineering for the
first time. Talk in their language (usually German): short, snappy, plain sentences, a few lines
per message. Say what you're doing and why. No jargon without a one-line explanation. Don't make
suggestions or offer options; the adult steers.

## Stack

- Go + raylib-go (cgo, OpenGL 3.3).
- Only target: NixOS. `./run.sh` builds and starts the game; on NixOS it enters the Nix shell
  (`shell.nix`) by itself. `fhs.nix` is a bubblewrap FHS env alternative:
  `$(nix-build fhs.nix --no-out-link)/bin/paket-blitz-fhs ./run.sh`. The setup is copied from
  github.com/lmis/wheel-loader, which works on the target machine; keep it in line with that
  rather than inventing a different one. Flakes are disabled on the target machine.

## Code map

| File | What |
| --- | --- |
| `main.go` | flags, window, main loop, screenshot capture |
| `game.go` | game states (title, play, game over, scores), jobs, coins, camera |
| `draw.go` | all rendering: 3D scene, markers, HUD, menus |
| `city.go` | town grid generation, road tile rotation, collision |
| `car.go` | truck driving physics, banana spin |
| `monkeys.go` | monkeys walking the roads, throwing bananas the truck slips on |
| `input.go` | keyboard/gamepad input and the replay script format |
| `assets.go` | model, font, sound loading; the lit shader |
| `scores.go` | highscore table (`scores.json`) |

`assets/custom/monkey.glb` is a cartoon monkey (about 1.75 m tall, faces +Z), loaded as `"monkey"`
and drawn by `monkeys.go`. To change its look, edit `scripts/make-monkey.py` and rerun it.

## Seeing the game: replay scripts and screenshots

You can't see the window, so test by replaying scripted input and looking at screenshots:

```bash
scripts/shoot.sh scripts/delivery.txt
```

This builds the game, replays the script with a fixed seed and fixed 60 fps timestep, and writes
PNGs to `shots/`. Read them with the Read tool and show them to the user. The script format is
documented on `Script` in `input.go` (`hold`, `release`, `enter`, `type`, `warp`, `shot`, `quit`).
`warp pickup|dropoff` teleports the truck next to the current target, `warp banana` parks it before a
banana on the road, so you don't have to steer.

## Rules

- Assets: CC0 only (Kenney kits from github.com/KenneyNL), fonts under OFL. No "free for
  non-commercial", no CC-BY-SA, nothing that forces a license on the game. Record new assets in
  `README.md` under Credits.
- Tests: don't write them upfront. Write a test only for a bug you found while screenshotting, and
  name the bug in a comment on the test.
- To add a feature, use the `feature` skill (plan → build → screenshot → fix).
