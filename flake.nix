{
  description = "Paket-Blitz – a small 3D delivery game for Zukunftstag";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

  outputs = { self, nixpkgs }:
    let
      system = "x86_64-linux";
      pkgs = nixpkgs.legacyPackages.${system};

      # Bubblewrap FHS environment: raylib's cgo build and miniaudio's
      # dlopen of ALSA/PulseAudio find everything under /usr/lib like on
      # any other distro. GL drivers come from the host's /run/opengl-driver.
      fhs = pkgs.buildFHSEnv {
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
        profile = ''
          export CGO_ENABLED=1
          export GOFLAGS=-tags=x11
        '';
        runScript = "bash";
      };
    in
    {
      devShells.${system}.default = fhs.env;

      # `nix run` from the repo root: builds and starts the game.
      apps.${system}.default = {
        type = "app";
        program = toString (pkgs.writeShellScript "paket-blitz" ''
          exec ${fhs}/bin/paket-blitz-env -c 'go run . "$@"' paket-blitz "$@"
        '');
      };
    };
}
