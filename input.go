package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Input is one frame's worth of player input.
type Input struct {
	Up, Down, Left, Right bool
	Enter, Back           bool // edge-triggered
	Chars                 []rune
	Warp                  string // script only: teleport the car, see Script
}

func KeyboardInput() Input {
	var in Input
	in.Up = rl.IsKeyDown(rl.KeyUp) || rl.IsKeyDown(rl.KeyW)
	in.Down = rl.IsKeyDown(rl.KeyDown) || rl.IsKeyDown(rl.KeyS)
	in.Left = rl.IsKeyDown(rl.KeyLeft) || rl.IsKeyDown(rl.KeyA)
	in.Right = rl.IsKeyDown(rl.KeyRight) || rl.IsKeyDown(rl.KeyD)
	in.Enter = rl.IsKeyPressed(rl.KeyEnter) || rl.IsKeyPressed(rl.KeyKpEnter) || rl.IsKeyPressed(rl.KeySpace)
	in.Back = rl.IsKeyPressed(rl.KeyBackspace)
	for r := rl.GetCharPressed(); r != 0; r = rl.GetCharPressed() {
		in.Chars = append(in.Chars, rune(r))
	}
	// Gamepad 0, if any.
	if rl.IsGamepadAvailable(0) {
		ax := rl.GetGamepadAxisMovement(0, rl.GamepadAxisLeftX)
		in.Left = in.Left || ax < -0.3 || rl.IsGamepadButtonDown(0, rl.GamepadButtonLeftFaceLeft)
		in.Right = in.Right || ax > 0.3 || rl.IsGamepadButtonDown(0, rl.GamepadButtonLeftFaceRight)
		in.Up = in.Up || rl.IsGamepadButtonDown(0, rl.GamepadButtonRightTrigger2) || rl.IsGamepadButtonDown(0, rl.GamepadButtonRightFaceDown)
		in.Down = in.Down || rl.IsGamepadButtonDown(0, rl.GamepadButtonLeftTrigger2) || rl.IsGamepadButtonDown(0, rl.GamepadButtonRightFaceLeft)
		in.Enter = in.Enter || rl.IsGamepadButtonPressed(0, rl.GamepadButtonMiddleRight) || rl.IsGamepadButtonPressed(0, rl.GamepadButtonRightFaceDown)
	}
	return in
}

// Script replays input for headless screenshots. Format, one command per line:
//
//	<frame> hold up|down|left|right ...
//	<frame> release up|down|left|right ...
//	<frame> enter | back
//	<frame> type TEXT
//	<frame> shot NAME          (writes NAME.png into the shot directory)
//	<frame> warp pickup|dropoff  (teleport the car next to the current target)
//	<frame> warp affe|banane     (park the car 7 m behind the first monkey or banana, facing it)
//	<frame> quit
type Script struct {
	events map[int][][]string
	held   map[string]bool
	Shots  []string // shots requested this frame
	Quit   bool
}

func LoadScript(path string) (*Script, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	s := &Script{events: map[int][][]string{}, held: map[string]bool{}}
	sc := bufio.NewScanner(f)
	for line := 1; sc.Scan(); line++ {
		fields := strings.Fields(sc.Text())
		if len(fields) == 0 || strings.HasPrefix(fields[0], "#") {
			continue
		}
		if len(fields) < 2 {
			return nil, fmt.Errorf("%s:%d: need <frame> <command>", path, line)
		}
		fr, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("%s:%d: %v", path, line, err)
		}
		s.events[fr] = append(s.events[fr], fields[1:])
	}
	return s, sc.Err()
}

func (s *Script) Input(frame int) Input {
	var in Input
	s.Shots = s.Shots[:0]
	for _, ev := range s.events[frame] {
		switch ev[0] {
		case "hold":
			for _, k := range ev[1:] {
				s.held[k] = true
			}
		case "release":
			for _, k := range ev[1:] {
				delete(s.held, k)
			}
		case "enter":
			in.Enter = true
		case "back":
			in.Back = true
		case "type":
			in.Chars = []rune(strings.Join(ev[1:], " "))
		case "shot":
			s.Shots = append(s.Shots, ev[1])
		case "warp":
			in.Warp = ev[1]
		case "quit":
			s.Quit = true
		}
	}
	in.Up, in.Down, in.Left, in.Right = s.held["up"], s.held["down"], s.held["left"], s.held["right"]
	return in
}
