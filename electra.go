package main

import (
    "math"
    "math/rand"
    "time"
    ws "github.com/fohristiwhirl/wsworld"
)

const (
    FPS = 121
    WIDTH = 1750
    HEIGHT = 850

    SUBLINES = 20
    RADIUS = 380
)

type Line struct {
    x1 float64
    y1 float64
    x2 float64
    y2 float64
}

func main() {
    ws.RegisterSprite("globe.png")
    ws.Start("Electra", "127.0.0.1:8000", "/", "resources", WIDTH, HEIGHT, FPS, false)

    var ticker = time.Tick(time.Second / FPS)

    var centre_x float64 = WIDTH / 2
    var centre_y float64 = HEIGHT / 2

    c := ws.NewCanvas()

    var angle float64
    var i int

    var lines []Line

    for {
        i++
        c.Clear()

        angle += 0.005
        orbiter1_x := centre_x + RADIUS * math.Cos(angle)
        orbiter1_y := centre_y + RADIUS * math.Sin(angle)
        orbiter2_x := centre_x - RADIUS * math.Cos(angle)
        orbiter2_y := centre_y - RADIUS * math.Sin(angle)

        var x, y float64 = orbiter1_x, orbiter1_y
        var next_x, next_y float64

        if i % 5 == 0 {

            lines = nil

            for n := 0 ; n < SUBLINES ; n++ {

                vecx, vecy := unit_vector(x, y, orbiter2_x, orbiter2_y)

                dx := orbiter2_x - x
                dy := orbiter2_y - y
                distance := math.Sqrt(dx * dx + dy * dy)

                if n == SUBLINES - 1 {
                    next_x = orbiter2_x
                    next_y = orbiter2_y
                } else {
                    next_x = x + (vecx * distance / (SUBLINES - float64(n))) + (rand.Float64() * 40) - 20
                    next_y = y + (vecy * distance / (SUBLINES - float64(n))) + (rand.Float64() * 40) - 20
                }

                lines = append(lines, Line{x, y, next_x, next_y})

                x = next_x
                y = next_y
            }
        }

        for _, line := range lines {
            c.AddLine("#00ccff", line.x1, line.y1, line.x2, line.y2, 0, 0)
        }

        c.AddSprite("globe.png", orbiter1_x, orbiter1_y, 0, 0)
        c.AddSprite("globe.png", orbiter2_x, orbiter2_y, 0, 0)

        <- ticker
        c.SendToAll()
    }
}

func unit_vector(x1, y1, x2, y2 float64) (float64, float64) {
    dx := x2 - x1
    dy := y2 - y1

    if (dx == 0 && dy == 0) {
        return 0, 0
    }

    distance := math.Sqrt(dx * dx + dy * dy)
    return dx / distance, dy / distance
}
