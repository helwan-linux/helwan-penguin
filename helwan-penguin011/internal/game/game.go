package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	gameaudio "github.com/helwan-linux/helwan-penguin/internal/audio"
	"github.com/helwan-linux/helwan-penguin/internal/graphics"
	"github.com/helwan-linux/helwan-penguin/internal/player"
	"github.com/helwan-linux/helwan-penguin/internal/progress"
	"github.com/helwan-linux/helwan-penguin/internal/world"
)

const (
	screenWidth, screenHeight = 960, 540
	maxAttempts               = 5
)

type Game struct {
	player     player.Player
	world      world.World
	camera     Camera
	state      State
	background *graphics.WorldBackground
	animation  graphics.PenguinAnimation
	sound      *gameaudio.Manager
	score      int
	attempts   int
	save       progress.Progress
	saveError  string
}

func New() *Game {
	saved, err := progress.Load()
	if err != nil {
		saved = progress.Default()
	}

	g := &Game{
		state:      StateMenu,
		background: graphics.NewWorldBackground(),
		animation:  graphics.NewPenguinAnimation(),
		sound:      gameaudio.NewManager(),
		save:       saved,
	}

	g.startLevel(saved.UnlockedLevel)
	g.state = StateMenu
	return g
}

func (g *Game) startLevel(level int) {
	if level < 1 {
		level = 1
	}
	if level > 30 {
		level = 30
	}

	g.attempts = maxAttempts
	g.resetCurrentLevel(level)
}

// resetCurrentLevel restarts a level without restoring spent attempts.
func (g *Game) resetCurrentLevel(level int) {
	g.world = world.NewLevel(level)
	g.player = player.New(55, g.world.Platforms[0].Y-38)
	g.camera = Camera{}
	g.animation = graphics.NewPenguinAnimation()
	g.state = StatePlaying
}

// loseAttempt removes one attempt and restarts the current level if possible.
func (g *Game) loseAttempt() {
	g.playSound("hit")
	if g.attempts > 0 {
		g.attempts--
	}

	if g.attempts <= 0 {
		g.attempts = 0
		g.state = StateGameOver
		return
	}

	g.resetCurrentLevel(g.world.Level)
}

var previousKeys = make(map[ebiten.Key]bool)

func justPressed(key ebiten.Key) bool {
	return ebiten.IsKeyPressed(key) && !previousKeys[key]
}

func updateKeyHistory() {
	keys := []ebiten.Key{
		ebiten.KeyEnter,
		ebiten.KeyEscape,
		ebiten.KeyP,
		ebiten.KeySpace,
	}

	for _, key := range keys {
		previousKeys[key] = ebiten.IsKeyPressed(key)
	}
}

func (g *Game) saveProgress() {
	if err := progress.Save(g.save); err != nil {
		g.saveError = "Could not save progress"
		return
	}

	g.saveError = ""
}

func (g *Game) playSound(name string) {
	if g.sound != nil {
		g.sound.Play(name)
	}
}

func (g *Game) completeLevel() {
	level := g.world.Level
	stars := g.world.CollectedCount()

	wasCompleted := g.save.IsCompleted(level)
	if !wasCompleted {
		progress.MarkCompleted(&g.save, level, stars)
		g.saveProgress()
	}

	if level == 30 {
		g.playSound("win")
		g.state = StateGameWon
	} else {
		g.playSound("level")
		g.state = StateLevelComplete
	}
}

