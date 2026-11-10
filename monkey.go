package main

import (
	"math"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	colFur  = rl.NewColor(120, 72, 38, 255)
	colSkin = rl.NewColor(232, 190, 140, 255)
	colEye  = rl.NewColor(25, 20, 20, 255)
)

const (
	numMonkeys   = 5
	monkeySpeed  = 2.5 // metres per second
	maxBananas   = 8
	stealRadius  = 1.7
	bananaRadius = 1.4
	rideTime     = 1.2 // seconds a monkey sits on the truck after stealing
	monkeyRest   = 3.0 // seconds after jumping off before it can steal again
	spinTime     = 1.0 // seconds the truck spins after a banana
)

// Monkey walks the roads, steals coins and drops bananas.
type Monkey struct {
	Pos        rl.Vector3
	Heading    float32 // radians, 0 faces +Z (like the car)
	Phase      float32 // walk animation
	Prev, Next Cell    // walking from Prev to Next
	Wait       float32 // seconds standing still before walking
	Ride       float32 // seconds left sitting on the truck
	Rest       float32 // seconds left before it can steal again
	Drop       float32 // seconds until the next banana
}

// Banana lies on the road. Heading is the way the monkey walked when it dropped it.
type Banana struct {
	Pos     rl.Vector3
	Heading float32
}

func cellAt(p rl.Vector3) Cell {
	return Cell{int(math.Floor(float64(p.X / TileSize))), int(math.Floor(float64(p.Z / TileSize)))}
}

func headingTo(from, to rl.Vector3) float32 {
	return float32(math.Atan2(float64(to.X-from.X), float64(to.Z-from.Z)))
}

// pickNext chooses a random road cell next to m.Next to walk to, not going back unless it must.
func (m *Monkey) pickNext(rng *rand.Rand) {
	var opts []Cell
	for _, d := range []int{DirN, DirE, DirS, DirW} {
		dx, dz := dirVec(d)
		c := Cell{m.Next.X + dx, m.Next.Z + dz}
		if isRoadCell(c.X, c.Z) && c != m.Prev {
			opts = append(opts, c)
		}
	}
	m.Prev = m.Next
	if len(opts) == 0 {
		m.Next = m.Prev
		return
	}
	m.Next = opts[rng.Intn(len(opts))]
}

