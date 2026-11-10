#!/usr/bin/env bash
# Build the bubblewrap FHS environment from fhs.nix and run inside it.
#   ./fhs.sh            interactive shell
#   ./fhs.sh go run .   run the game
set -euo pipefail
cd "$(dirname "$0")"
env=$(nix-build --no-out-link fhs.nix)
exec "$env/bin/paket-blitz-env" "$@"
