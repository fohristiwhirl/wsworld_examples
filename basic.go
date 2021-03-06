package main

import (
    "fmt"
    "math"
    "time"
    ws "github.com/fohristiwhirl/wsworld"
)

const (
    FPS = 121       // Server speed can exceed client, that's fine
    WIDTH = 1750
    HEIGHT = 850
)

func main() {
    ws.RegisterSprite("space ship.png")
    ws.RegisterSprite("globe.png")
    ws.RegisterSound("shot.wav")
    ws.Start("Basic Example", "127.0.0.1:8000", "/", "resources", WIDTH, HEIGHT, FPS, false)

    var ticker = time.Tick(time.Second / FPS)

    var x, y, speedx, speedy, angle float64 = 100, 100, 0, 0, 0

    c := ws.NewCanvas()
    z := ws.NewSoundscape()

    for {

        c.Clear()
        z.Clear()       // Or sounds will play repeatedly every frame...

        if ws.KeyDown(-1, "w") && speedy > -2 && y > 16          { speedy -= 0.1 }
        if ws.KeyDown(-1, "a") && speedx > -2 && x > 16          { speedx -= 0.1 }
        if ws.KeyDown(-1, "s") && speedy <  2 && y < HEIGHT - 16 { speedy += 0.1 }
        if ws.KeyDown(-1, "d") && speedx <  2 && x <  WIDTH - 16 { speedx += 0.1 }

        if (x > WIDTH - 16 && speedx > 0) || (x < 16 && speedx < 0) {
            speedx *= -1
            z.PlaySound("shot.wav")
        }
        if (y > HEIGHT - 16 && speedy > 0) || (y < 16 && speedy < 0) {
            speedy *= -1
            z.PlaySound("shot.wav")
        }

        x += speedx
        y += speedy

        angle += 0.03
        orbiter_x := 50 * math.Cos(angle)
        orbiter_y := 50 * math.Sin(angle)

        c.AddLine("#ffff00", x, y, x + orbiter_x, y + orbiter_y, 0, 0)
        c.AddSprite("space ship.png", x, y, speedx, speedy)
        c.AddSprite("globe.png", x + orbiter_x, y + orbiter_y, 0, 0)
        c.AddText("Hello there: this works!", "#ffff00", 30, "Arial", x - orbiter_x, y - orbiter_y, 0, 0)

        <- ticker
        c.SendToAll()
        z.SendToAll()

        clicks := ws.PollClicks(-1)
        if len(clicks) > 0 {
            last_click := clicks[len(clicks) - 1]
            ws.SendDebugToAll(fmt.Sprintf("Click at %d, %d (Button %d)", last_click.X, last_click.Y, last_click.Button))
        }
    }
}