func (g *Game) Update() error {
	defer updateKeyHistory()

	if g.background != nil {
		g.background.Update()
	}

	if justPressed(ebiten.KeyP) {
		if g.state == StatePlaying {
			g.state = StatePaused
			return nil
		}

		if g.state == StatePaused {
			g.state = StatePlaying
			return nil
		}
	}

	if justPressed(ebiten.KeyEscape) {
		switch g.state {
		case StatePlaying:
			g.state = StatePaused
		case StatePaused:
			g.state = StatePlaying
		case StateMenu:
			return ebiten.Termination
		default:
			g.state = StateMenu
		}

		return nil
	}

	switch g.state {
	case StateMenu:
		if justPressed(ebiten.KeyEnter) ||
			justPressed(ebiten.KeySpace) {
			level := g.save.UnlockedLevel

			if level < 1 {
				level = 1
			}
			if level > 30 {
				level = 30
			}

			g.startLevel(level)
		}

		return nil

	case StatePaused:
		if justPressed(ebiten.KeyEnter) {
			g.state = StatePlaying
		}

		return nil

	case StateGameWon:
		if justPressed(ebiten.KeyEnter) || justPressed(ebiten.KeySpace) {
			g.state = StateMenu
		}
		return nil

	case StateLevelComplete:
		if justPressed(ebiten.KeyEnter) {
			if g.world.Level < 30 {
				next := g.world.Level + 1

				if g.save.IsUnlocked(next) {
					g.startLevel(next)
				} else {
					g.state = StateMenu
				}
			} else {
				g.state = StateMenu
			}
		}

		return nil

	case StateGameOver:
		if justPressed(ebiten.KeyEnter) {
			g.startLevel(g.world.Level)
		}

		return nil
	}

	g.world.Update(&g.player)
	previousBottom := g.player.Y + g.player.Height

	wasOnGround := g.player.OnGround
	g.player.Update()
	if wasOnGround && g.player.VY < 0 {
		g.playSound("jump")
	}
	if g.world.Number == 6 {
		// Space world: low gravity creates longer, floatier jumps.
		g.player.VY *= 0.94
	}
	wasFalling := g.player.VY >= 0

	starsBefore := g.world.CollectedCount()
	g.world.ResolvePlayer(&g.player)
	g.world.ApplyHazardEffects(&g.player)
	if g.world.HazardAt(g.player.X, g.player.Y, g.player.Width, g.player.Height) {
		g.loseAttempt()
		return nil
	}
	starsAfter := g.world.CollectedCount()

	if starsAfter > starsBefore {
		g.score += (starsAfter - starsBefore) * 50
		g.playSound("collect")
	}

	for i := range g.world.Enemies {
		e := &g.world.Enemies[i]

		if !e.Active {
			continue
		}

		e.Update()

		overlapsX := g.player.X+g.player.Width > e.X &&
			g.player.X < e.X+e.Width

		overlapsY := g.player.Y+g.player.Height > e.Y &&
			g.player.Y < e.Y+e.Height

		if !overlapsX || !overlapsY {
			continue
		}

		stomped := wasFalling && previousBottom <= e.Y+10

		if stomped {
			g.playSound("stomp")
			if e.Boss || e.Kind == world.EnemyArmored {
				e.Health--
				if e.Health <= 0 {
					e.Active = false
					if e.Boss {
						g.score += 500
						g.playSound("boss")
					} else {
						g.score += 150
					}
				} else {
					g.score += 100
				}
			} else {
				e.Active = false
				g.score += 100
			}

			g.player.VY = -7.5
			g.player.OnGround = false
		} else {
			g.loseAttempt()
			return nil
		}
	}

	moving := ebiten.IsKeyPressed(ebiten.KeyA) ||
		ebiten.IsKeyPressed(ebiten.KeyD) ||
		ebiten.IsKeyPressed(ebiten.KeyArrowLeft) ||
		ebiten.IsKeyPressed(ebiten.KeyArrowRight)

	jumping := g.player.VY < 0 && !g.player.OnGround
	falling := g.player.VY > 0 && !g.player.OnGround

	g.animation.Update(moving, jumping, falling)

	if ebiten.IsKeyPressed(ebiten.KeyA) ||
		ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.animation.Facing = -1
	} else if ebiten.IsKeyPressed(ebiten.KeyD) ||
		ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.animation.Facing = 1
	}

	if g.player.Y > screenHeight+120 {
		g.loseAttempt()
		return nil
	}

	if g.player.X >= g.world.EndX {
		if g.world.Level%5 == 0 && !g.world.BossDefeated() {
			g.player.X = g.world.EndX - g.player.Width
		} else {
			g.completeLevel()
			return nil
		}
	}

	g.camera.Follow(g.player.X, g.player.Y)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.background != nil {
		g.background.Draw(screen, g.camera.X, g.world.Number)
	} else {
		screen.Fill(color.RGBA{
			R: 166,
			G: 218,
			B: 245,
			A: 255,
		})
	}

	groundColor, topColor := platformColors(g.world.Number)

	for _, p := range g.world.Platforms {
		if p.Broken {
			continue
		}
		x := p.X - g.camera.X

		if x < -p.W || x > screenWidth {
			continue
		}

		platformColor, platformTop := groundColor, topColor
		switch p.Kind {
		case world.PlatformIce:
			platformColor = color.RGBA{R: 120, G: 205, B: 238, A: 255}
			platformTop = color.RGBA{R: 235, G: 250, B: 255, A: 255}
		case world.PlatformMoving:
			platformColor = color.RGBA{R: 130, G: 160, B: 180, A: 255}
			platformTop = color.RGBA{R: 255, G: 220, B: 100, A: 255}
		case world.PlatformCrumbling:
			platformColor = color.RGBA{R: 155, G: 105, B: 65, A: 255}
			platformTop = color.RGBA{R: 235, G: 165, B: 85, A: 255}
		case world.PlatformConveyor:
			platformColor = color.RGBA{R: 95, G: 105, B: 115, A: 255}
			platformTop = color.RGBA{R: 245, G: 185, B: 55, A: 255}
		}
		ebitenutil.DrawRect(screen, x, p.Y, p.W, p.H, platformColor)
		ebitenutil.DrawRect(screen, x, p.Y, p.W, 6, platformTop)
		if p.Kind == world.PlatformConveyor {
			for stripe := 0.0; stripe < p.W; stripe += 28 {
				ebitenutil.DrawRect(screen, x+stripe, p.Y+8, 12, 3, color.RGBA{R: 35, G: 40, B: 45, A: 255})
			}
		}
	}

	for _, h := range g.world.Hazards {
		x := h.X - g.camera.X
		if x < -h.W || x > screenWidth {
			continue
		}
		hazardColor := color.RGBA{R: 35, G: 35, B: 45, A: 255}
		switch h.Kind {
		case world.HazardIceSpikes:
			hazardColor = color.RGBA{R: 190, G: 240, B: 255, A: 255}
		case world.HazardThorns:
			hazardColor = color.RGBA{R: 50, G: 150, B: 65, A: 255}
		case world.HazardLava:
			hazardColor = color.RGBA{R: 245, G: 75, B: 25, A: 255}
		case world.HazardQuicksand:
			hazardColor = color.RGBA{R: 210, G: 165, B: 70, A: 255}
		case world.HazardLaser:
			hazardColor = color.RGBA{R: 245, G: 35, B: 55, A: 255}
		case world.HazardEnergy:
			hazardColor = color.RGBA{R: 70, G: 210, B: 255, A: 255}
		case world.HazardFallingRock:
			hazardColor = color.RGBA{R: 105, G: 91, B: 85, A: 255}
		case world.HazardWind:
			hazardColor = color.RGBA{R: 190, G: 230, B: 255, A: 255}
		}
		if (h.Kind == world.HazardLaser || h.Kind == world.HazardEnergy || h.Kind == world.HazardFallingRock) && !h.ActiveAt(g.world.Frame) {
			hazardColor = color.RGBA{R: 85, G: 85, B: 95, A: 255}
		}
		if h.Kind == world.HazardFallingRock {
			ebitenutil.DrawCircle(screen, x+h.W/2, h.Y+h.H/2, h.W/2, hazardColor)
		} else if h.Kind == world.HazardQuicksand {
			ebitenutil.DrawRect(screen, x, h.Y, h.W, h.H, hazardColor)
			for stripe := 0.0; stripe < h.W; stripe += 16 {
				ebitenutil.DrawRect(screen, x+stripe, h.Y+4, 8, 3, color.RGBA{R: 245, G: 205, B: 110, A: 255})
			}
		} else if h.Kind == world.HazardWind {
			for stripe := 0.0; stripe < h.W; stripe += 20 {
				ebitenutil.DrawRect(screen, x+stripe, h.Y+3, 13, 3, hazardColor)
			}
		} else if h.Kind == world.HazardIceSpikes || h.Kind == world.HazardThorns {
			for spike := 0.0; spike < h.W; spike += 14 {
				ebitenutil.DrawRect(screen, x+spike, h.Y+5, 7, h.H-5, hazardColor)
			}
		} else {
			ebitenutil.DrawRect(screen, x, h.Y, h.W, h.H, hazardColor)
		}
	}

	for _, e := range g.world.Enemies {
		if !e.Active {
			continue
		}

		x := e.X - g.camera.X

		if x < -e.Width || x > screenWidth {
			continue
		}

		bodyColor := color.RGBA{
			R: 190,
			G: 55,
			B: 65,
			A: 255,
		}

		headColor := color.RGBA{
			R: 220,
			G: 75,
			B: 80,
			A: 255,
		}

		switch e.Kind {
		case world.EnemyRunner:
			bodyColor = color.RGBA{R: 215, G: 75, B: 45, A: 255}
			headColor = color.RGBA{R: 250, G: 115, B: 55, A: 255}
		case world.EnemyFlyer:
			bodyColor = color.RGBA{R: 45, G: 125, B: 190, A: 255}
			headColor = color.RGBA{R: 75, G: 185, B: 235, A: 255}
		case world.EnemyArmored:
			bodyColor = color.RGBA{R: 85, G: 95, B: 105, A: 255}
			headColor = color.RGBA{R: 155, G: 165, B: 175, A: 255}
		}
		if e.Boss {
			switch e.BossVariant {
			case 1: // ice guardian
				bodyColor, headColor = color.RGBA{R: 55, G: 130, B: 190, A: 255}, color.RGBA{R: 140, G: 225, B: 255, A: 255}
			case 2: // forest guardian
				bodyColor, headColor = color.RGBA{R: 45, G: 125, B: 55, A: 255}, color.RGBA{R: 125, G: 210, B: 85, A: 255}
			case 3: // volcano guardian
				bodyColor, headColor = color.RGBA{R: 145, G: 45, B: 25, A: 255}, color.RGBA{R: 255, G: 115, B: 35, A: 255}
			case 4: // desert guardian
				bodyColor, headColor = color.RGBA{R: 160, G: 105, B: 35, A: 255}, color.RGBA{R: 240, G: 195, B: 85, A: 255}
			case 5: // factory guardian
				bodyColor, headColor = color.RGBA{R: 65, G: 75, B: 90, A: 255}, color.RGBA{R: 80, G: 225, B: 225, A: 255}
			default: // space guardian
				bodyColor, headColor = color.RGBA{R: 90, G: 35, B: 145, A: 255}, color.RGBA{R: 180, G: 105, B: 255, A: 255}
			}
		}

		ebitenutil.DrawRect(
			screen,
			x,
			e.Y+7,
			e.Width,
			e.Height-7,
			bodyColor,
		)

		ebitenutil.DrawCircle(
			screen,
			x+e.Width/2,
			e.Y+9,
			e.Width/2,
			headColor,
		)

		ebitenutil.DrawCircle(
			screen,
			x+7,
			e.Y+9,
			3,
			color.RGBA{R: 255, G: 255, B: 255, A: 255},
		)

		ebitenutil.DrawCircle(
			screen,
			x+e.Width-7,
			e.Y+9,
			3,
			color.RGBA{R: 255, G: 255, B: 255, A: 255},
		)

		ebitenutil.DrawCircle(
			screen,
			x+7,
			e.Y+9,
			1.3,
			color.RGBA{R: 20, G: 20, B: 20, A: 255},
		)

		ebitenutil.DrawCircle(
			screen,
			x+e.Width-7,
			e.Y+9,
			1.3,
			color.RGBA{R: 20, G: 20, B: 20, A: 255},
		)

		if e.Boss {
			healthText := "Boss HP: " + itoa(e.Health)

			ebitenutil.DebugPrintAt(
				screen,
				healthText,
				int(x),
				int(e.Y)-18,
			)
		}
	}

	for _, c := range g.world.Collectibles {
		if c.Taken {
			continue
		}

		ebitenutil.DrawCircle(
			screen,
			c.X-g.camera.X,
			c.Y,
			9,
			color.RGBA{R: 255, G: 194, B: 40, A: 255},
		)
	}

	goalX := g.world.EndX - g.camera.X
	lastPlatform := g.world.Platforms[len(g.world.Platforms)-1]
	// The finish marker is a trophy, not a flag. Level 30 gets a larger golden cup.
	trophyGold := color.RGBA{R: 255, G: 190, B: 35, A: 255}
	trophyLight := color.RGBA{R: 255, G: 225, B: 95, A: 255}
	cupWidth := 34.0
	cupTop := lastPlatform.Y - 58
	if g.world.Level == 30 {
		cupWidth = 48
		cupTop = lastPlatform.Y - 72
	}
	// Cup bowl, handles, stem and pedestal are drawn from native shapes.
	ebitenutil.DrawRect(screen, goalX+cupWidth*0.28, cupTop+24, cupWidth*0.44, 18, trophyGold)
	ebitenutil.DrawRect(screen, goalX+cupWidth*0.40, cupTop+42, cupWidth*0.20, 12, trophyLight)
	ebitenutil.DrawRect(screen, goalX+cupWidth*0.24, cupTop+54, cupWidth*0.52, 6, trophyGold)
	ebitenutil.DrawRect(screen, goalX+cupWidth*0.08, cupTop+55, cupWidth*0.84, 6, trophyGold)
	ebitenutil.DrawCircle(screen, goalX+cupWidth/2, cupTop+18, cupWidth*0.40, trophyGold)
	ebitenutil.DrawRect(screen, goalX+cupWidth*0.05, cupTop+8, cupWidth*0.20, 18, trophyGold)
	ebitenutil.DrawRect(screen, goalX+cupWidth*0.75, cupTop+8, cupWidth*0.20, 18, trophyGold)
	ebitenutil.DrawRect(screen, goalX+cupWidth*0.27, cupTop+8, cupWidth*0.46, 18, trophyLight)

	if g.world.Level == 30 && g.state == StatePlaying {
		ebitenutil.DebugPrintAt(screen, "FINAL TROPHY", int(goalX-12), int(lastPlatform.Y-92))
	}

	graphics.DrawPenguin(
		screen,
		float32(g.player.X-g.camera.X),
		float32(g.player.Y),
		&g.animation,
	)

	ebitenutil.DebugPrintAt(
		screen,
		"HELWAN PENGUIN | WORLD "+itoa(g.world.Number),
		18,
		16,
	)

	ebitenutil.DebugPrintAt(
		screen,
		"Level: "+itoa(g.world.Level)+
			"/30  Attempts: "+itoa(g.attempts)+"/"+itoa(maxAttempts)+
			"  Stars: "+itoa(g.world.CollectedCount())+
			"  Score: "+itoa(g.score),
		18,
		38,
	)

	if g.saveError != "" {
		ebitenutil.DebugPrintAt(screen, g.saveError, 18, 82)
	}

	switch g.state {
	case StateMenu:
		drawOverlay(
			screen,
			"HELWAN PENGUIN",
			"Press ENTER to continue",
			"ESC to quit",
		)

		ebitenutil.DebugPrintAt(
			screen,
			"Unlocked level: "+itoa(g.save.UnlockedLevel)+
				"  Total stars: "+itoa(g.save.TotalStars),
			340,
			330,
		)

	case StatePaused:
		drawOverlay(
			screen,
			"PAUSED",
			"Press P or ENTER to resume",
			"ESC to resume",
		)

	case StateLevelComplete:
		drawOverlay(
			screen,
			"LEVEL COMPLETE!",
			"Press ENTER to continue",
			"",
		)

	case StateGameOver:
		drawOverlay(
			screen,
			"NO ATTEMPTS LEFT",
			"Press ENTER to retry this level",
			"A new game gives you 5 attempts",
		)

	case StateGameWon:
		drawOverlay(
			screen,
			"CONGRATULATIONS!",
			"YOU FINISHED HELWAN PENGUIN!",
			"All 30 levels complete - Press ENTER",
		)

	default:
		ebitenutil.DebugPrintAt(
			screen,
			"Move: A/D or arrows | Jump: Space/W/Up | Pause: P",
			18,
			60,
		)
	}
}