// spawnMonkeys puts one monkey two tiles in front of the truck, looking at it, and the rest on
// random road cells away from the start.
func (g *Game) spawnMonkeys() {
	g.Bananas = g.Bananas[:0]
	start := cellAt(g.Car.Pos)
	near := Cell{start.X, start.Z + 2}
	g.Monkeys = []Monkey{{Pos: CellCenter(near), Heading: math.Pi, Prev: near, Next: near, Wait: 4,
		Drop: 6 + 4*g.Rng.Float32()}}
	for len(g.Monkeys) < numMonkeys {
		c := g.City.Roads[g.Rng.Intn(len(g.City.Roads))]
		if abs(c.X-start.X)+abs(c.Z-start.Z) < 6 {
			continue
		}
		m := Monkey{Pos: CellCenter(c), Prev: c, Next: c, Drop: 3 + 6*g.Rng.Float32(), Phase: g.Rng.Float32() * 6}
		m.pickNext(g.Rng)
		g.Monkeys = append(g.Monkeys, m)
	}
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func (g *Game) updateMonkeys(dt float32) {
	for i := range g.Monkeys {
		m := &g.Monkeys[i]
		if m.Ride > 0 {
			// Sitting on the roof, behind the parcel.
			m.Ride -= dt
			f := g.Car.Forward()
			m.Pos = rl.Vector3Add(g.Car.Pos, rl.NewVector3(-f.X*0.7, 1.45, -f.Z*0.7))
			m.Heading = g.Car.Heading
			if m.Ride <= 0 {
				// Jump off and run away down the road.
				m.Pos = rl.NewVector3(g.Car.Pos.X, 0, g.Car.Pos.Z)
				m.Prev, m.Next = cellAt(m.Pos), cellAt(m.Pos)
				m.pickNext(g.Rng)
				m.Rest = monkeyRest
				g.A.Play("land")
			}
			continue
		}
		m.Rest -= dt
		if m.Wait > 0 {
			m.Wait -= dt
		} else {
			to := CellCenter(m.Next)
			d := rl.Vector3Subtract(to, m.Pos)
			d.Y = 0
			dist := rl.Vector3Length(d)
			step := monkeySpeed * dt
			if dist <= step {
				m.Pos = to
				m.pickNext(g.Rng)
			} else {
				m.Pos = rl.Vector3Add(m.Pos, rl.Vector3Scale(d, step/dist))
				m.Heading = headingTo(m.Pos, to)
			}
			m.Phase += 9 * dt
		}
		m.Pos.Y = 0
		if m.Rest > monkeyRest-0.4 {
			// Little hop when landing from the truck.
			m.Pos.Y = 1.2 * (m.Rest - (monkeyRest - 0.4))
		}

		if m.Rest <= 0 && rl.Vector3Distance(m.Pos, rl.NewVector3(g.Car.Pos.X, 0, g.Car.Pos.Z)) < stealRadius {
			m.Ride = rideTime
			g.A.Play("jump")
			if g.Money > 0 {
				g.Money--
				g.popup("-1  Affe!", g.Car.Pos, rl.NewColor(255, 80, 70, 255))
			} else {
				g.popup("Affe! Nix zu klauen", g.Car.Pos, rl.NewColor(255, 80, 70, 255))
			}
			continue
		}

		m.Drop -= dt
		if m.Drop <= 0 {
			m.Drop = 7 + 6*g.Rng.Float32()
			far := rl.Vector3Distance(m.Pos, g.Car.Pos) > 5
			if far && len(g.Bananas) < maxBananas {
				g.Bananas = append(g.Bananas, Banana{Pos: rl.NewVector3(m.Pos.X, 0, m.Pos.Z), Heading: m.Heading})
			}
		}
	}

	if g.Car.Spin <= 0 {
		for i, b := range g.Bananas {
			if rl.Vector3Distance(b.Pos, rl.NewVector3(g.Car.Pos.X, 0, g.Car.Pos.Z)) < bananaRadius {
				g.Car.StartSpin()
				g.A.Play("skid")
				g.popup("Ausgerutscht!", g.Car.Pos, rl.Yellow)
				g.Bananas = append(g.Bananas[:i], g.Bananas[i+1:]...)
				break
			}
		}
	}
}

var (
	colBanana    = rl.NewColor(250, 215, 50, 255)
	colBananaEnd = rl.NewColor(90, 60, 25, 255)
)

// Draw draws a banana lying on the road: a curved row of yellow balls with brown tips.
func (b *Banana) Draw(a *Assets, clock float32) {
	rl.DrawCylinder(rl.NewVector3(b.Pos.X, 0.07, b.Pos.Z), 0.9, 0.9, 0.01, 16, rl.NewColor(0, 0, 0, 50))
	place := rl.MatrixMultiply(rl.MatrixRotateY(b.Heading+math.Pi/2), rl.MatrixTranslate(b.Pos.X, 0, b.Pos.Z))
	const n = 14
	for i := 0; i < n; i++ {
		t := float64(i)/(n-1)*2 - 1 // -1 .. 1 along the banana
		x := float32(t * 0.9)
		z := float32(0.35 * (1 - t*t)) // the bend
		w := float32(0.22 * (1 - 0.5*t*t))
		col := colBanana
		if i == 0 || i == n-1 {
			col, w = colBananaEnd, 0.1
		}
		a.Ball(rl.NewVector3(x, 0.22, z), rl.NewVector3(w, w, w), col, place)
	}
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
	if m.Ride <= 0 {
		rl.DrawCylinder(rl.NewVector3(m.Pos.X, 0.07, m.Pos.Z), 0.6, 0.6, 0.01, 16, rl.NewColor(0, 0, 0, 70))
	}
	swing := float32(math.Sin(float64(m.Phase)))
	bob := float32(math.Abs(float64(swing))) * 0.08
	place := rl.MatrixMultiply(rl.MatrixRotateY(m.Heading), rl.MatrixTranslate(m.Pos.X, m.Pos.Y+bob, m.Pos.Z))
	for _, p := range monkeyParts(swing) {
		a.Ball(p.off, p.size, p.col, place)
	}
}
