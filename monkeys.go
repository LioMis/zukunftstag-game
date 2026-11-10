package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	numMonkeys   = 6
	monkeySpeed  = 2.2 // m/s
	throwEvery   = 4.0 // seconds between throws, across all monkeys
	throwTime    = 0.9 // seconds a banana flies
	maxBananas   = 12  // oldest banana disappears beyond this
	bananaRadius = 1.4 // how close the truck must come to slip
	slipTime     = 1.0 // seconds the truck spins
	monkeyScale  = 1.15
)

type Monkey struct {
	Pos     rl.Vector3
	From    Cell // road cell it came from, so it doesn't turn back
	Goal    Cell // road cell it walks to
	Heading float32
	Throw   float32 // seconds since its last throw, for the throwing pose
}

type Banana struct {
	From, To rl.Vector3
	Age      float32 // seconds since thrown
	Spin     float32 // degrees around Y, for variety
}

func (b *Banana) Landed() bool { return b.Age >= throwTime }

// Pos is where the banana is now: on a parabola while flying, on the road after.
func (b *Banana) Pos() rl.Vector3 {
	t := clamp(b.Age/throwTime, 0, 1)
	p := rl.Vector3Lerp(b.From, b.To, t)
	p.Y = b.From.Y*(1-t) + 0.15*t + 4*t*(1-t)
	return p
}

func cellOf(p rl.Vector3) Cell {
	return Cell{int(math.Floor(float64(p.X / TileSize))), int(math.Floor(float64(p.Z / TileSize)))}
}

func (g *Game) roadAt(c Cell) bool {
	return c.X >= 0 && c.Z >= 0 && c.X < GridN && c.Z < GridN && g.City.T[c.X][c.Z].Road
}

func (g *Game) spawnMonkeys() {
	g.Monkeys = g.Monkeys[:0]
	g.Bananas = g.Bananas[:0]
	g.NextThrow = throwEvery
	for len(g.Monkeys) < numMonkeys {
		c := g.City.Roads[g.Rng.Intn(len(g.City.Roads))]
		if rl.Vector3Distance(CellCenter(c), g.Car.Pos) < 3*TileSize {
			continue // not right on top of the truck at the start
		}
		m := Monkey{Pos: CellCenter(c), From: c, Goal: c, Throw: 99}
		g.Monkeys = append(g.Monkeys, m)
	}
}

// nextCell picks a neighbouring road cell, turning back only at a dead end.
func (g *Game) nextCell(m *Monkey) Cell {
	var opts []Cell
	for _, d := range []int{DirN, DirE, DirS, DirW} {
		dx, dz := dirVec(d)
		c := Cell{m.Goal.X + dx, m.Goal.Z + dz}
		if g.roadAt(c) && c != m.From {
			opts = append(opts, c)
		}
	}
	if len(opts) == 0 {
		return m.From
	}
	return opts[g.Rng.Intn(len(opts))]
}

func (g *Game) updateMonkeys(dt float32) {
	for i := range g.Monkeys {
		m := &g.Monkeys[i]
		m.Throw += dt
		if m.Throw < 0.5 {
			continue // stands still while throwing
		}
		goal := CellCenter(m.Goal)
		d := rl.Vector3Subtract(goal, m.Pos)
		dist := rl.Vector3Length(d)
		if dist < 0.1 {
			next := g.nextCell(m)
			m.From, m.Goal = m.Goal, next
			continue
		}
		step := float32(math.Min(float64(monkeySpeed*dt), float64(dist)))
		m.Pos = rl.Vector3Add(m.Pos, rl.Vector3Scale(d, step/dist))
		m.Heading = float32(math.Atan2(float64(d.X), float64(d.Z)))
	}

	g.NextThrow -= dt
	if g.NextThrow <= 0 && len(g.Monkeys) > 0 {
		g.NextThrow = throwEvery
		g.throwBanana(g.nearestMonkey())
	}

	for i := range g.Bananas {
		g.Bananas[i].Age += dt
	}
	// Slip: the truck runs over a banana on the road.
	for i := 0; i < len(g.Bananas); i++ {
		b := g.Bananas[i]
		p := b.Pos()
		if !b.Landed() || rl.Vector3Distance(rl.NewVector3(p.X, 0, p.Z), rl.NewVector3(g.Car.Pos.X, 0, g.Car.Pos.Z)) > bananaRadius {
			continue
		}
		g.Bananas = append(g.Bananas[:i], g.Bananas[i+1:]...)
		i--
		if g.Car.Slip <= 0 {
			g.Car.StartSlip()
			g.A.Play("skid")
			g.popup("Ausgerutscht!", g.Car.Pos, rl.NewColor(255, 230, 60, 255))
		}
	}
}

