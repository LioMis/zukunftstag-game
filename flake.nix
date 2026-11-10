# Optional: only for machines with flakes enabled. Without flakes use shell.nix.
{
  description = "Paket-Blitz – a small 3D delivery game for Zukunftstag";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

  outputs = { self, nixpkgs }:
    let
      system = "x86_64-linux";
      pkgs = nixpkgs.legacyPackages.${system};

      fhs = import ./fhs.nix { inherit pkgs; };
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
