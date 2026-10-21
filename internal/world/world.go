package world

import (
	"math"

	"github.com/helwan-linux/helwan-penguin/internal/levels"
	"github.com/helwan-linux/helwan-penguin/internal/player"
)

const (
	TotalLevels       = 45
	GhostHitPoints    = 3
	ProjectileCost    = 20
	GhostKillCost     = 60
	GasDurationFrames = 300
)

type World struct {
	Number       int
	Level        int
	Platforms    []Platform
	Collectibles []Collectible
	Enemies      []Enemy
	Hazards      []Hazard
	Projectiles  []Projectile
	EndX         float64
	Frame        float64
	GasFrames    int
}

type EnemyKind int

const (
	EnemyWalker EnemyKind = iota
	EnemyRunner
	EnemyFlyer
	EnemyArmored
	EnemyBoss
	EnemyGhost
)

type Enemy struct {
	X, Y, Width, Height float64
	Active              bool
	Boss                bool
	Health              int
	Kind                EnemyKind
	BossVariant         int
	StartX              float64
	Patrol              float64
	Speed               float64
	Direction           float64
	BaseY               float64
	Phase               float64

	// World 8 and 9 enemy state.
	Ghost      bool
	Invisible  bool
	Stone      bool
	Visibility float64
	HitFlash   int
}

type Projectile struct {
	X, Y, Width, Height float64
	VX                  float64
	Active              bool
	Water               bool
}

func (e *Enemy) Update() {
	if !e.Active {
		return
	}

	if e.Stone {
		return
	}

	if e.Direction == 0 {
		e.Direction = 1
	}

	if e.Ghost {
		e.Phase += 0.045
		e.Visibility += 1

		// Ghosts alternate between visible and invisible.
		// Invisible ghosts cannot be hit or hurt the player.
		e.Invisible = int(e.Visibility/75)%2 == 1

		if e.Patrol > 0 {
			e.X += e.Speed * e.Direction

			if e.X < e.StartX-e.Patrol {
				e.X = e.StartX - e.Patrol
				e.Direction = 1
			} else if e.X > e.StartX+e.Patrol {
				e.X = e.StartX + e.Patrol
				e.Direction = -1
			}
		}

		if e.HitFlash > 0 {
			e.HitFlash--
		}

		return
	}

	if e.Boss {
		e.Phase += 0.045
		speed := e.Speed

		switch e.BossVariant {
		case 1:
			speed *= 1.45
		case 2:
			e.Y = e.BaseY + math.Sin(e.Phase)*30
		case 3:
			e.Y = e.BaseY - math.Abs(math.Sin(e.Phase))*25
		case 4:
			speed *= 0.55 + math.Abs(math.Sin(e.Phase))*1.7
		case 5:
			if math.Sin(e.Phase) > 0.55 {
				speed *= 2.1
			} else {
				speed *= 0.45
			}
		case 6:
			e.Y = e.BaseY + math.Sin(e.Phase)*36
			speed *= 0.9
		case 7:
			e.Y = e.BaseY + math.Sin(e.Phase*1.4)*42
			speed *= 0.7 + math.Abs(math.Sin(e.Phase))*1.3
		case 8:
			e.Y = e.BaseY + math.Sin(e.Phase*1.7)*22
			speed *= 1.1
		case 9:
			e.Y = e.BaseY - math.Abs(math.Sin(e.Phase))*35
			speed *= 1.25
		}

		e.X += speed * e.Direction
	} else if e.Kind == EnemyFlyer {
		e.X += e.Speed * e.Direction
		e.Y = e.BaseY + math.Sin(e.Phase)*24
		e.Phase += 0.055
	} else {
		e.X += e.Speed * e.Direction
	}

	if e.Patrol > 0 &&
		(e.X < e.StartX-e.Patrol || e.X > e.StartX+e.Patrol) {
		e.Direction *= -1

		if e.X < e.StartX-e.Patrol {
			e.X = e.StartX - e.Patrol
		} else {
			e.X = e.StartX + e.Patrol
		}
	}
}

func platformAt(x, y, width float64) Platform {
	return Platform{
		X:     x,
		Y:     y,
		W:     width,
		H:     540 - y,
		BaseX: x,
		BaseY: y,
	}
}

