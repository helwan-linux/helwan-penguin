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

	for _, p := range g.world.Platforms {
		graphics.DrawPlatform(screen, p, g.camera.X, g.world.Number)
	}

	for _, h := range g.world.Hazards {
		graphics.DrawHazard(screen, h, g.camera.X)
	}

	for _, e := range g.world.Enemies {
		graphics.DrawEnemy(screen, e, g.camera.X)
	}

	for _, c := range g.world.Collectibles {
		graphics.DrawCollectible(screen, c, g.camera.X)
	}

	graphics.DrawGoal(screen, g.world, g.camera.X, maxLevel)

	graphics.DrawPenguin(screen,
		float32(g.player.X-g.camera.X),
		float32(g.player.Y),
		&g.animation)

	graphics.DrawHUD(screen, g.world.Number, g.world.Level,
		g.attempts, maxAttempts, g.world.CollectedCount(),
		g.score, g.saveError)

	if g.infoPage != "" {
		graphics.DrawInfo(screen, g.infoPage)
		return
	}

	switch g.state {
	case StateMenu:
		graphics.DrawOverlay(screen, "HELWAN PENGUIN",
			"ENTER: Continue | N: New Game",
			"O: About | H: Help | ESC: Quit")
		ebitenutil.DebugPrintAt(screen,
			"Unlocked: "+itoa(g.save.UnlockedLevel)+
				"  Total stars: "+itoa(g.save.TotalStars), 340, 330)

	case StatePaused:
		graphics.DrawOverlay(screen, "PAUSED",
			"P or ENTER: Resume | N: New Game",
			"O: About | H: Help | ESC: Resume")

	case StateLevelComplete:
		graphics.DrawOverlay(screen, "LEVEL COMPLETE!",
			"Press ENTER to continue", "")

	case StateGameOver:
		graphics.DrawOverlay(screen, "NO ATTEMPTS LEFT",
			"Press ENTER to retry this level",
			"A new game gives you 5 attempts")

	case StateGameWon:
		graphics.DrawOverlay(screen, "CONGRATULATIONS!",
			"YOU FINISHED HELWAN PENGUIN!",
			"All 35 levels complete - Press ENTER")

	default:
		graphics.DrawControls(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
