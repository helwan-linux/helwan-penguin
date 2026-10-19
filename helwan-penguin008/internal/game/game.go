package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/helwan-linux/helwan-penguin/internal/graphics"
	"github.com/helwan-linux/helwan-penguin/internal/player"
	"github.com/helwan-linux/helwan-penguin/internal/progress"
	"github.com/helwan-linux/helwan-penguin/internal/world"
)

const screenWidth, screenHeight = 960, 540

type Game struct {
	player     player.Player
	world      world.World
	camera     Camera
	state      State
	background *graphics.WorldBackground
	animation  graphics.PenguinAnimation
	score      int
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
		save:       saved,
	}

	g.startLevel(saved.UnlockedLevel)
	return g
}

func (g *Game) startLevel(level int) {
	if level < 1 {
		level = 1
	}
	if level > 30 {
		level = 30
	}

	g.world = world.NewLevel(level)
	g.player = player.New(55, g.world.Platforms[0].Y-38)
	g.camera = Camera{}
	g.animation = graphics.NewPenguinAnimation()
	g.state = StatePlaying
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

func (g *Game) completeLevel() {
	level := g.world.Level
	stars := g.world.CollectedCount()

	wasCompleted := g.save.IsCompleted(level)
	if !wasCompleted {
		progress.MarkCompleted(&g.save, level, stars)
		g.saveProgress()
	}

	g.state = StateLevelComplete
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

	previousBottom := g.player.Y + g.player.Height

	g.player.Update()
	wasFalling := g.player.VY >= 0

	starsBefore := g.world.CollectedCount()
	g.world.ResolvePlayer(&g.player)
	starsAfter := g.world.CollectedCount()

	if starsAfter > starsBefore {
		g.score += (starsAfter - starsBefore) * 50
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
			if e.Boss {
				e.Health--

				if e.Health <= 0 {
					e.Active = false
					g.score += 500
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
			g.state = StateGameOver
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
		g.state = StateGameOver
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
		x := p.X - g.camera.X

		if x < -p.W || x > screenWidth {
			continue
		}

		ebitenutil.DrawRect(
			screen,
			x,
			p.Y,
			p.W,
			p.H,
			groundColor,
		)

		ebitenutil.DrawRect(
			screen,
			x,
			p.Y,
			p.W,
			6,
			topColor,
		)
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

		if e.Boss {
			bodyColor = color.RGBA{
				R: 95,
				G: 35,
				B: 135,
				A: 255,
			}

			headColor = color.RGBA{
				R: 145,
				G: 65,
				B: 185,
				A: 255,
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

	flagX := g.world.EndX - g.camera.X
	lastPlatform := g.world.Platforms[len(g.world.Platforms)-1]

	ebitenutil.DrawRect(
		screen,
		flagX,
		lastPlatform.Y-70,
		5,
		70,
		color.RGBA{R: 50, G: 60, B: 70, A: 255},
	)

	ebitenutil.DrawRect(
		screen,
		flagX+5,
		lastPlatform.Y-70,
		35,
		22,
		color.RGBA{R: 55, G: 190, B: 105, A: 255},
	)

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
			"/30  Stars: "+itoa(g.world.CollectedCount())+
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
			"TRY AGAIN",
			"Press ENTER to restart",
			"",
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
