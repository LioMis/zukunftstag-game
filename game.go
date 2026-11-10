package main

import (
	"fmt"
	"math"
	"math/rand"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type State int

const (
	StateTitle State = iota
	StatePlay
	StateOver
	StateScores
)

const (
	roundTime     = 75.0
	deliveryBonus = 12.0
	numCoins      = 30
	triggerRadius = 2.4
)

type Popup struct {
	Text string
	Pos  rl.Vector3
	Age  float32
	Col  rl.Color
}

type Game struct {
	A     *Assets
	Rng   *rand.Rand
	State State
	City  *City
	Car   Car
	Cam   rl.Camera3D
	Time  float32 // seconds left
	Clock float32 // seconds since start of state
	Money int
	Done  int // deliveries

	HasParcel bool
	Pickup    int // index into City.Buildings
	Dropoff   int
	Coins     []rl.Vector3
	Monkeys   []Monkey
	Popups    []Popup

	Scores     []Score
	ScoresPath string
	Name       string
	Rank       int
	Qualified  bool
}

func NewGame(a *Assets, rng *rand.Rand, scoresPath string) *Game {
	g := &Game{A: a, Rng: rng, ScoresPath: scoresPath, Rank: -1}
	g.Scores = LoadScores(scoresPath)
	g.reset()
	g.State = StateTitle
	return g
}

func (g *Game) reset() {
	g.City = NewCity(g.Rng)
	trucks := []string{"vehicle-truck-yellow", "vehicle-truck-red", "vehicle-truck-green", "vehicle-truck-purple"}
	start := CellCenter(Cell{GridN / 2, GridN / 2})
	g.Car = Car{Pos: start, Model: trucks[g.Rng.Intn(len(trucks))]}
	g.Time = roundTime
	g.Money, g.Done = 0, 0
	g.HasParcel = false
	g.Popups = nil
	g.newJob()
	g.Coins = g.Coins[:0]
	for len(g.Coins) < numCoins {
		g.Coins = append(g.Coins, g.randomRoadSpot())
	}
	// One monkey two tiles in front of the truck, looking at it.
	g.Monkeys = []Monkey{{Pos: rl.NewVector3(start.X, 0, start.Z+2*TileSize), Heading: math.Pi}}
	g.Cam = rl.Camera3D{Position: rl.NewVector3(start.X, 8, start.Z-10), Target: start,
		Up: rl.NewVector3(0, 1, 0), Fovy: 55, Projection: rl.CameraPerspective}
}

func (g *Game) randomRoadSpot() rl.Vector3 {
	c := g.City.Roads[g.Rng.Intn(len(g.City.Roads))]
	p := CellCenter(c)
	p.Y = 0.6
	return p
}

// newJob picks a pickup that is not too close to the car and a drop-off far from the pickup.
func (g *Game) newJob() {
	b := g.City.Buildings
	for i := 0; i < 50; i++ {
		g.Pickup = g.Rng.Intn(len(b))
		if rl.Vector3Distance(b[g.Pickup].Front, g.Car.Pos) > 3*TileSize {
			break
		}
	}
	for i := 0; i < 50; i++ {
		g.Dropoff = g.Rng.Intn(len(b))
		d := rl.Vector3Distance(b[g.Dropoff].Front, b[g.Pickup].Front)
		if d > 4*TileSize && d < 12*TileSize {
			break
		}
	}
}

func (g *Game) Target() rl.Vector3 {
	if g.HasParcel {
		return g.City.Buildings[g.Dropoff].Front
	}
	return g.City.Buildings[g.Pickup].Front
}

func (g *Game) popup(text string, pos rl.Vector3, col rl.Color) {
	g.Popups = append(g.Popups, Popup{Text: text, Pos: rl.NewVector3(pos.X, pos.Y+2, pos.Z), Col: col})
}

func (g *Game) Update(in Input, dt float32) {
	g.Clock += dt
	switch g.State {
	case StateTitle:
		g.idleCamera(dt)
		if in.Enter {
			g.reset()
			g.State, g.Clock = StatePlay, 0
			g.A.Play("toggle")
		}
	case StatePlay:
		g.updatePlay(in, dt)
	case StateOver:
		g.idleCamera(dt)
		if g.Qualified {
			for _, r := range in.Chars {
				if len([]rune(g.Name)) < 12 && (r == ' ' || r == '-' || isNameRune(r)) {
					g.Name += string(r)
				}
			}
			if in.Back && len(g.Name) > 0 {
				rs := []rune(g.Name)
				g.Name = string(rs[:len(rs)-1])
			}
			if in.Enter && g.Clock > 0.5 && strings.TrimSpace(g.Name) != "" {
				g.Scores, g.Rank = Insert(g.Scores, Score{Name: strings.TrimSpace(g.Name), Score: g.Money})
				if err := SaveScores(g.ScoresPath, g.Scores); err != nil {
					fmt.Println("saving scores:", err)
				}
				g.State, g.Clock = StateScores, 0
			}
		} else if in.Enter && g.Clock > 1 {
			g.State, g.Clock = StateScores, 0
		}
	case StateScores:
		g.idleCamera(dt)
		if in.Enter && g.Clock > 0.5 {
			g.State, g.Clock = StateTitle, 0
		}
	}
	for i := range g.Popups {
		g.Popups[i].Age += dt
		g.Popups[i].Pos.Y += 1.5 * dt
	}
	for len(g.Popups) > 0 && g.Popups[0].Age > 1.5 {
		g.Popups = g.Popups[1:]
	}
}

func isNameRune(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || strings.ContainsRune("äöüÄÖÜéèàç", r)
}

func (g *Game) updatePlay(in Input, dt float32) {
	g.Time -= dt
	if g.Time <= 0 {
		g.Time = 0
		g.Qualified = Qualifies(g.Scores, g.Money)
		g.Name, g.Rank = "", -1
		g.State, g.Clock = StateOver, 0
		g.A.Play("rotate")
		return
	}
	if in.Warp != "" {
		b := g.City.Buildings[g.Pickup]
		if in.Warp == "dropoff" {
			b = g.City.Buildings[g.Dropoff]
		}
		// Park a few metres up the road from the marker, facing it.
		dx := b.Front.X - CellCenter(b.Cell).X
		dz := b.Front.Z - CellCenter(b.Cell).Z
		side := rl.Vector3Normalize(rl.NewVector3(dz, 0, -dx))
		g.Car.Pos = rl.Vector3Add(b.Front, rl.Vector3Scale(side, 7))
		g.Car.Heading = float32(math.Atan2(float64(-side.X), float64(-side.Z)))
		g.Car.Speed = 0
		f := g.Car.Forward()
		g.Cam.Position = rl.Vector3Add(g.Car.Pos, rl.NewVector3(-f.X*7, 13, -f.Z*7))
		g.Cam.Target = g.Car.Pos
	}
	if hit := g.Car.Update(in, g.City, dt); hit > 0 {
		g.A.Play("impact")
	}

	if rl.Vector3Distance(g.Car.Pos, g.Target()) < triggerRadius {
		if !g.HasParcel {
			g.HasParcel = true
			g.A.Play("placement-a")
			g.popup("Paket geladen!", g.Car.Pos, rl.Yellow)
		} else {
			b := g.City.Buildings
			dist := rl.Vector3Distance(b[g.Pickup].Front, b[g.Dropoff].Front)
			pay := 10 + int(dist/TileSize)*2
			g.Money += pay
			g.Done++
			g.Time += deliveryBonus
			g.HasParcel = false
			g.A.Play("toggle")
			g.popup(fmt.Sprintf("+%d Fr.  +%ds", pay, int(deliveryBonus)), g.Car.Pos, rl.Green)
			g.newJob()
		}
	}
	for i := range g.Coins {
		if rl.Vector3Distance(g.Car.Pos, rl.NewVector3(g.Coins[i].X, 0, g.Coins[i].Z)) < 1.8 {
			g.Money++
			g.A.Play("coin")
			g.popup("+1", g.Coins[i], rl.Gold)
			g.Coins[i] = g.randomRoadSpot()
		}
	}

	// Chase camera.
	f := g.Car.Forward()
	want := rl.Vector3Add(g.Car.Pos, rl.NewVector3(-f.X*7, 13, -f.Z*7))
	k := 1 - float32(math.Exp(float64(-4*dt)))
	g.Cam.Position = rl.Vector3Lerp(g.Cam.Position, want, k)
	g.Cam.Target = rl.Vector3Lerp(g.Cam.Target, rl.Vector3Add(g.Car.Pos, rl.NewVector3(f.X*5, 0, f.Z*5)), k*2)
}

// idleCamera slowly circles the town for menus.
func (g *Game) idleCamera(dt float32) {
	mid := float32(GridN) * TileSize / 2
	a := float64(g.Clock) * 0.08
	g.Cam.Position = rl.NewVector3(mid+float32(math.Cos(a))*48, 48, mid+float32(math.Sin(a))*48)
	g.Cam.Target = rl.NewVector3(mid, 0, mid)
}