func platformColors(worldNumber int) (color.RGBA, color.RGBA) {
	switch worldNumber {
	case 2:
		return color.RGBA{R: 58, G: 105, B: 62, A: 255},
			color.RGBA{R: 145, G: 205, B: 105, A: 255}

	case 3:
		return color.RGBA{R: 105, G: 49, B: 45, A: 255},
			color.RGBA{R: 255, G: 125, B: 55, A: 255}

	case 4:
		return color.RGBA{R: 194, G: 151, B: 83, A: 255},
			color.RGBA{R: 255, G: 220, B: 135, A: 255}

	case 5:
		return color.RGBA{R: 65, G: 75, B: 91, A: 255},
			color.RGBA{R: 105, G: 205, B: 220, A: 255}

	case 6:
		return color.RGBA{R: 65, G: 48, B: 105, A: 255},
			color.RGBA{R: 165, G: 125, B: 255, A: 255}

	default:
		return color.RGBA{R: 116, G: 173, B: 194, A: 255},
			color.RGBA{R: 245, G: 252, B: 255, A: 255}
	}
}

func drawOverlay(
	screen *ebiten.Image,
	title string,
	subtitle string,
	footer string,
) {
	ebitenutil.DrawRect(
		screen,
		180,
		175,
		600,
		190,
		color.RGBA{R: 20, G: 35, B: 50, A: 235},
	)

	ebitenutil.DebugPrintAt(screen, title, 390, 220)
	ebitenutil.DebugPrintAt(screen, subtitle, 310, 260)

	if footer != "" {
		ebitenutil.DebugPrintAt(screen, footer, 390, 295)
	}
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}

	if n < 0 {
		return "-" + itoa(-n)
	}

	var b [20]byte
	i := len(b)

	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}

	return string(b[i:])
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
