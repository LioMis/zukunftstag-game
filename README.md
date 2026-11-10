# Paket-Blitz

A small 3D delivery game for Zukunftstag. Drive through town, pick up parcels (yellow), deliver
them (green), collect coins. Each delivery adds time. Highest Franken wins.

## Run (NixOS)

```bash
nix-shell --run 'go run .'
```

Or open the shell and work in it:

```bash
nix-shell
```

```bash
go run .
```

`shell.nix` enters a bubblewrap FHS environment (`fhs.nix`) and needs no flakes. With flakes
enabled, `nix develop` and `nix run` do the same.

Controls: arrow keys or WASD, Enter/Space to confirm, Esc to quit. Gamepads work too.
Flags: `-fullscreen`, `-w`/`-h`, `-scores <file>`.

## Screenshots without playing

```bash
scripts/shoot.sh scripts/delivery.txt
```

writes PNGs to `shots/`. See `CLAUDE.md` for the script format.

## Credits

- 3D models and sounds: [Kenney](https://kenney.nl), CC0, from the City Builder, Racing and 3D
  Platformer starter kits on github.com/KenneyNL.
- Font: Lilita One by Juan Montoreano, SIL Open Font License (`assets/fonts/license.txt`).
