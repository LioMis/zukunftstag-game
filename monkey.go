package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// monkeySpeed is how fast the monkey runs, in metres per second. The truck does up to 17.
const monkeySpeed = 4.0

type Monkey struct {
	Pos rl.Vector3
	Rot float32 // which way it looks, in degrees
}

// Update moves the monkey a little towards the truck. Returns true if it caught the truck.
func (m *Monkey) Update(truck rl.Vector3, city *City, dt float32) bool {
	toTruck := rl.Vector3Subtract(truck, m.Pos)
	toTruck.Y = 0
	if rl.Vector3Length(toTruck) < 2 {
		return true
	}
	dir := rl.Vector3Normalize(toTruck)
	m.Pos = rl.Vector3Add(m.Pos, rl.Vector3Scale(dir, monkeySpeed*dt))
	m.Pos, _ = city.Collide(m.Pos, 0.5)
	m.Rot = float32(math.Atan2(float64(dir.X), float64(dir.Z))) * rl.Rad2deg
	return false
}

func (m *Monkey) Draw(a *Assets) {
	a.Draw("monkey", m.Pos, m.Rot, 1)
}

// updateMonkey runs the monkey. When it catches the truck it steals a coin and jumps far away
// to start the chase again.
func (g *Game) updateMonkey(dt float32) {
	if !g.Monkey.Update(g.Car.Pos, g.City, dt) {
		return
	}
	if g.Money > 0 {
		g.Money--
	}
	g.A.Play("impact")
	g.popup("Affe klaut 1 Fr.!", g.Car.Pos, rl.Orange)
	g.Monkey.Pos = g.farRoadSpot()
}

// farRoadSpot picks a road spot far from the truck, so the monkey starts a long way off.
func (g *Game) farRoadSpot() rl.Vector3 {
	p := g.randomRoadSpot()
	for i := 0; i < 50 && rl.Vector3Distance(p, g.Car.Pos) < 6*TileSize; i++ {
		p = g.randomRoadSpot()
	}
	p.Y = 0
	return p
}
