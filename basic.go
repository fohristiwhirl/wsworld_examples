package main

import (
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

    var angle float64

    c := ws.NewCanvas()

    player := c.NewSprite("space ship.png", 100, 100, 0, 0)
    orbiter := c.NewSprite("globe.png", player.X, player.Y, 0, 0)

    for {
        if ws.KeyDown(-1, "w") && player.Speedy > -2 && player.Y > 16          { player.Speedy -= 0.1 }
        if ws.KeyDown(-1, "a") && player.Speedx > -2 && player.X > 16          { player.Speedx -= 0.1 }
        if ws.KeyDown(-1, "s") && player.Speedy <  2 && player.Y < HEIGHT - 16 { player.Speedy += 0.1 }
        if ws.KeyDown(-1, "d") && player.Speedx <  2 && player.X <  WIDTH - 16 { player.Speedx += 0.1 }

        if (player.X > WIDTH - 16 && player.Speedx > 0) || (player.X < 16 && player.Speedx < 0) {
            player.Speedx *= -1
            c.PlaySound("shot.wav")
        }
        if (player.Y > HEIGHT - 16 && player.Speedy > 0) || (player.Y < 16 && player.Speedy < 0) {
            player.Speedy *= -1
            c.PlaySound("shot.wav")
        }

        player.X += player.Speedx
        player.Y += player.Speedy

        angle += 0.03
        orbiter.X = player.X + 50 * math.Cos(angle)
        orbiter.Y = player.Y + 50 * math.Sin(angle)

        <- ticker
        c.SendToAll()
    }
}
