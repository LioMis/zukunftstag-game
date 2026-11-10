# nix-shell            -> bubblewrap FHS shell, then `go run .`
# nix-shell --run 'go run .'
{ pkgs ? import <nixpkgs> { } }:
(import ./fhs.nix { inherit pkgs; }).env
