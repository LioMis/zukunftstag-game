package main

import (
	"math"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	GridN    = 17  // tiles per side
	TileSize = 4.0 // world units per tile
	BlockGap = 4   // a road every BlockGap tiles
)

// Directions as bits: N = -Z, E = +X, S = +Z, W = -X.
const (
	DirN = 1 << iota
	DirE
	DirS
	DirW
)

type Cell struct{ X, Z int }

type Tile struct {
	Road  bool
	Model string
	Rot   float32 // degrees around Y
	Solid bool
}

// Building is a delivery address: a building tile and the road spot in front of it.
type Building struct {
	Cell  Cell
	Front rl.Vector3 // marker position on the road in front of the door
	Name  string
}

type City struct {
	T         [GridN][GridN]Tile
	Roads     []Cell
	Buildings []Building
}

// Base connection masks of the road models at rotation 0.
var roadBase = map[string]int{
	"road-straight":     DirN | DirS,
	"road-corner":       DirS | DirE,
	"road-split":        DirE | DirS | DirW,
	"road-intersection": DirN | DirE | DirS | DirW,
}

// Direction a building model's door faces at rotation 0.
const buildingFront = DirS

// rot90 maps a direction mask through a +90° rotation about Y (E->N->W->S->E).
func rot90(m int) int {
	r := 0
	if m&DirE != 0 {
		r |= DirN
	}
	if m&DirN != 0 {
		r |= DirW
	}
	if m&DirW != 0 {
		r |= DirS
	}
	if m&DirS != 0 {
		r |= DirE
	}
	return r
}

func rotationFor(base, want int) float32 {
	m := base
	for i := 0; i < 4; i++ {
		if m == want {
			return float32(90 * i)
		}
		m = rot90(m)
	}
	return 0
}

func isRoadCell(x, z int) bool {
	return x >= 0 && z >= 0 && x < GridN && z < GridN && (x%BlockGap == 0 || z%BlockGap == 0)
}

func CellCenter(c Cell) rl.Vector3 {
	return rl.NewVector3((float32(c.X)+0.5)*TileSize, 0, (float32(c.Z)+0.5)*TileSize)
}

func dirVec(d int) (int, int) {
	switch d {
	case DirN:
		return 0, -1
	case DirE:
		return 1, 0
	case DirS:
		return 0, 1
	default:
		return -1, 0
	}
}

var streetNames = []string{"Bahnhofstrasse", "Seeweg", "Rosengasse", "Lindenplatz", "Marktgasse", "Bergstrasse",
	"Sonnenweg", "Kirchgasse", "Mühleweg", "Schulhausstrasse", "Bachweg", "Gartenstrasse", "Postgasse", "Parkweg"}

func NewCity(rng *rand.Rand) *City {
	c := &City{}
	buildings := []string{"building-small-a", "building-small-b", "building-small-c", "building-small-d", "building-garage"}
	for z := 0; z < GridN; z++ {
		for x := 0; x < GridN; x++ {
			t := &c.T[x][z]
			if isRoadCell(x, z) {
				t.Road = true
				mask := 0
				for _, d := range []int{DirN, DirE, DirS, DirW} {
					dx, dz := dirVec(d)
					if isRoadCell(x+dx, z+dz) {
						mask |= d
					}
				}
				n := 0
				for m := mask; m != 0; m &= m - 1 {
					n++
				}
				switch {
				case n == 4:
					t.Model = "road-intersection"
				case n == 3:
					t.Model = "road-split"
				case mask == DirN|DirS || mask == DirE|DirW:
					t.Model = "road-straight"
				default:
					t.Model = "road-corner"
				}
				t.Rot = rotationFor(roadBase[t.Model], mask)
				if t.Model == "road-straight" && (x+z)%2 == 1 {
					t.Model = "road-straight-lightposts"
				}
				c.Roads = append(c.Roads, Cell{x, z})
				continue
			}
			t.Solid = true
			// Position inside the 3x3 block.
			bx, bz := x%BlockGap, z%BlockGap
			if bx == 2 && bz == 2 {
				if rng.Intn(2) == 0 {
					t.Model = "pavement-fountain"
				} else {
					t.Model = "grass-trees-tall"
				}
				continue
			}
			if rng.Intn(9) == 0 {
				t.Model = "grass-trees"
				continue
			}
			// Face a neighbouring road.
			var front int
			for _, d := range []int{DirN, DirE, DirS, DirW} {
				dx, dz := dirVec(d)
				if isRoadCell(x+dx, z+dz) {
					front = d
					if rng.Intn(2) == 0 {
						break
					}
				}
			}
			t.Model = buildings[rng.Intn(len(buildings))]
			t.Rot = rotationFor(buildingFront, front)
			dx, dz := dirVec(front)
			road := Cell{x + dx, z + dz}
			p := CellCenter(road)
			p.X -= float32(dx) * TileSize * 0.22
			p.Z -= float32(dz) * TileSize * 0.22
			c.Buildings = append(c.Buildings, Building{Cell: Cell{x, z}, Front: p,
				Name: streetNames[rng.Intn(len(streetNames))]})
		}
	}
	return c
}

func (c *City) Draw(a *Assets) {
	for z := 0; z < GridN; z++ {
		for x := 0; x < GridN; x++ {
			t := c.T[x][z]
			a.Draw(t.Model, CellCenter(Cell{x, z}), t.Rot, TileSize)
		}
	}
	// Ring of trees around the town.
	for i := -1; i <= GridN; i++ {
		for _, cl := range []Cell{{i, -1}, {i, GridN}, {-1, i}, {GridN, i}} {
			a.Draw("grass-trees-tall", CellCenter(cl), float32((cl.X*7+cl.Z*3)%4)*90, TileSize)
		}
	}
}

func (c *City) solid(x, z int) bool {
	if x < 0 || z < 0 || x >= GridN || z >= GridN {
		return true
	}
	return c.T[x][z].Solid
}

// Collide pushes a circle at pos out of solid tiles. Returns the corrected
// position and the push normal (zero if no contact).
func (c *City) Collide(pos rl.Vector3, r float32) (rl.Vector3, rl.Vector3) {
	var normal rl.Vector3
	cx := int(math.Floor(float64(pos.X / TileSize)))
	cz := int(math.Floor(float64(pos.Z / TileSize)))
	for z := cz - 1; z <= cz+1; z++ {
		for x := cx - 1; x <= cx+1; x++ {
			if !c.solid(x, z) {
				continue
			}
			minX, minZ := float32(x)*TileSize, float32(z)*TileSize
			px := clamp(pos.X, minX, minX+TileSize)
			pz := clamp(pos.Z, minZ, minZ+TileSize)
			dx, dz := pos.X-px, pos.Z-pz
			d2 := dx*dx + dz*dz
			if d2 >= r*r {
				continue
			}
			d := float32(math.Sqrt(float64(d2)))
			if d < 1e-4 {
				continue // centre inside a tile; neighbours will resolve
			}
			push := r - d
			nx, nz := dx/d, dz/d
			pos.X += nx * push
			pos.Z += nz * push
			normal.X += nx
			normal.Z += nz
		}
	}
	return pos, rl.Vector3Normalize(normal)
}

func clamp(v, lo, hi float32) float32 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
