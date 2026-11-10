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
	// Tail: many small overlapping balls, so it looks like one smooth tube. It leaves the bottom,
	// sweeps back and up, and curls in at the tip like a question mark. Thick at the base, thin at the tip.
	const n = 40
	for i := 0; i < n; i++ {
		t := float64(i) / (n - 1)
		a := 1.03 - t*4.3         // angle around the curl centre: down, back, up, then forward
		r := 0.58 * (1 - 0.6*t)   // spiral: the curl gets tighter towards the tip
		y := 0.8 - r*math.Cos(a)  // curl centre 0.8 up ...
		z := -0.9 + r*math.Sin(a) // ... and 0.9 behind the monkey
		w := float32(0.085 - 0.04*t)
		ps = append(ps, part{v(float32(0.1*t), float32(y), float32(z)), v(w, w, w), colFur})
	}
	// Lighter fur patch on the back of the head, so it isn't a plain brown ball from behind.
	ps = append(ps, part{v(0, 1.48, -0.2), v(0.26, 0.24, 0.18), rl.NewColor(150, 95, 55, 255)})
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