func NewLevel(level int) World {
	if level < 1 {
		level = 1
	}
	if level > TotalLevels {
		level = TotalLevels
	}

	worldNumber := (level-1)/5 + 1
	stage := (level - 1) % 5

	specs := levels.Route(worldNumber, stage+1)

	if len(specs) < 2 {
		specs = levels.Route(1, 1)
	}

	platforms := make([]Platform, 0, len(specs))

	for _, spec := range specs {
		platforms = append(platforms, platformAt(
			spec.X,
			spec.Y,
			spec.W,
		))
	}

	for i := range platforms {
		p := &platforms[i]
		p.BaseX = p.X
		p.BaseY = p.Y

		if i == 0 || i == len(platforms)-1 {
			continue
		}

		switch worldNumber {
		case 1:
			if stage == 0 {
				continue
			}
			if (i+stage)%3 == 0 {
				p.Kind = PlatformIce
			}

		case 2:
			if i%2 == 1 {
				p.Kind = PlatformMoving
				p.MoveRange = 12 + float64(stage*3)
				p.Speed = 0.35 + float64(stage)*0.05
			}

		case 3:
			if (i+stage)%2 == 0 {
				p.Kind = PlatformCrumbling
			}

		case 4:
			if (i+stage)%2 == 1 {
				p.Kind = PlatformCrumbling
			}

		case 5:
			if i%2 == 1 {
				p.Kind = PlatformConveyor
				p.ConveyorSpeed = 0.45 + float64(stage)*0.1
			}

		case 6:
			if i%2 == 0 {
				p.Kind = PlatformMoving
				p.MoveRange = 18 + float64(stage*4)
				p.Speed = 0.45 + float64(stage)*0.06
				p.MoveVertical = true
			}

		case 7:
			if i%2 == 0 {
				p.Kind = PlatformMoving
				p.MoveRange = 22 + float64(stage*5)
				p.Speed = 0.55 + float64(stage)*0.07
				p.MoveVertical = true
			} else if (i+stage)%3 == 0 {
				p.Kind = PlatformIce
			}

		case 8:
			// Shadow Realm: stable platforms; the danger comes from ghosts.
			if (i+stage)%4 == 2 {
				p.Kind = PlatformMoving
				p.MoveRange = 12 + float64(stage*2)
				p.Speed = 0.4 + float64(stage)*0.04
				p.MoveVertical = true
			}

		case 9:
			// Volcanic Fury: moving platforms above lava.
			if i%3 == 1 {
				p.Kind = PlatformMoving
				p.MoveRange = 10 + float64(stage*3)
				p.Speed = 0.45 + float64(stage)*0.05
			}
		}
	}

	last := platforms[len(platforms)-1]

	w := World{
		Number:    worldNumber,
		Level:     level,
		Platforms: platforms,
	}

	w.EndX = last.X + last.W - 48

	w.makeCollectibles(stage)
	w.makeEnemies(stage)
	w.makeHazards(stage)

	return w
}

func (w *World) makeCollectibles(stage int) {
	w.Collectibles = make([]Collectible, 0, len(w.Platforms))

	// Worlds 8 and 9 have no collectibles. Ammunition is limited.
	if w.Number >= 8 {
		return
	}

	for i, p := range w.Platforms[:len(w.Platforms)-1] {
		if (i+stage)%4 == 3 {
			continue
		}

		y := p.Y - 28 - float64((i%2)*12)

		if p.Kind == PlatformMoving && p.MoveVertical {
			y -= 10
		}

		w.Collectibles = append(w.Collectibles, Collectible{
			X:     p.X + p.W*0.55,
			Y:     y,
			Taken: false,
		})
	}
}

