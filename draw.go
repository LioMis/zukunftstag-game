package main

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	colPickup  = rl.NewColor(255, 200, 40, 255)
	colDropoff = rl.NewColor(60, 220, 90, 255)
)

func (g *Game) Draw() {
	rl.ClearBackground(fogColor)
	g.A.SetView(g.Cam)
	rl.BeginMode3D(g.Cam)
	mid := float32(GridN) * TileSize / 2
	rl.DrawPlane(rl.NewVector3(mid, -0.01, mid), rl.NewVector2(600, 600), rl.NewColor(120, 190, 90, 255))
	g.City.Draw(g.A)

	t := float64(g.Clock)
	for i, c := range g.Coins {
		p := c
		p.Y = 0.4 + 0.15*float32(math.Sin(t*3+float64(i)))
		g.A.Draw("coin", p, float32(t*120+float64(i)*40), 2.2)
	}
	if g.State == StatePlay {
		b := g.City.Buildings
		if !g.HasParcel {
			g.drawMarker(b[g.Pickup].Front, colPickup)
		}
		g.drawMarker(b[g.Dropoff].Front, colDropoff)
		g.Car.Draw(g.A)
		for i := range g.Bananas {
			g.Bananas[i].Draw(g.A, g.Clock)
		}
		for i := range g.Monkeys {
			g.Monkeys[i].Draw(g.A)
		}
		if g.HasParcel {
			p := rl.Vector3Add(g.Car.Pos, rl.NewVector3(0, 1.9, 0))
			rl.DrawCube(p, 0.8, 0.6, 0.8, rl.NewColor(196, 140, 80, 255))
			rl.DrawCubeWires(p, 0.8, 0.6, 0.8, rl.NewColor(110, 70, 30, 255))
		}
		g.drawArrow()
	}
	rl.EndMode3D()

	for _, p := range g.Popups {
		s := rl.GetWorldToScreen(p.Pos, g.Cam)
		alpha := uint8(255 * clamp(1.5-p.Age, 0, 1))
		col := p.Col
		col.A = alpha
		g.textCentered(p.Text, s.X, s.Y, 40, col)
	}

	switch g.State {
	case StateTitle:
		g.drawTitle()
	case StatePlay:
		g.drawHUD()
	case StateOver:
		g.drawOver()
	case StateScores:
		g.drawScores()
	}
}

func (g *Game) drawMarker(p rl.Vector3, col rl.Color) {
	pulse := float32(0.5 + 0.5*math.Sin(float64(g.Clock)*5))
	ring := col
	ring.A = 160
	rl.DrawCylinder(rl.NewVector3(p.X, 0.08, p.Z), triggerRadius, triggerRadius, 0.04, 32, ring)
	beam := col
	beam.A = uint8(110 + 70*pulse)
	rl.DrawCylinder(rl.NewVector3(p.X, 0.1, p.Z), 0.8, 0.8, 40, 16, beam)
}

// drawArrow draws a flat arrow hovering above the car, pointing at the target.
func (g *Game) drawArrow() {
	d := rl.Vector3Subtract(g.Target(), g.Car.Pos)
	d.Y = 0
	if rl.Vector3Length(d) < 3 {
		return
	}
	d = rl.Vector3Normalize(d)
	side := rl.NewVector3(d.Z, 0, -d.X)
	c := rl.Vector3Add(g.Car.Pos, rl.NewVector3(0, 3.2, 0))
	tip := rl.Vector3Add(c, rl.Vector3Scale(d, 1.6))
	l := rl.Vector3Add(c, rl.Vector3Scale(side, 0.9))
	r := rl.Vector3Add(c, rl.Vector3Scale(side, -0.9))
	col := colPickup
	if g.HasParcel {
		col = colDropoff
	}
	rl.DrawTriangle3D(tip, l, r, col)
	rl.DrawTriangle3D(tip, r, l, col)
	tail := rl.Vector3Subtract(c, rl.Vector3Scale(d, 0.9))
	sl := rl.Vector3Scale(side, 0.35)
	q := []rl.Vector3{rl.Vector3Add(c, sl), rl.Vector3Subtract(c, sl), rl.Vector3Subtract(tail, sl), rl.Vector3Add(tail, sl)}
	rl.DrawTriangle3D(q[0], q[1], q[2], col)
	rl.DrawTriangle3D(q[0], q[2], q[3], col)
	rl.DrawTriangle3D(q[2], q[1], q[0], col)
	rl.DrawTriangle3D(q[3], q[2], q[0], col)
}

func (g *Game) text(s string, x, y, size float32, col rl.Color) {
	sh := rl.NewColor(0, 0, 0, col.A/2)
	rl.DrawTextEx(g.A.Font, s, rl.NewVector2(x+3, y+3), size, 1, sh)
	rl.DrawTextEx(g.A.Font, s, rl.NewVector2(x, y), size, 1, col)
}

func (g *Game) textCentered(s string, cx, y, size float32, col rl.Color) {
	w := rl.MeasureTextEx(g.A.Font, s, size, 1).X
	g.text(s, cx-w/2, y, size, col)
}

func (g *Game) textRight(s string, rx, y, size float32, col rl.Color) {
	w := rl.MeasureTextEx(g.A.Font, s, size, 1).X
	g.text(s, rx-w, y, size, col)
}

