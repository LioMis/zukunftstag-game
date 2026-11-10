# Bubblewrap FHS environment for building and running Paket-Blitz.
# raylib's cgo build and miniaudio's dlopen of ALSA/PulseAudio find everything
# under /usr/lib like on any other distro. GL drivers come from the host's
# /run/opengl-driver.
#
# Build and enter it with ./fhs.sh (no flakes needed):
#   ./fhs.sh            interactive shell inside the FHS env
#   ./fhs.sh go run .   run one command inside it
{ pkgs ? import <nixpkgs> { } }:
pkgs.buildFHSEnv {
  name = "paket-blitz-env";
  targetPkgs = p: with p; [
    go
    gcc
    pkg-config
    git
    libGL
    libglvnd
    libxkbcommon
    xorg.libX11
    xorg.libXcursor
    xorg.libXext
    xorg.libXi
    xorg.libXinerama
    xorg.libXrandr
    xorg.libXrender
    xorg.libXxf86vm
    alsa-lib
    libpulseaudio
    # Headless screenshots (scripts/shoot.sh falls back to these without a display).
    xvfb-run
    mesa
    ffmpeg
  ];
  # Set the Go env here rather than in `profile`: /etc/profile is only read
  # by login shells, and `./fhs.sh go run .` starts no shell at all.
  runScript = pkgs.writeShellScript "paket-blitz-run" ''
    export CGO_ENABLED=1
    export GOFLAGS=-tags=x11
    if [ $# -eq 0 ]; then
      exec bash
    fi
    exec "$@"
  '';
}