func (w *World) makeEnemies(stage int) {
	w.Enemies = nil

	if w.Level == 1 {
		return
	}

	// Worlds 8 and 9 use ghosts with three hit points.
	if w.Number >= 8 {
		for i, p := range w.Platforms[1 : len(w.Platforms)-1] {
			if i%2 == 1 && stage < 2 {
				continue
			}

			width, height := 34.0, 36.0
			health := GhostHitPoints
			boss := false

			if w.Level%5 == 0 && i == len(w.Platforms[1:len(w.Platforms)-1])-1 {
				width, height = 62, 62
				boss = true
				health = GhostHitPoints
			}

			w.Enemies = append(w.Enemies, Enemy{
				X:          p.X + p.W*0.5 - width/2,
				Y:          p.Y - height,
				Width:      width,
				Height:     height,
				Active:     true,
				Boss:       boss,
				Kind:       EnemyGhost,
				Health:     health,
				Ghost:      true,
				StartX:     p.X + p.W*0.5,
				Patrol:     p.W * 0.3,
				Speed:      0.45 + float64(stage)*0.08 + float64(w.Number)*0.04,
				Direction:  1,
				BaseY:      p.Y - height,
				Phase:      float64(i+stage) * 0.7,
				Visibility: float64(i*37 + stage*19),
			})
		}

		if w.Level%5 == 0 {
			p := w.Platforms[len(w.Platforms)-1]
			w.Enemies = append(w.Enemies, Enemy{
				X:          p.X + p.W*0.55,
				Y:          p.Y - 68,
				Width:      62,
				Height:     62,
				Active:     true,
				Boss:       true,
				Kind:       EnemyGhost,
				Health:     GhostHitPoints,
				Ghost:      true,
				StartX:     p.X + p.W*0.58,
				Patrol:     65,
				Speed:      0.8 + float64(w.Number)*0.08,
				Direction:  -1,
				BaseY:      p.Y - 68,
				Phase:      float64(w.Number) * 0.8,
				Visibility: float64(w.Number * 23),
			})
		}

		return
	}

	for i, p := range w.Platforms[1 : len(w.Platforms)-1] {
		if i%2 == 1 && stage < 2 {
			continue
		}

		kind := EnemyWalker

		switch w.Number {
		case 1:
			kind = EnemyWalker
		case 2:
			kind = EnemyFlyer
		case 3:
			kind = EnemyArmored
		case 4:
			kind = EnemyRunner
		case 5:
			kind = EnemyArmored
		case 6:
			kind = EnemyFlyer
		case 7:
			if i%2 == 0 {
				kind = EnemyArmored
			} else {
				kind = EnemyFlyer
			}
		}

		if stage >= 2 && i%3 == 1 {
			kind = EnemyRunner
		}

		width, height := 30.0, 28.0
		y := p.Y - height

		if kind == EnemyFlyer {
			y = p.Y - 72
			width, height = 34, 24
		}

		speed := 0.5 + float64(stage)*0.12 + float64(w.Number)*0.06

		if kind == EnemyRunner {
			speed *= 1.4
		}
		if kind == EnemyArmored {
			speed *= 0.65
		}
		if w.Number == 7 {
			speed *= 1.15
		}

		health := 1
		if kind == EnemyArmored {
			health = 2
		}
		if w.Number == 7 && kind == EnemyArmored {
			health = 3
		}

		w.Enemies = append(w.Enemies, Enemy{
			X:         p.X + p.W*0.5 - width/2,
			Y:         y,
			Width:     width,
			Height:    height,
			Active:    true,
			Kind:      kind,
			Health:    health,
			StartX:    p.X + p.W*0.5,
			Patrol:    p.W * 0.25,
			Speed:     speed,
			Direction: 1,
			BaseY:     y,
			Phase:     float64(i+stage) * 0.7,
		})
	}

	if w.Level%5 == 0 {
		p := w.Platforms[len(w.Platforms)-1]
		health := 3 + w.Number/2

		w.Enemies = append(w.Enemies, Enemy{
			X:           p.X + p.W*0.55,
			Y:           p.Y - 68,
			Width:       62,
			Height:      62,
			Active:      true,
			Boss:        true,
			Kind:        EnemyBoss,
			BossVariant: w.Number,
			Health:      health,
			StartX:      p.X + p.W*0.58,
			Patrol:      65,
			Speed:       0.7 + float64(w.Number)*0.12,
			Direction:   -1,
			BaseY:       p.Y - 68,
			Phase:       float64(w.Number) * 0.8,
		})
	}
}

