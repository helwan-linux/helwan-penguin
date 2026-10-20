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
	maxLevel                  = 35
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
	infoPage   string
	infoReturn State
}

func New() *Game {
	saved, err := progress.Load()
	if err != nil {
		saved = progress.Default()
	}
	if saved.UnlockedLevel < 1 {
		saved.UnlockedLevel = 1
	}
	if saved.UnlockedLevel > maxLevel {
		saved.UnlockedLevel = maxLevel
	}

	g := &Game{
		state:      StateMenu,
		background: graphics.NewWorldBackground(),
		animation:  graphics.NewPenguinAnimation(),
		sound:      gameaudio.NewManager(),
		save:       saved,
	}
	g.startLevel(1)
	g.state = StateMenu
	return g
}

func (g *Game) startLevel(level int) {
	if level < 1 {
		level = 1
	}
	if level > maxLevel {
		level = maxLevel
	}
	g.attempts = maxAttempts
	g.resetCurrentLevel(level)
}

func (g *Game) resetCurrentLevel(level int) {
	g.world = world.NewLevel(level)
	if len(g.world.Platforms) == 0 {
		g.world = world.NewLevel(1)
	}
	g.player = player.New(55, g.world.Platforms[0].Y-38)
	g.camera = Camera{}
	g.animation = graphics.NewPenguinAnimation()
	g.state = StatePlaying
}

func (g *Game) startNewGame() {
	g.save = progress.Default()
	g.score = 0
	g.saveProgress()
	g.startLevel(1)
}

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
		ebiten.KeyEnter, ebiten.KeyEscape, ebiten.KeyP,
		ebiten.KeySpace, ebiten.KeyN, ebiten.KeyO,
		ebiten.KeyH,
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

func (g *Game) openInfo(page string) {
	g.infoReturn = g.state
	g.infoPage = page
}