// nearestMonkey is the monkey closest to the truck, so throws happen where the player can see them.
func (g *Game) nearestMonkey() *Monkey {
	best := &g.Monkeys[0]
	for i := range g.Monkeys {
		if rl.Vector3Distance(g.Monkeys[i].Pos, g.Car.Pos) < rl.Vector3Distance(best.Pos, g.Car.Pos) {
			best = &g.Monkeys[i]
		}
	}
	return best
}

// throwBanana throws a banana a few metres ahead of the monkey onto its road.
func (g *Game) throwBanana(m *Monkey) {
	ahead := rl.Vector3Subtract(CellCenter(m.Goal), m.Pos)
	if rl.Vector3Length(ahead) < 0.1 {
		ahead = rl.NewVector3(0, 0, 1)
	}
	ahead = rl.Vector3Normalize(ahead)
	side := rl.NewVector3(ahead.Z, 0, -ahead.X)
	to := rl.Vector3Add(m.Pos, rl.Vector3Scale(ahead, 4+g.Rng.Float32()*4))
	to = rl.Vector3Add(to, rl.Vector3Scale(side, (g.Rng.Float32()-0.5)*2))
	if !g.roadAt(cellOf(to)) {
		to = CellCenter(m.Goal) // keep it on the road
	}
	from := rl.Vector3Add(m.Pos, rl.NewVector3(0, 1.6*monkeyScale, 0))
	g.Bananas = append(g.Bananas, Banana{From: from, To: to, Spin: g.Rng.Float32() * 360})
	if len(g.Bananas) > maxBananas {
		g.Bananas = g.Bananas[1:]
	}
	m.Throw = 0
	m.Heading = float32(math.Atan2(float64(ahead.X), float64(ahead.Z)))
}

// warpToBanana parks the truck a few metres up the road from a banana on the road, facing it.
// Script only, so tests don't have to steer.
func (g *Game) warpToBanana() {
	for i := len(g.Bananas) - 1; i >= 0; i-- {
		b := g.Bananas[i]
		if !b.Landed() {
			continue
		}
		p := b.Pos()
		p.Y = 0
		for _, d := range []int{DirN, DirE, DirS, DirW} {
			dx, dz := dirVec(d)
			start := rl.NewVector3(p.X+float32(dx)*7, 0, p.Z+float32(dz)*7)
			if !g.roadAt(cellOf(start)) || !g.roadAt(cellOf(rl.NewVector3(p.X+float32(dx)*3.5, 0, p.Z+float32(dz)*3.5))) {
				continue
			}
			g.parkCar(start, rl.Vector3Subtract(p, start))
			return
		}
	}
}

func (g *Game) drawMonkeys() {
	for _, m := range g.Monkeys {
		p := m.Pos
		walking := m.Throw >= 0.5
		if walking {
			p.Y = 0.08 * float32(math.Abs(math.Sin(float64(g.Clock)*9+float64(m.Pos.X))))
		}
		rl.DrawCylinder(rl.NewVector3(m.Pos.X, 0.07, m.Pos.Z), 0.5, 0.5, 0.01, 16, rl.NewColor(0, 0, 0, 60))
		g.A.Draw("monkey", p, m.Heading*rl.Rad2deg, monkeyScale)
	}
	for _, b := range g.Bananas {
		rot := b.Spin
		if !b.Landed() {
			rot += b.Age * 720
		}
		drawBanana(b.Pos(), rot)
	}
}

// drawBanana draws a curved banana from a few spheres, about 0.9 m long so it reads from the chase camera.
func drawBanana(p rl.Vector3, rotDeg float32) {
	yellow := rl.NewColor(250, 215, 40, 255)
	brown := rl.NewColor(90, 60, 20, 255)
	a := float64(rotDeg * rl.Deg2rad)
	dir := rl.NewVector3(float32(math.Sin(a)), 0, float32(math.Cos(a)))
	side := rl.NewVector3(dir.Z, 0, -dir.X)
	const n = 7
	for i := 0; i < n; i++ {
		t := float32(i)/(n-1)*2 - 1 // -1..1 along the banana
		bend := 0.25 * (1 - t*t)    // curve to one side
		q := rl.Vector3Add(p, rl.Vector3Add(rl.Vector3Scale(dir, t*0.45), rl.Vector3Scale(side, bend)))
		r := 0.16 * float32(math.Sqrt(float64(1-0.6*t*t)))
		col := yellow
		if i == 0 || i == n-1 {
			col = brown
			r = 0.07
		}
		rl.DrawSphere(q, r, col)
	}
}