func (w *World) makeHazards(stage int) {
	w.Hazards = nil

	if w.Level == 1 {
		return
	}

	for i := 0; i < len(w.Platforms)-1; i++ {
		left, right := w.Platforms[i], w.Platforms[i+1]
		gapStart := left.X + left.W
		gapWidth := right.X - gapStart

		if gapWidth < 40 {
			continue
		}

		kind := HazardPit

		switch w.Number {
		case 1:
			if (i+stage)%3 == 0 {
				kind = HazardIceSpikes
			}
		case 2:
			if (i+stage)%2 == 0 {
				kind = HazardThorns
			}
		case 3:
			kind = HazardLava
		case 4:
			if (i+stage)%2 == 0 {
				kind = HazardQuicksand
			}
		case 5:
			if (i+stage)%2 == 0 {
				kind = HazardLaser
			}
		case 6:
			if (i+stage)%2 == 0 {
				kind = HazardEnergy
			}
		case 7:
			if (i+stage)%2 == 0 {
				kind = HazardEnergy
			} else {
				kind = HazardLaser
			}
		case 8:
			// Shadow Realm gaps are dangerous pits.
			kind = HazardPit
		case 9:
			// Lava fills the gaps in the volcanic world.
			kind = HazardLava
		}

		if (i+stage)%3 != 1 ||
			w.Number == 3 ||
			w.Number == 5 ||
			w.Number == 7 ||
			w.Number == 8 ||
			w.Number == 9 {
			y, h := 438.0, 102.0

			if kind == HazardIceSpikes ||
				kind == HazardThorns ||
				kind == HazardLaser ||
				kind == HazardEnergy {
				y, h = left.Y-18, 18
			}

			hazard := Hazard{
				X:    gapStart + 8,
				Y:    y,
				W:    gapWidth - 16,
				H:    h,
				Kind: kind,
			}

			if kind == HazardLaser || kind == HazardEnergy {
				hazard.Period = 150 - float64(stage*8)
				hazard.ActiveFor = 82 - float64(stage*5)
				hazard.Phase = float64(i*29 + w.Level*7)
			}

			w.Hazards = append(w.Hazards, hazard)
		}
	}

	if w.Number >= 3 {
		for i, p := range w.Platforms[1 : len(w.Platforms)-1] {
			if (i+stage)%3 != 0 || p.W <= 100 {
				continue
			}

			kind := HazardLava

			switch w.Number {
			case 4:
				kind = HazardQuicksand
			case 5:
				kind = HazardLaser
			case 6, 7:
				kind = HazardEnergy
			case 8:
				kind = HazardEnergy
			case 9:
				kind = HazardLava
			}

			h := Hazard{
				X:    p.X + p.W*0.35,
				Y:    p.Y - 12,
				W:    p.W * 0.3,
				H:    12,
				Kind: kind,
			}

			if kind == HazardLaser || kind == HazardEnergy {
				h.Period = 140 - float64(stage*7)
				h.ActiveFor = 68 - float64(stage*4)
				h.Phase = float64(i*37 + w.Level*11)
			}

			w.Hazards = append(w.Hazards, h)
		}
	}

	if w.Number >= 3 {
		for i, p := range w.Platforms[1 : len(w.Platforms)-1] {
			if (i+stage+w.Level)%2 != 0 {
				continue
			}

			kind := HazardFallingRock

			switch w.Number {
			case 4:
				kind = HazardQuicksand
			case 5:
				kind = HazardLaser
			case 6, 7:
				kind = HazardEnergy
			case 8:
				kind = HazardEnergy
			case 9:
				kind = HazardFallingRock
			}

			h := Hazard{
				X:    p.X + p.W*0.42,
				Y:    p.Y - 36,
				W:    24,
				H:    24,
				Kind: kind,
			}

			if kind == HazardFallingRock {
				h.BaseY = p.Y - 80
				h.Period = 105 - float64(stage*6)
				h.ActiveFor = 75
				h.Phase = float64(i*31 + w.Level*9)
				h.Speed = 3.2 + float64(stage)*0.4
				h.Falling = true
			} else if kind == HazardLaser || kind == HazardEnergy {
				h.Y = p.Y - 14
				h.H = 14
				h.W = p.W * 0.55
				h.Period = 140 - float64(stage*7)
				h.ActiveFor = 68 - float64(stage*4)
				h.Phase = float64(i*37 + w.Level*11)
			}

			w.Hazards = append(w.Hazards, h)

			if w.Number == 4 && (i+stage)%2 == 1 {
				w.Hazards = append(w.Hazards, Hazard{
					X:    p.X + p.W*0.12,
					Y:    p.Y - 86,
					W:    p.W * 0.72,
					H:    42,
					Kind: HazardWind,
				})
			}

			// Volcanic vents periodically release gas.
			if w.Number == 9 && (i+stage)%2 == 0 {
				w.Hazards = append(w.Hazards, Hazard{
					X:         p.X + p.W*0.2,
					Y:         p.Y - 52,
					W:         p.W * 0.6,
					H:         48,
					Kind:      HazardWind,
					Period:    240,
					ActiveFor: 55,
					Phase:     float64(i*41 + w.Level*13),
				})
			}
		}
	}
}

