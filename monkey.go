package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	colFur  = rl.NewColor(120, 72, 38, 255)
	colSkin = rl.NewColor(232, 190, 140, 255)
	colEye  = rl.NewColor(25, 20, 20, 255)
)

// Monkey walks the roads, steals coins and drops bananas.
type Monkey struct {
	Pos     rl.Vector3
	Heading float32 // radians, 0 faces +Z (like the car)
	Phase   float32 // walk animation
}

// part is one ball of a monkey, in the monkey's own frame (+Z is its front, units are metres).
type part struct {
	off, size rl.Vector3
	col       rl.Color
}

// monkeyParts builds the monkey from squashed balls. swing moves arms and legs while walking.
func monkeyParts(swing float32) []part {
	v := rl.NewVector3
	ps := []part{
		{v(0, 0.75, 0), v(0.45, 0.55, 0.4), colFur},              // body
		{v(0, 0.72, 0.25), v(0.3, 0.38, 0.18), colSkin},          // belly
		{v(0, 1.45, 0), v(0.38, 0.36, 0.36), colFur},             // head
		{v(0, 1.4, 0.24), v(0.28, 0.24, 0.16), colSkin},          // face
		{v(0, 1.32, 0.33), v(0.17, 0.11, 0.1), colSkin},          // muzzle
		{v(-0.11, 1.5, 0.36), v(0.06, 0.07, 0.04), colEye},       // eyes
		{v(0.11, 1.5, 0.36), v(0.06, 0.07, 0.04), colEye},        //
		{v(-0.38, 1.5, 0), v(0.15, 0.15, 0.07), colFur},          // ears
		{v(0.38, 1.5, 0), v(0.15, 0.15, 0.07), colFur},           //
		{v(-0.41, 1.5, 0.03), v(0.09, 0.09, 0.05), colSkin},      // inner ears
		{v(0.41, 1.5, 0.03), v(0.09, 0.09, 0.05), colSkin},       //
		{v(-0.44, 0.8, swing*0.2), v(0.12, 0.38, 0.12), colFur},  // arms
		{v(0.44, 0.8, -swing*0.2), v(0.12, 0.38, 0.12), colFur},  //
		{v(-0.2, 0.28, -swing*0.15), v(0.14, 0.3, 0.14), colFur}, // legs
		{v(0.2, 0.28, swing*0.15), v(0.14, 0.3, 0.14), colFur},   //
	}
	// Tail: a row of small balls curling up behind.
	for i := 0; i < 16; i++ {
		t := float64(i) / 15
		a := t * math.Pi * 1.3
		ps = append(ps, part{v(0, 0.5+float32(0.6*(1-math.Cos(a))), -0.35-float32(0.45*math.Sin(a))),
			v(0.075, 0.075, 0.075), colFur})
	}
	return ps
}

func (m *Monkey) Draw(a *Assets) {
	rl.DrawCylinder(rl.NewVector3(m.Pos.X, 0.07, m.Pos.Z), 0.6, 0.6, 0.01, 16, rl.NewColor(0, 0, 0, 70))
	swing := float32(math.Sin(float64(m.Phase)))
	bob := float32(math.Abs(float64(swing))) * 0.08
	place := rl.MatrixMultiply(rl.MatrixRotateY(m.Heading), rl.MatrixTranslate(m.Pos.X, m.Pos.Y+bob, m.Pos.Z))
	for _, p := range monkeyParts(swing) {
		a.Ball(p.off, p.size, p.col, place)
	}
}
