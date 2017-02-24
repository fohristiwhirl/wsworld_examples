package main

// This example is a bit of a scam, since programming the whole
// thing in JavaScript would be more performant than this hybrid.

import (
    "math"
    "math/rand"

    ws "github.com/fohristiwhirl/wsworld"
)

const (
    FPS = 60        // This value is a lie in this example; the actual FPS is based on how often the client requests a frame.
    WIDTH = 1750
    HEIGHT = 850
    QUEENS = 8      // If your computer can handle it, try 30 queens and 4500 beasts.
    BEASTS = 1900
    BEAST_MAX_SPEED = 7
    QUEEN_MAX_SPEED = 5.5
    BEAST_ACCEL_MODIFIER = 0.55
    QUEEN_ACCEL_MODIFIER = 0.7
    QUEEN_TURN_PROB = 0.001
    BEAST_TURN_PROB = 0.002
    AVOID_STRENGTH = 4000
    MAX_PLAYER_SPEED = 10
    MARGIN = 50
)

const (
    QUEEN = iota
    BEAST
)

type Sim struct {
    tick int
    queens []*Dood
    beasts []*Dood
    player *ws.Entity
}

type Dood struct {
    ent *ws.Entity      // This is the part that the engine knows about, containing X, Y, Speedx, Speedy
    species int
    target *ws.Entity
    sim *Sim
}

func main() {

    ws.RegisterSprite("space ship.png")
    ws.Start("Swarmz 4.0", "127.0.0.1:8000", "/", "resources", WIDTH, HEIGHT)

    c := ws.NewCanvas(FPS)

    s := Sim{}
    s.Init(c)

    for {
        s.Iterate()
        ws.WaitForRequest()
        c.Send()
    }
}

func (s *Sim) Init(c *ws.Canvas) {

    for n := 0 ; n < QUEENS ; n++ {
        new_ent := c.NewPoint("#ffffff", WIDTH / 2, HEIGHT / 2, 0, 0)
        new_ent.Hidden = true
        s.queens = append(s.queens, &Dood{ent: new_ent, species: QUEEN, target: nil, sim: s})
    }

    for n := 0 ; n < BEASTS ; n++ {
        new_ent := c.NewPoint("#00ff00", WIDTH / 2, HEIGHT / 2, 0, 0)
        s.beasts = append(s.beasts, &Dood{ent: new_ent, species: BEAST, target: nil, sim: s})
    }

    s.player = c.NewSprite("space ship.png", 100, 100, 0, 0)
}

func (s *Sim) Iterate() {
    s.tick += 1
    s.MoveDoods()
    s.MovePlayer()
}

func (s *Sim) MoveDoods() {
    for _, d := range s.queens {
        d.Move()
    }
    for _, d := range s.beasts {
        d.Move()
    }
}

func (s *Sim) MovePlayer() {

    x, y, speedx, speedy := s.player.X, s.player.Y, s.player.Speedx, s.player.Speedy

    // Respond to input...

    if ws.KeyDown("w") { speedy -= 0.2 }
    if ws.KeyDown("a") { speedx -= 0.2 }
    if ws.KeyDown("s") { speedy += 0.2 }
    if ws.KeyDown("d") { speedx += 0.2 }

    // Bounce off walls...

    if (x < 16 && speedx < 0) || (x >  WIDTH - 16 && speedx > 0) { speedx *= -1 }
    if (y < 16 && speedy < 0) || (y > HEIGHT - 16 && speedy > 0) { speedy *= -1 }

    // Throttle speed...

    speed := math.Sqrt(speedx * speedx + speedy * speedy)

    if speed > MAX_PLAYER_SPEED {
        speedx *= MAX_PLAYER_SPEED / speed
        speedy *= MAX_PLAYER_SPEED / speed
    }

    // Update entity...

    s.player.Speedx = speedx
    s.player.Speedy = speedy
    s.player.Move()
}

func (d *Dood) Move() {

    ent := d.ent
    x, y, speedx, speedy := ent.X, ent.Y, ent.Speedx, ent.Speedy

    var turnprob, maxspeed, accelmod float64
    switch d.species {
    case QUEEN:
        turnprob = QUEEN_TURN_PROB
        maxspeed = QUEEN_MAX_SPEED
        accelmod = QUEEN_ACCEL_MODIFIER
    case BEAST:
        turnprob = BEAST_TURN_PROB
        maxspeed = BEAST_MAX_SPEED
        accelmod = BEAST_ACCEL_MODIFIER
    }

    // Chase target...

    if d.target == nil || rand.Float64() < turnprob || d.target == ent {
        tar_id := rand.Intn(QUEENS)
        d.target = d.sim.queens[tar_id].ent
    }

    vecx, vecy := unit_vector(x, y, d.target.X, d.target.Y)

    if vecx == 0 && vecy == 0 {
        speedx += rand.Float64() * 2 - 1 * accelmod
        speedy += rand.Float64() * 2 - 1 * accelmod
    } else {
        speedx += vecx * rand.Float64() * accelmod
        speedy += vecy * rand.Float64() * accelmod
    }

    // Wall avoidance...

    if (x < MARGIN) {
        speedx += rand.Float64() * 2
    }
    if (x >= WIDTH - MARGIN) {
        speedx -= rand.Float64() * 2
    }
    if (y < MARGIN) {
        speedy += rand.Float64() * 2
    }
    if (y >= HEIGHT - MARGIN) {
        speedy -= rand.Float64() * 2
    }

    // Player avoidance...

    dx := d.sim.player.X - x
    dy := d.sim.player.Y - y

    distance_squared := dx * dx + dy * dy
    distance := math.Sqrt(distance_squared)

    if distance > 1 {
        adjusted_force := AVOID_STRENGTH / (distance_squared * distance)
        speedx -= dx * adjusted_force * rand.Float64()
        speedy -= dy * adjusted_force * rand.Float64()
    }

    // Throttle speed...

    speed := math.Sqrt(ent.Speedx * ent.Speedx + ent.Speedy * ent.Speedy)

    if speed > maxspeed {
        speedx *= maxspeed / speed
        speedy *= maxspeed / speed
    }

    // Update entity...

    ent.Speedx = speedx
    ent.Speedy = speedy
    ent.Move()
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