func (g *Game) panel(x, y, w, h float32) {
	rl.DrawRectangleRounded(rl.NewRectangle(x, y, w, h), 0.15, 8, rl.NewColor(20, 30, 50, 180))
}

func (g *Game) drawHUD() {
	sw := float32(rl.GetScreenWidth())
	timeCol := rl.White
	if g.Time < 10 && int(g.Time*4)%2 == 0 {
		timeCol = rl.NewColor(255, 90, 80, 255)
	}
	g.panel(20, 20, 220, 70)
	g.text(fmt.Sprintf("Zeit %d", int(math.Ceil(float64(g.Time)))), 40, 28, 52, timeCol)
	g.panel(sw-260, 20, 240, 70)
	g.textRight(fmt.Sprintf("%d Fr.", g.Money), sw-40, 28, 52, rl.Gold)

	b := g.City.Buildings
	var task string
	col := colPickup
	if g.HasParcel {
		task = "Liefern an: " + b[g.Dropoff].Name
		col = colDropoff
	} else {
		task = "Abholen bei: " + b[g.Pickup].Name
	}
	dist := rl.Vector3Distance(g.Car.Pos, g.Target())
	task += fmt.Sprintf("  (%d m)", int(dist*5))
	w := rl.MeasureTextEx(g.A.Font, task, 40, 1).X + 50
	g.panel(sw/2-w/2, 20, w, 60)
	g.textCentered(task, sw/2, 30, 40, col)
}

func (g *Game) drawTitle() {
	sw, sh := float32(rl.GetScreenWidth()), float32(rl.GetScreenHeight())
	rl.DrawRectangle(0, 0, int32(sw), int32(sh), rl.NewColor(10, 20, 40, 90))
	bob := float32(math.Sin(float64(g.Clock)*2)) * 6
	g.textCentered("PAKET-BLITZ", sw/2, sh*0.18+bob, 140, rl.NewColor(255, 210, 60, 255))
	g.textCentered("Hol Pakete ab, liefere sie aus, sammle Münzen!", sw/2, sh*0.18+150, 44, rl.White)
	g.textCentered("Pfeiltasten / WASD: fahren", sw/2, sh*0.55, 40, rl.White)
	g.textCentered("Achtung: Affen klauen Münzen und werfen Bananen!", sw/2, sh*0.55+50, 36, rl.NewColor(255, 200, 120, 255))
	if int(g.Clock*2)%2 == 0 {
		g.textCentered("ENTER drücken zum Starten", sw/2, sh*0.72, 52, rl.NewColor(120, 255, 140, 255))
	}
	if len(g.Scores) > 0 {
		g.textCentered(fmt.Sprintf("Rekord: %s – %d Fr.", g.Scores[0].Name, g.Scores[0].Score), sw/2, sh*0.86, 36, rl.Gold)
	}
}

func (g *Game) drawOver() {
	sw, sh := float32(rl.GetScreenWidth()), float32(rl.GetScreenHeight())
	rl.DrawRectangle(0, 0, int32(sw), int32(sh), rl.NewColor(10, 20, 40, 140))
	g.textCentered("ZEIT UM!", sw/2, sh*0.15, 120, rl.NewColor(255, 210, 60, 255))
	pakete := "Pakete"
	if g.Done == 1 {
		pakete = "Paket"
	}
	g.textCentered(fmt.Sprintf("%d %s geliefert – %d Franken", g.Done, pakete, g.Money), sw/2, sh*0.15+135, 52, rl.White)
	if g.Qualified {
		g.textCentered("Neuer Highscore! Wie heisst du?", sw/2, sh*0.52, 48, rl.NewColor(120, 255, 140, 255))
		cursor := ""
		if int(g.Clock*2)%2 == 0 {
			cursor = "_"
		}
		g.panel(sw/2-300, sh*0.62, 600, 90)
		g.textCentered(g.Name+cursor, sw/2, sh*0.62+12, 64, rl.White)
		g.textCentered("ENTER zum Speichern", sw/2, sh*0.82, 36, rl.LightGray)
	} else {
		g.textCentered("ENTER zum Weiter", sw/2, sh*0.7, 44, rl.LightGray)
	}
}

func (g *Game) drawScores() {
	sw, sh := float32(rl.GetScreenWidth()), float32(rl.GetScreenHeight())
	rl.DrawRectangle(0, 0, int32(sw), int32(sh), rl.NewColor(10, 20, 40, 140))
	g.textCentered("BESTENLISTE", sw/2, sh*0.06, 90, rl.NewColor(255, 210, 60, 255))
	g.panel(sw/2-360, sh*0.22, 720, sh*0.62)
	for i, s := range g.Scores {
		y := sh*0.24 + float32(i)*sh*0.058
		col := rl.White
		if i == g.Rank {
			col = rl.NewColor(120, 255, 140, 255)
		}
		g.text(fmt.Sprintf("%2d.", i+1), sw/2-320, y, 40, col)
		g.text(s.Name, sw/2-230, y, 40, col)
		g.textRight(fmt.Sprintf("%d Fr.", s.Score), sw/2+320, y, 40, col)
	}
	g.textCentered("ENTER: nochmal spielen", sw/2, sh*0.89, 40, rl.LightGray)
}
