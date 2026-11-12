package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	carRadius   = 1.0
	carMaxSpeed = 17.0
	carMaxRev   = 6.0
)

type Car struct {
	Pos     rl.Vector3
	Heading float32 // radians, 0 faces +Z
	Speed   float32
	Model   string
	Bump    float32 // seconds since the last crash, for sound throttling
}

func (c *Car) Forward() rl.Vector3 {
	return rl.NewVector3(float32(math.Sin(float64(c.Heading))), 0, float32(math.Cos(float64(c.Heading))))
}

// Update moves the car. Returns the impact speed if it hit something this frame.
func (c *Car) Update(in Input, city *City, dt float32) float32 {
	switch {
	case in.Up:
		if c.Speed < 0 {
			c.Speed += 30 * dt
		} else {
			c.Speed += 16 * dt
		}
	case in.Down:
		if c.Speed > 0 {
			c.Speed -= 30 * dt
		} else {
			c.Speed -= 10 * dt
		}
	default:
		c.Speed -= c.Speed * 1.5 * dt
	}
	c.Speed = clamp(c.Speed, -carMaxRev, carMaxSpeed)

	steer := float32(0)
	if in.Left {
		steer++
	}
	if in.Right {
		steer--
	}
	// Standing still still turns (on the spot), just a bit slower than at speed.
	grip := 0.6 + 0.4*clamp(float32(math.Abs(float64(c.Speed)))/5, 0, 1)
	sign := float32(1)
	if c.Speed < 0 {
		sign = -1
	}
	c.Heading += steer * 2.4 * grip * sign * dt

	c.Pos = rl.Vector3Add(c.Pos, rl.Vector3Scale(c.Forward(), c.Speed*dt))
	pos, n := city.Collide(c.Pos, carRadius)
	c.Pos = pos
	c.Bump += dt
	if n.X != 0 || n.Z != 0 {
		// Lose the part of the velocity going into the wall.
		into := -rl.Vector3DotProduct(c.Forward(), n) * c.Speed
		if into > 0 {
			c.Speed *= 0.6
			if into > 6 && c.Bump > 0.4 {
				c.Bump = 0
				return into
			}
		}
	}
	return 0
}

func (c *Car) Draw(a *Assets) {
	// Blob shadow.
	rl.DrawCylinder(rl.NewVector3(c.Pos.X, 0.07, c.Pos.Z), 1.3, 1.3, 0.01, 20, rl.NewColor(0, 0, 0, 70))
	a.Draw(c.Model, c.Pos, c.Heading*rl.Rad2deg, 1)
}
