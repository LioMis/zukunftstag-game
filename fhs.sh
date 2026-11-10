#!/usr/bin/env bash
# Run a command (or, without arguments, a shell) inside the FHS env from fhs.nix.
#   ./fhs.sh            interactive shell
#   ./fhs.sh go run .   run the game
set -euo pipefail
cd "$(dirname "$0")"
env=$(nix-build fhs.nix --no-out-link)
PAKET_BLITZ_FHS=1 exec "$env/bin/paket-blitz-fhs" "$@"
