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
// dt is the time since the last frame in seconds (1/60 at 60 fps), so "x * dt" means "x per second".
func (c *Car) Update(in Input, city *City, dt float32) float32 {
	switch { // pick the first case that is true
	case in.Up: // gas pedal pressed
		if c.Speed < 0 { // still rolling backwards?
			c.Speed += 30 * dt // then brake hard: 30 speed units per second
		} else { // standing or already going forward
			c.Speed += 16 * dt // speed up more gently: 16 per second
		}
	case in.Down: // brake / reverse pressed
		if c.Speed > 0 { // still rolling forwards?
			c.Speed -= 30 * dt // then brake hard
		} else { // standing or already going backwards
			c.Speed -= 10 * dt // reverse slowly: 10 per second
		}
	default: // no pedal pressed
		c.Speed -= c.Speed * 1.5 * dt // roll out: lose 1.5x the current speed per second, so fast slows quickly and slow creeps
	}
	c.Speed = clamp(c.Speed, -carMaxRev, carMaxSpeed) // never faster than 17 forwards or 6 backwards

	steer := float32(0) // 0 = wheel straight
	if in.Left {        // left pressed
		steer++ // +1 = turn left
	}
	if in.Right { // right pressed
		steer-- // -1 = turn right; both pressed cancels out to 0
	}
	grip := clamp(float32(math.Abs(float64(c.Speed)))/5, 0, 1) // standing car can't turn; full turning from speed 5 up
	sign := float32(1)                                         // driving forwards: steering works normally
	if c.Speed < 0 {                                           // driving backwards?
		sign = -1 // then flip steering, like a real car reversing
	}
	c.Heading += steer * 2.4 * grip * sign * dt // turn the truck: at most 2.4 radians (about 140°) per second

	c.Pos = rl.Vector3Add(c.Pos, rl.Vector3Scale(c.Forward(), c.Speed*dt)) // move along the way the truck faces, by speed × time
	pos, n := city.Collide(c.Pos, carRadius)                               // ask the town: are we inside a building or tree? get a fixed position and the wall direction n
	c.Pos = pos                                                            // put the truck back out of the wall
	c.Bump += dt                                                           // count time since the last crash sound
	if n.X != 0 || n.Z != 0 {                                              // n is not zero: we touched a wall
		// Lose the part of the velocity going into the wall.
		into := -rl.Vector3DotProduct(c.Forward(), n) * c.Speed // how fast we drove straight into the wall (0 if we only scraped along it)
		if into > 0 {                                           // really moving towards the wall, not away from it
			c.Speed *= 0.6                // crash: keep only 60% of the speed
			if into > 6 && c.Bump > 0.4 { // hard enough, and the last crash sound was over 0.4 s ago?
				c.Bump = 0  // restart the crash-sound timer
				return into // tell the game how hard we hit, so it plays the crash sound
			}
		}
	}
	return 0 // no loud crash this frame
}

func (c *Car) Draw(a *Assets) {
	// Blob shadow.
	rl.DrawCylinder(rl.NewVector3(c.Pos.X, 0.07, c.Pos.Z), 1.3, 1.3, 0.01, 20, rl.NewColor(0, 0, 0, 70))
	a.Draw(c.Model, c.Pos, c.Heading*rl.Rad2deg, 1)
}
