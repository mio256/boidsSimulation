package main

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	winWidth, winHeight int = 1900, 1200
	numBoids            int = 500
)

type Vector struct {
	x, y float64
}

func (v *Vector) Add(v2 Vector) {
	v.x += v2.x
	v.y += v2.y
}

func (v *Vector) Sub(v2 Vector) {
	v.x -= v2.x
	v.y -= v2.y
}

func (v *Vector) Mul(scalar float64) {
	v.x *= scalar
	v.y *= scalar
}

func (v *Vector) Div(scalar float64) {
	v.x /= scalar
	v.y /= scalar
}

func (v *Vector) Limit(max float64) {
	mag := v.Mag()
	if mag > max {
		v.Div(mag)
		v.Mul(max)
	}
}

func (v *Vector) Mag() float64 {
	return math.Sqrt(v.x*v.x + v.y*v.y)
}

func (v *Vector) Distance(v2 Vector) float64 {
	dx := v2.x - v.x
	dy := v2.y - v.y
	return math.Sqrt(dx*dx + dy*dy)
}

type Boid struct {
	position Vector
	velocity Vector
	accel    Vector
}

func (b *Boid) Update() {
	b.velocity.Add(b.accel)
	b.velocity.Limit(2)
	b.position.Add(b.velocity)
	b.accel = Vector{0, 0}
}

func (b *Boid) ApplyForce(force Vector) {
	b.accel.Add(force)
}

func (b *Boid) Align(boids []Boid) Vector {
	perceptionRadius := 50.0
	steering := Vector{0, 0}
	total := 0
	for _, other := range boids {
		d := b.position.Distance(other.position)
		if other != *b && d < perceptionRadius {
			steering.Add(other.velocity)
			total++
		}
	}
	if total > 0 {
		steering.Div(float64(total))
		steering.Limit(2)
		steering.Sub(b.velocity)
		steering.Limit(0.05)
	}
	return steering
}

func (b *Boid) Cohesion(boids []Boid) Vector {
	perceptionRadius := 30.0
	steering := Vector{0, 0}
	total := 0
	for _, other := range boids {
		d := b.position.Distance(other.position)
		if other != *b && d < perceptionRadius {
			steering.Add(other.position)
			total++
		}
	}
	if total > 0 {
		steering.Div(float64(total))
		steering.Sub(b.position)
		steering.Limit(0.05)
	}
	return steering
}

func (b *Boid) Separation(boids []Boid) Vector {
	perceptionRadius := 20.0
	steering := Vector{0, 0}
	total := 0
	for _, other := range boids {
		d := b.position.Distance(other.position)
		if other != *b && d < perceptionRadius {
			diff := Vector{b.position.x - other.position.x, b.position.y - other.position.y}
			diff.Div(d)
			steering.Add(diff)
			total++
		}
	}
	if total > 0 {
		steering.Div(float64(total))
	}
	steering.Limit(0.05)
	return steering
}

func (b *Boid) Edges() {
	if b.position.x > float64(winWidth) {
		b.position.x = 0
	} else if b.position.x < 0 {
		b.position.x = float64(winWidth)
	}
	if b.position.y > float64(winHeight) {
		b.position.y = 0
	} else if b.position.y < 0 {
		b.position.y = float64(winHeight)
	}
}

func NewBoid(x, y float64) Boid {
	angle := rand.Float64() * 2 * math.Pi
	return Boid{
		position: Vector{x, y},
		velocity: Vector{math.Cos(angle), math.Sin(angle)},
		accel:    Vector{0, 0},
	}
}

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		fmt.Printf("Error initializing SDL: %s\n", err)
		return
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("Boids Simulation", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, int32(winWidth), int32(winHeight), sdl.WINDOW_SHOWN|sdl.WINDOW_FULLSCREEN)
	if err != nil {
		fmt.Printf("Error creating window: %s\n", err)
		return
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Printf("Error creating renderer: %s\n", err)
		return
	}
	defer renderer.Destroy()

	boids := make([]Boid, numBoids)
	for i := range boids {
		boids[i] = NewBoid(rand.Float64()*float64(winWidth), rand.Float64()*float64(winHeight))
	}

	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}

		for i := range boids {
			boids[i].Edges()
			align := boids[i].Align(boids)
			cohesion := boids[i].Cohesion(boids)
			separation := boids[i].Separation(boids)
			boids[i].ApplyForce(align)
			boids[i].ApplyForce(cohesion)
			boids[i].ApplyForce(separation)
			boids[i].Update()
		}

		renderer.SetDrawColor(0, 0, 0, 255)
		renderer.Clear()

		renderer.SetDrawColor(255, 255, 255, 255)
		for _, boid := range boids {
			renderer.DrawPoint(int32(boid.position.x), int32(boid.position.y))
		}

		renderer.Present()
		sdl.Delay(16)
	}
}
