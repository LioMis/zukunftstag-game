// This file is the truck: what it knows about itself, how it drives, and how it is drawn.
//
// Every Go file starts by saying which package it belongs to. All files of this game are in
// package "main", so they can use each other's things without any extra work.
package main

// import lists the other code this file uses: "math" for sine, cosine and absolute values, and
// raylib (short name "rl") for 3D points and drawing.
import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Constants are fixed numbers with a name. If we want a faster truck, we change the number here
// once, instead of hunting for 17 everywhere in the code.
const (
	carRadius   = 1.0  // the truck counts as a circle this big when it bumps into walls
	carMaxSpeed = 17.0 // top speed forwards
	carMaxRev   = 6.0  // top speed backwards
)

// A struct is a bundle of values that belong together, like a form with fields to fill in.
// A Car knows where it is, which way it faces, how fast it goes, which truck model it looks like,
// and how long ago it last crashed.
type Car struct {
	Pos     rl.Vector3 // position: three numbers X (left/right), Y (up), Z (forward/back)
	Heading float32    // which way it faces, as an angle in radians; 0 faces +Z
	Speed   float32    // positive = forwards, negative = backwards
	Model   string     // name of the 3D model, e.g. "vehicle-truck-red"
	Bump    float32    // seconds since the last crash sound, so it doesn't play 60 times a second
}

// A method is a function that belongs to a struct. "(c *Car)" means: this works on one car,
// called c inside the method. c.Speed is that car's speed.
//
// Forward turns the heading angle into an arrow of length 1 pointing where the truck faces.
// Sine and cosine are the maths that turn an angle into "how much sideways, how much forwards".
func (c *Car) Forward() rl.Vector3 {
	return rl.NewVector3(float32(math.Sin(float64(c.Heading))), 0, float32(math.Cos(float64(c.Heading))))
}

// Update is called once per frame, 60 times a second. It gets the keys pressed this frame (in),
// the town (city), and the time since the last frame (dt, about 1/60 of a second). It moves the
// truck a tiny bit, and returns how hard it crashed, or 0 if it didn't.
//
// Why "* dt" everywhere: "Speed += 16 * dt" adds 16 per second, split into 60 small steps. The
// truck then drives equally fast on a slow and on a fast computer.
func (c *Car) Update(in Input, city *City, dt float32) float32 {
	// Step 1: gas and brake.
	// A switch checks its cases from top to bottom and runs only the first one that is true.
	// Gas while rolling backwards brakes hard first; gas while standing or going forwards
	// speeds up. The brake key does the same the other way round, so holding it long enough
	// drives backwards. With no key pressed the truck rolls out: it loses a share of its speed
	// every frame, so it slows down fast at first and gently at the end.
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
	// clamp keeps a number between a lowest and a highest value: here, the speed limits.
	c.Speed = clamp(c.Speed, -carMaxRev, carMaxSpeed)

	// Step 2: steering.
	// steer becomes +1 for left, -1 for right, and 0 for none or both keys.
	// grip makes steering depend on speed, like a real car: standing still you can't turn,
	// and from speed 5 up you turn fully. sign flips the steering when reversing, because a
	// car going backwards turns the other way. Then the heading changes by all of that together.
	steer := float32(0)
	if in.Left {
		steer++
	}
	if in.Right {
		steer--
	}
	grip := clamp(float32(math.Abs(float64(c.Speed)))/5, 0, 1)
	sign := float32(1)
	if c.Speed < 0 {
		sign = -1
	}
	c.Heading += steer * 2.4 * grip * sign * dt

	// Step 3: move, then bump into walls.
	// First we move the truck along its Forward arrow by speed × time, without asking whether
	// there is a house in the way. Then we ask the town (city.Collide) to push the truck back out
	// of any house or tree it now overlaps. The town gives back the fixed position, and n, an
	// arrow pointing away from the wall it touched (or zero if it touched nothing).
	// Bump counts up the time since the last crash sound.
	c.Pos = rl.Vector3Add(c.Pos, rl.Vector3Scale(c.Forward(), c.Speed*dt))
	pos, n := city.Collide(c.Pos, carRadius)
	c.Pos = pos
	c.Bump += dt

	// Step 4: did we crash?
	// If n is not zero, we touched a wall. "into" measures how much of our speed went straight
	// at the wall: driving head-on gives a big number, sliding along it a small one. Hitting the
	// wall costs 40% of the speed. Only a hard hit (into > 6) that comes more than 0.4 seconds
	// after the last one is reported back with return, so the game plays the crash sound.
	// return ends the method right there and hands the number to whoever called Update.
	if n.X != 0 || n.Z != 0 {
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

// Draw puts the truck on the screen. First a dark see-through disc on the ground as its shadow,
// then the 3D truck model at its position, turned by its heading. Rad2deg converts the angle from
// radians to degrees, because the drawing function wants degrees.
func (c *Car) Draw(a *Assets) {
	rl.DrawCylinder(rl.NewVector3(c.Pos.X, 0.07, c.Pos.Z), 1.3, 1.3, 0.01, 20, rl.NewColor(0, 0, 0, 70))
	a.Draw(c.Model, c.Pos, c.Heading*rl.Rad2deg, 1)
}
