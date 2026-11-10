# A bubblewrap FHS environment with everything the game needs in the usual
# places (/usr/include, /usr/lib), so ./run.sh and go run . work unchanged.
#
#   ./run.sh                                  (enters this env by itself on NixOS)
#   ./fhs.sh go run .                         (run one command in it)
#   ./fhs.sh                                  (an interactive shell in it)
#   nix-shell -A env fhs.nix                  (same, via nix-shell)
#
# The host's OpenGL driver (/run/opengl-driver on NixOS) is used as is.
{ pkgs ? import <nixpkgs> { } }:

pkgs.buildFHSEnv {
  name = "paket-blitz-fhs";
  targetPkgs = p: with p; [
    go
    gcc
    pkg-config
    git
    libGL
    libglvnd
    mesa
    libx11
    libxcursor
    libxrandr
    libxinerama
    libxi
    libxext
    libxrender
    libxfixes
    libxkbcommon
    # miniaudio (raylib's audio) dlopens these
    alsa-lib
    libpulseaudio
    # headless screenshots: scripts/shoot.sh uses xvfb-run when there is no display
    xvfb-run
  ];
  # raylib's build needs the headers, which live in the -dev outputs
  extraOutputsToInstall = [ "dev" ];
  # with arguments, run them (e.g. ./run.sh); without, an interactive shell.
  # The Go env is set here, not in `profile`, so it also applies without a login shell.
  runScript = pkgs.writeShellScript "paket-blitz-fhs-run" ''
    export CGO_ENABLED=1
    export GOFLAGS=-tags=x11
    if [ $# -eq 0 ]; then exec bash; fi
    exec "$@"
  '';
  profile = ''
    export PAKET_BLITZ_FHS=1
  '';
}
