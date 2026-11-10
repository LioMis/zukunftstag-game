// Paket-Blitz: a small delivery game for Zukunftstag.
package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	assetDir := flag.String("assets", "assets", "asset directory")
	scoresPath := flag.String("scores", "scores.json", "highscore file")
	script := flag.String("script", "", "replay input from this script instead of the keyboard (headless screenshots)")
	shotDir := flag.String("shots", "shots", "directory for screenshots written by the script")
	seed := flag.Int64("seed", 0, "random seed (0 = time based)")
	width := flag.Int("w", 1280, "window width")
	height := flag.Int("h", 720, "window height")
	fullscreen := flag.Bool("fullscreen", false, "start fullscreen")
	flag.Parse()

	if *seed == 0 {
		*seed = time.Now().UnixNano()
	}
	rng := rand.New(rand.NewSource(*seed))

	var sc *Script
	if *script != "" {
		var err error
		if sc, err = LoadScript(*script); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if err := os.MkdirAll(*shotDir, 0o755); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	rl.SetTraceLogLevel(rl.LogWarning)
	flags := uint32(rl.FlagMsaa4xHint)
	if sc == nil {
		flags |= rl.FlagVsyncHint | rl.FlagWindowResizable
	}
	rl.SetConfigFlags(flags)
	rl.InitWindow(int32(*width), int32(*height), "Paket-Blitz")
	defer rl.CloseWindow()
	if *fullscreen {
		rl.ToggleFullscreen()
	}
	rl.SetExitKey(rl.KeyEscape)

	audio := sc == nil
	if audio {
		rl.InitAudioDevice()
		audio = rl.IsAudioDeviceReady()
		defer rl.CloseAudioDevice()
	}
	a := LoadAssets(*assetDir, audio)
	g := NewGame(a, rng, *scoresPath)
	if audio {
		rl.PlayMusicStream(a.Engine)
		rl.SetMusicVolume(a.Engine, 0)
	}

	for frame := 0; !rl.WindowShouldClose(); frame++ {
		var in Input
		dt := float32(1.0 / 60)
		if sc != nil {
			in = sc.Input(frame)
		} else {
			in = KeyboardInput()
			dt = clamp(rl.GetFrameTime(), 0, 0.05)
		}
		g.Update(in, dt)
		if audio {
			rl.UpdateMusicStream(a.Engine)
			v := float32(0)
			if g.State == StatePlay {
				v = 0.15 + 0.25*clamp(abs32(g.Car.Speed)/carMaxSpeed, 0, 1)
			}
			rl.SetMusicVolume(a.Engine, v)
			rl.SetMusicPitch(a.Engine, 0.8+0.7*clamp(abs32(g.Car.Speed)/carMaxSpeed, 0, 1))
		}

		rl.BeginDrawing()
		g.Draw()
		if sc != nil && len(sc.Shots) > 0 {
			rl.DrawRenderBatchActive() // flush 2D batch so the capture includes the HUD
			for _, name := range sc.Shots {
				img := rl.LoadImageFromScreen()
				p := filepath.Join(*shotDir, name+".png")
				rl.ExportImage(*img, p)
				rl.UnloadImage(img)
				fmt.Println("wrote", p)
			}
		}
		rl.EndDrawing()
		if sc != nil && sc.Quit {
			break
		}
	}
}

func abs32(v float32) float32 {
	if v < 0 {
		return -v
	}
	return v
}
