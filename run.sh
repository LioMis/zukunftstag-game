#!/usr/bin/env bash
# Build and start Paket-Blitz. Just run:   ./run.sh
# On NixOS it re-runs itself inside the FHS env (fhs.nix). Arguments go to the game,
# e.g. ./run.sh -fullscreen
set -euo pipefail
cd "$(dirname "$0")"

if [ -e /etc/NIXOS ] && [ -z "${PAKET_BLITZ_FHS:-}" ]; then
	exec ./fhs.sh ./run.sh "$@"
fi

if [ ! -x paket-blitz ] || [ -n "$(find . -path ./shots -prune -o \( -name '*.go' -o -name go.mod -o -name go.sum \) -newer paket-blitz -print -quit)" ]; then
	echo "==> Building (the first build compiles raylib and takes a minute)"
	go build -tags x11 -o paket-blitz .
fi
exec ./paket-blitz "$@"
