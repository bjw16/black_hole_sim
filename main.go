package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"time"
)

const (
	width        = 80
	height       = 24
	cx           = 40.0
	cy           = 12.0
	numParticles = 500
	inSpeed      = 0.25
	spiralSpeed  = 0.08
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
	blackHolePos := "\x1b[12;40H@"

	for {
		fmt.Print(clear)

		grid := make([][]byte, height)
		for i := range grid {
			grid[i] = bytes.Repeat([]byte{' '}, width)
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
			spiralFactor := 15.0 / math.Max(dist, 1.0)
			vx += math.Cos(perpAngle) * spiralSpeed * spiralFactor
			vy += math.Sin(perpAngle) * spiralSpeed * spiralFactor

			p.x += vx
			p.y += vy

			// Plot particle, clamped to grid
			ix := int(math.Max(0, math.Min(float64(width-1), p.x)))
			iy := int(math.Max(0, math.Min(float64(height-1), p.y)))
			grid[iy][ix] = '*'
		}

		// Render grid with newlines
		for _, row := range grid {
			fmt.Println(string(row))
		}

		// Draw black hole at center
		fmt.Print(blackHolePos)

		time.Sleep(60 * time.Millisecond)
	}
}