// FireProjectile launches one shot in the direction the player is facing.
// The caller must deduct ProjectileCost from the player's score first.
func (w *World) FireProjectile(x, y, direction float64) bool {
	if w.Number != 8 && w.Number != 9 {
		return false
	}

	if direction == 0 {
		direction = 1
	}

	w.Projectiles = append(w.Projectiles, Projectile{
		X:      x,
		Y:      y,
		Width:  18,
		Height: 10,
		VX:     direction * 8,
		Active: true,
		Water:  w.Number == 9,
	})

	return true
}

// Update advances hazards, enemies, projectiles and moving platforms.
func (w *World) Update(playerState *player.Player) {
	w.Frame++

	if w.GasFrames > 0 {
		w.GasFrames--
	}

	for i := range w.Hazards {
		h := &w.Hazards[i]

		if h.Kind == HazardWind &&
			w.Number == 9 &&
			h.Period > 0 &&
			h.ActiveAt(w.Frame) &&
			playerState != nil &&
			playerState.X+playerState.Width > h.X &&
			playerState.X < h.X+h.W &&
			playerState.Y+playerState.Height > h.Y &&
			playerState.Y < h.Y+h.H {
			w.GasFrames = GasDurationFrames
		}

		if h.Kind != HazardFallingRock || !h.Falling {
			continue
		}

		if h.ActiveAt(w.Frame) {
			h.Y += h.Speed
		} else {
			h.Y = h.BaseY
		}
	}

	for i := range w.Enemies {
		w.Enemies[i].Update()
	}

	w.updateProjectiles()

	for i := range w.Platforms {
		p := &w.Platforms[i]

		if p.Kind == PlatformCrumbling && p.CrumbleTimer > 0 {
			p.CrumbleTimer--
			if p.CrumbleTimer <= 0 {
				p.Broken = true
			}
		}

		if p.Kind != PlatformMoving || p.MoveRange <= 0 || p.Broken {
			continue
		}

		oldX, oldY := p.X, p.Y
		offset := math.Sin(w.Frame*p.Speed/30+float64(i)) * p.MoveRange

		if p.MoveVertical {
			p.Y = p.BaseY + offset
		} else {
			p.X = p.BaseX + offset
		}

		if playerState != nil &&
			playerState.OnGround &&
			playerState.X+playerState.Width > oldX &&
			playerState.X < oldX+p.W &&
			math.Abs(playerState.Y+playerState.Height-oldY) < 5 {
			playerState.X += p.X - oldX
			playerState.Y += p.Y - oldY
		}
	}
}

