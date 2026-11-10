#!/usr/bin/env bash
# Build the game and replay an input script, writing screenshots to shots/.
#   scripts/shoot.sh scripts/delivery.txt [extra game flags...]
# Uses the real display if there is one, otherwise a virtual one (xvfb-run).
set -euo pipefail
cd "$(dirname "$0")/.."
script=${1:?usage: scripts/shoot.sh <script.txt> [flags]}
shift
go build -o bin/paket-blitz .
run=(./bin/paket-blitz -script "$script" -seed 1 -scores shots/scores.json "$@")
mkdir -p shots
rm -f shots/scores.json
if [[ -n "${DISPLAY:-}" || -n "${WAYLAND_DISPLAY:-}" ]]; then
	"${run[@]}"
else
	xvfb-run -a -s "-screen 0 1280x720x24" "${run[@]}"
fi 2>&1 | grep -v '^WARNING: FONT'
