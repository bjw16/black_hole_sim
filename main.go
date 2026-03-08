package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

const (
	width        = 320
	height       = 90
	cx           = 160.0
	cy           = 45.0
	numParticles = 15000
	inSpeed      = 0.20
	spiralSpeed  = 0.06
)

type Particle struct {
	x, y float64
}

func main() {
	rand.Seed(time.Now().UnixNano())

	particles := make([]Particle, numParticles)
	for i := range particles {
		particles[i].x = float64(2 + rand.Intn(width-4))
		particles[i].y = float64(1 + rand.Intn(height-2))
	}

	clear := "\x1b[2J\x1b[H"
	blackHolePos := "\x1b[46;161H@"

	for {
		fmt.Print(clear)

		grid := make([][]int, height)
		for i := range grid {
			grid[i] = make([]int, width)
		}

		for i := range particles {
			p := &particles[i]

			dx := cx - p.x
			dy := cy - p.y
			dist := math.Sqrt(dx*dx + dy*dy)
			if dist < 1.5 {
				// Respawn particle outside
				p.x = float64(2 + rand.Intn(width-4))
				p.y = float64(1 + rand.Intn(height-2))
				continue
			}

			// Radial velocity towards center
			vx := (dx / dist) * inSpeed
			vy := (dy / dist) * inSpeed

			// Spiral: perpendicular velocity (counter-clockwise)
			angle := math.Atan2(dy, dx)
			perpAngle := angle + math.Pi/2
			spiralFactor := 25.0 / math.Max(dist, 1.0)
			vx += math.Cos(perpAngle) * spiralSpeed * spiralFactor
			vy += math.Sin(perpAngle) * spiralSpeed * spiralFactor

			p.x += vx
			p.y += vy

			// Plot particle, clamped to grid
			ix := int(math.Max(0, math.Min(float64(width-1), p.x)))
			iy := int(math.Max(0, math.Min(float64(height-1), p.y)))
			grid[iy][ix]++
		}

		// Zero black hole region
		bhRadius := 2.5
		for iy := 0; iy < height; iy++ {
			for ix := 0; ix < width; ix++ {
				dx := float64(ix) - cx
				dy := float64(iy) - cy
				if math.Sqrt(dx*dx+dy*dy) < bhRadius {
					grid[iy][ix] = 0
				}
			}
		}

		// Render grid with density chars
		for y := range grid {
			rowB := make([]byte, 0, width)
			for x := range grid[y] {
				count := grid[y][x]
				var ch byte
				switch {
				case count >= 12:
					ch = '*'
				case count >= 8:
					ch = '%'
				case count >= 5:
					ch = '#'
				case count >= 3:
					ch = '@'
				case count >= 1:
					ch = '.'
				default:
					ch = ' '
				}
				rowB = append(rowB, ch)
			}
			fmt.Println(string(rowB))
		}

		// Draw black hole at center
		fmt.Print(blackHolePos)

		time.Sleep(50 * time.Millisecond)
	}
}