func (w *World) updateProjectiles() {
	for i := range w.Projectiles {
		p := &w.Projectiles[i]

		if !p.Active {
			continue
		}

		p.X += p.VX

		if p.X < 0 || p.X > w.EndX+120 {
			p.Active = false
			continue
		}

		for j := range w.Enemies {
			e := &w.Enemies[j]

			if !e.Active || e.Stone || (e.Ghost && e.Invisible) {
				continue
			}

			if p.X+p.Width <= e.X ||
				p.X >= e.X+e.Width ||
				p.Y+p.Height <= e.Y ||
				p.Y >= e.Y+e.Height {
				continue
			}

			// Both fire and water need three hits to defeat a ghost.
			if e.Ghost {
				e.Health--
				e.HitFlash = 8

				if e.Health <= 0 {
					if w.Number == 9 {
						// A defeated volcanic ghost becomes a deadly stone.
						e.Stone = true
						e.Ghost = false
						e.Kind = EnemyArmored
						e.Width = 28
						e.Height = 22
						e.Y = e.BaseY + e.Height - 22
						e.Speed = 0
						e.Patrol = 0
						e.Invisible = false
					} else {
						e.Active = false
					}
				}
			}

			p.Active = false
			break
		}
	}

	active := w.Projectiles[:0]
	for _, p := range w.Projectiles {
		if p.Active {
			active = append(active, p)
		}
	}
	w.Projectiles = active
}

// MovementMultiplier returns the player movement multiplier during volcanic gas.
func (w World) MovementMultiplier() float64 {
	if w.Number == 9 && w.GasFrames > 0 {
		return 0.45
	}

	return 1
}

func (w *World) ResolvePlayer(p *player.Player) {
	p.OnGround = false
	previousY := p.Y - p.VY

	for i := range w.Platforms {
		platform := &w.Platforms[i]

		if platform.Broken {
			continue
		}

		overlapsX := p.X+p.Width > platform.X &&
			p.X < platform.X+platform.W
		previousBottom := previousY + p.Height
		currentBottom := p.Y + p.Height

		if overlapsX &&
			p.VY >= 0 &&
			previousBottom <= platform.Y+4 &&
			currentBottom >= platform.Y {
			p.Y = platform.Y - p.Height
			p.VY = 0
			p.OnGround = true

			if platform.Kind == PlatformCrumbling &&
				platform.CrumbleTimer == 0 {
				platform.CrumbleTimer = 55
			}

			if platform.Kind == PlatformConveyor {
				p.X += platform.ConveyorSpeed
			}

			if platform.Kind == PlatformIce {
				p.X += p.VX * 0.25
			}
		}
	}

	for i := range w.Collectibles {
		c := &w.Collectibles[i]

		if !c.Taken &&
			p.X < c.X+12 &&
			p.X+p.Width > c.X-12 &&
			p.Y < c.Y+12 &&
			p.Y+p.Height > c.Y-12 {
			c.Taken = true
		}
	}
}

func (w World) CollectedCount() int {
	n := 0

	for _, c := range w.Collectibles {
		if c.Taken {
			n++
		}
	}

	return n
}

func (w World) BossDefeated() bool {
	for _, e := range w.Enemies {
		if e.Boss && e.Active && !e.Stone {
			return false
		}
	}

	return true
}

func (w World) HazardAt(x, y, width, height float64) bool {
	for _, h := range w.Hazards {
		if h.Kind == HazardPit ||
			h.Kind == HazardQuicksand ||
			h.Kind == HazardWind ||
			!h.ActiveAt(w.Frame) {
			continue
		}

		if x+width > h.X &&
			x < h.X+h.W &&
			y+height > h.Y &&
			y < h.Y+h.H {
			return true
		}
	}

	return false
}

func (w World) StoneAt(x, y, width, height float64) bool {
	for _, e := range w.Enemies {
		if !e.Active || !e.Stone {
			continue
		}

		if x+width > e.X &&
			x < e.X+e.Width &&
			y+height > e.Y &&
			y < e.Y+e.Height {
			return true
		}
	}

	return false
}

func (w *World) ApplyHazardEffects(p *player.Player) {
	for _, h := range w.Hazards {
		if p.X+p.Width <= h.X ||
			p.X >= h.X+h.W ||
			p.Y+p.Height <= h.Y ||
			p.Y >= h.Y+h.H {
			continue
		}

		switch h.Kind {
		case HazardQuicksand:
			p.VX *= 0.45
			if p.VY < 2.5 {
				p.VY += 0.35
			}

		case HazardWind:
			if w.Number != 9 {
				p.X += 0.8
			}
		}
	}
}

func FirstLevel() World {
	return NewLevel(1)
}