func (g *Game) completeLevel() {
	level := g.world.Level
	stars := g.world.CollectedCount()

	if !g.save.IsCompleted(level) {
		progress.MarkCompleted(&g.save, level, stars)
		if g.save.UnlockedLevel > maxLevel {
			g.save.UnlockedLevel = maxLevel
		}
		g.saveProgress()
	}

	if level >= maxLevel {
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

	if g.infoPage != "" {
		if justPressed(ebiten.KeyEscape) ||
			justPressed(ebiten.KeyEnter) ||
			justPressed(ebiten.KeyO) ||
			justPressed(ebiten.KeyH) {
			g.infoPage = ""
			g.state = g.infoReturn
		}
		return nil
	}

	if justPressed(ebiten.KeyO) {
		g.openInfo("about")
		return nil
	}
	if justPressed(ebiten.KeyH) {
		g.openInfo("help")
		return nil
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
		if justPressed(ebiten.KeyN) {
			g.startNewGame()
			return nil
		}
		if justPressed(ebiten.KeyEnter) ||
			justPressed(ebiten.KeySpace) {
			level := g.save.UnlockedLevel
			if level < 1 {
				level = 1
			}
			if level > maxLevel {
				level = maxLevel
			}
			g.startLevel(level)
		}
		return nil

	case StatePaused:
		if justPressed(ebiten.KeyN) {
			g.startNewGame()
			return nil
		}
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
			if g.world.Level < maxLevel {
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
		g.player.VY *= 0.94
	}

	wasFalling := g.player.VY >= 0
	starsBefore := g.world.CollectedCount()
	g.world.ResolvePlayer(&g.player)
	g.world.ApplyHazardEffects(&g.player)

	if g.world.HazardAt(g.player.X, g.player.Y,
		g.player.Width, g.player.Height) {
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
		screen.Fill(color.RGBA{R: 166, G: 218, B: 245, A: 255})
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
		pc, pt := groundColor, topColor
		switch p.Kind {
		case world.PlatformIce:
			pc = color.RGBA{R: 120, G: 205, B: 238, A: 255}
			pt = color.RGBA{R: 235, G: 250, B: 255, A: 255}
		case world.PlatformMoving:
			pc = color.RGBA{R: 130, G: 160, B: 180, A: 255}
			pt = color.RGBA{R: 255, G: 220, B: 100, A: 255}
		case world.PlatformCrumbling:
			pc = color.RGBA{R: 155, G: 105, B: 65, A: 255}
			pt = color.RGBA{R: 235, G: 165, B: 85, A: 255}
		case world.PlatformConveyor:
			pc = color.RGBA{R: 95, G: 105, B: 115, A: 255}
			pt = color.RGBA{R: 245, G: 185, B: 55, A: 255}
		}
		ebitenutil.DrawRect(screen, x, p.Y, p.W, p.H, pc)
		ebitenutil.DrawRect(screen, x, p.Y, p.W, 6, pt)
	}

	for _, h := range g.world.Hazards {
		x := h.X - g.camera.X
		if x < -h.W || x > screenWidth {
			continue
		}
		hc := color.RGBA{R: 35, G: 35, B: 45, A: 255}
		switch h.Kind {
		case world.HazardIceSpikes:
			hc = color.RGBA{R: 190, G: 240, B: 255, A: 255}
		case world.HazardThorns:
			hc = color.RGBA{R: 50, G: 150, B: 65, A: 255}
		case world.HazardLava:
			hc = color.RGBA{R: 245, G: 75, B: 25, A: 255}
		case world.HazardQuicksand:
			hc = color.RGBA{R: 210, G: 165, B: 70, A: 255}
		case world.HazardLaser:
			hc = color.RGBA{R: 245, G: 35, B: 55, A: 255}
		case world.HazardEnergy:
			hc = color.RGBA{R: 70, G: 210, B: 255, A: 255}
		case world.HazardFallingRock:
			hc = color.RGBA{R: 105, G: 91, B: 85, A: 255}
		case world.HazardWind:
			hc = color.RGBA{R: 190, G: 230, B: 255, A: 255}
		}
		ebitenutil.DrawRect(screen, x, h.Y, h.W, h.H, hc)
	}

	for _, e := range g.world.Enemies {
		if !e.Active {
			continue
		}
		x := e.X - g.camera.X
		if x < -e.Width || x > screenWidth {
			continue
		}
		body := color.RGBA{R: 190, G: 55, B: 65, A: 255}
		head := color.RGBA{R: 220, G: 75, B: 80, A: 255}
		switch e.Kind {
		case world.EnemyRunner:
			body = color.RGBA{R: 215, G: 75, B: 45, A: 255}
			head = color.RGBA{R: 250, G: 115, B: 55, A: 255}
		case world.EnemyFlyer:
			body = color.RGBA{R: 45, G: 125, B: 190, A: 255}
			head = color.RGBA{R: 75, G: 185, B: 235, A: 255}
		case world.EnemyArmored:
			body = color.RGBA{R: 85, G: 95, B: 105, A: 255}
			head = color.RGBA{R: 155, G: 165, B: 175, A: 255}
		}
		if e.Boss {
			body = color.RGBA{R: 90, G: 35, B: 145, A: 255}
			head = color.RGBA{R: 180, G: 105, B: 255, A: 255}
		}
		ebitenutil.DrawRect(screen, x, e.Y+7, e.Width, e.Height-7, body)
		ebitenutil.DrawCircle(screen, x+e.Width/2, e.Y+9, e.Width/2, head)
		if e.Boss {
			ebitenutil.DebugPrintAt(screen, "Boss HP: "+itoa(e.Health), int(x), int(e.Y)-18)
		}
	}

	for _, c := range g.world.Collectibles {
		if !c.Taken {
			ebitenutil.DrawCircle(screen, c.X-g.camera.X, c.Y, 9,
				color.RGBA{R: 255, G: 194, B: 40, A: 255})
		}
	}

	if len(g.world.Platforms) > 0 {
		last := g.world.Platforms[len(g.world.Platforms)-1]
		goalX := g.world.EndX - g.camera.X
		cupWidth := 34.0
		cupTop := last.Y - 58
		if g.world.Level == maxLevel {
			cupWidth = 48
			cupTop = last.Y - 72
		}
		gold := color.RGBA{R: 255, G: 190, B: 35, A: 255}
		light := color.RGBA{R: 255, G: 225, B: 95, A: 255}
		ebitenutil.DrawRect(screen, goalX+cupWidth*0.28, cupTop+24, cupWidth*0.44, 18, gold)
		ebitenutil.DrawRect(screen, goalX+cupWidth*0.40, cupTop+42, cupWidth*0.20, 12, light)
		ebitenutil.DrawRect(screen, goalX+cupWidth*0.24, cupTop+54, cupWidth*0.52, 6, gold)
		ebitenutil.DrawCircle(screen, goalX+cupWidth/2, cupTop+18, cupWidth*0.40, gold)
		if g.world.Level == maxLevel && g.state == StatePlaying {
			ebitenutil.DebugPrintAt(screen, "FINAL TROPHY", int(goalX-12), int(last.Y-92))
		}
	}

	graphics.DrawPenguin(screen, float32(g.player.X-g.camera.X),
		float32(g.player.Y), &g.animation)

	ebitenutil.DebugPrintAt(screen,
		"HELWAN PENGUIN | WORLD "+itoa(g.world.Number), 18, 16)
	ebitenutil.DebugPrintAt(screen,
		"Level: "+itoa(g.world.Level)+"/35  Attempts: "+
			itoa(g.attempts)+"/"+itoa(maxAttempts)+
			"  Stars: "+itoa(g.world.CollectedCount())+
			"  Score: "+itoa(g.score), 18, 38)

	if g.saveError != "" {
		ebitenutil.DebugPrintAt(screen, g.saveError, 18, 82)
	}

	if g.infoPage != "" {
		g.drawInfo(screen)
		return
	}

	switch g.state {
	case StateMenu:
		drawOverlay(screen, "HELWAN PENGUIN",
			"ENTER: Continue | N: New Game",
			"O: About | H: Help | ESC: Quit")
		ebitenutil.DebugPrintAt(screen,
			"Unlocked: "+itoa(g.save.UnlockedLevel)+
				"  Total stars: "+itoa(g.save.TotalStars), 340, 330)
	case StatePaused:
		drawOverlay(screen, "PAUSED",
			"P or ENTER: Resume | N: New Game",
			"O: About | H: Help | ESC: Resume")
	case StateLevelComplete:
		drawOverlay(screen, "LEVEL COMPLETE!", "Press ENTER to continue", "")
	case StateGameOver:
		drawOverlay(screen, "NO ATTEMPTS LEFT",
			"Press ENTER to retry this level",
			"A new game gives you 5 attempts")
	case StateGameWon:
		drawOverlay(screen, "CONGRATULATIONS!",
			"YOU FINISHED HELWAN PENGUIN!",
			"All 35 levels complete - Press ENTER")
	default:
		ebitenutil.DebugPrintAt(screen,
			"Move: A/D or arrows | Jump: Space/Up | Pause: P", 18, 60)
	}
}

func (g *Game) drawInfo(screen *ebiten.Image) {
	ebitenutil.DrawRect(screen, 120, 100, 720, 340,
		color.RGBA{R: 20, G: 35, B: 50, A: 245})
	if g.infoPage == "about" {
		ebitenutil.DebugPrintAt(screen, "ABOUT HELWAN PENGUIN", 350, 140)
		ebitenutil.DebugPrintAt(screen, "A platform game by Helwan Linux", 285, 195)
		ebitenutil.DebugPrintAt(screen, "Arch Power, Helwan Simplicity", 300, 230)
		ebitenutil.DebugPrintAt(screen, "One Human Family, One Open Future.", 270, 265)
	} else {
		ebitenutil.DebugPrintAt(screen, "HOW TO PLAY", 400, 130)
		ebitenutil.DebugPrintAt(screen, "Move: A/D or Left/Right", 300, 170)
		ebitenutil.DebugPrintAt(screen, "Jump: Space or Up Arrow", 300, 200)
		ebitenutil.DebugPrintAt(screen, "Jump on enemies from above to defeat them", 245, 230)
		ebitenutil.DebugPrintAt(screen, "Collect stars and reach the trophy", 275, 260)
		ebitenutil.DebugPrintAt(screen, "P: Pause | N: New Game while paused", 280, 290)
		ebitenutil.DebugPrintAt(screen, "O: About | H: Help", 350, 320)
	}
	ebitenutil.DebugPrintAt(screen, "Press ESC or ENTER to return", 330, 390)
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
	case 7:
		return color.RGBA{R: 70, G: 35, B: 55, A: 255},
			color.RGBA{R: 255, G: 90, B: 130, A: 255}
	default:
		return color.RGBA{R: 116, G: 173, B: 194, A: 255},
			color.RGBA{R: 245, G: 252, B: 255, A: 255}
	}
}

func drawOverlay(screen *ebiten.Image, title, subtitle, footer string) {
	ebitenutil.DrawRect(screen, 180, 175, 600, 190,
		color.RGBA{R: 20, G: 35, B: 50, A: 235})
	ebitenutil.DebugPrintAt(screen, title, 390, 220)
	ebitenutil.DebugPrintAt(screen, subtitle, 260, 260)
	if footer != "" {
		ebitenutil.DebugPrintAt(screen, footer, 240, 295)
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
