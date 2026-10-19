
package world

import (
	"github.com/helwan-linux/helwan-penguin/internal/enemies"
	"github.com/helwan-linux/helwan-penguin/internal/player"
)

type World struct {
	Number       int
	Level        int
	Platforms    []Platform
	Collectibles []Collectible
	Enemies      []enemies.Enemy
	EndX         float64
}

func NewLevel(level int) World {
	if level < 1 {
		level = 1
	}
	if level > 30 {
		level = 30
	}

	worldNumber := (level-1)/5 + 1
	levelInWorld := (level-1)%5 + 1

	platforms := buildPlatforms(worldNumber, levelInWorld)
	collectibles := make([]Collectible, 0, len(platforms))
	enemyList := make([]enemies.Enemy, 0, len(platforms))

	for i, p := range platforms {
		if i > 0 && i < len(platforms)-1 {
			collectibles = append(collectibles, Collectible{
				X: p.X + p.W*0.5,
				Y: p.Y - 30 - float64((i+levelInWorld)%3)*8,
			})
		}

		if i == 0 || i == len(platforms)-1 {
			continue
		}

		enemyChance := (i + levelInWorld + worldNumber) % 4
		wideEnough := p.W >= 125

		if enemyChance == 1 && wideEnough {
			speed := 0.8 +
				float64(worldNumber)*0.2 +
				float64(levelInWorld)*0.08

			if (i+levelInWorld+worldNumber)%2 == 0 {
				speed = -speed
			}

			enemyList = append(enemyList, enemies.Enemy{
				X:      p.X + p.W/2 - 12,
				Y:      p.Y - 28,
				Width:  24,
				Height: 28,
				Active: true,
				VX:     speed,
				MinX:   p.X + 5,
				MaxX:   p.X + p.W - 5,
			})
		}
	}

	last := platforms[len(platforms)-1]

	if level%5 == 0 {
		enemyList = append(enemyList, enemies.Enemy{
			X:      last.X + last.W - 145,
			Y:      last.Y - 65,
			Width:  75,
			Height: 65,
			Active: true,
			VX:     1.2 + float64(worldNumber)*0.15,
			MinX:   last.X + last.W - 230,
			MaxX:   last.X + last.W - 10,
			Boss:   true,
			Health: 3 + worldNumber/2,
		})
	}

	return World{
		Number:       worldNumber,
		Level:        level,
		Platforms:    platforms,
		Collectibles: collectibles,
		Enemies:      enemyList,
		EndX:         last.X + last.W - 55,
	}
}

func buildPlatforms(worldNumber, levelInWorld int) []Platform {
	patterns := [][]int{
		{410, 385, 405, 370, 395, 355, 390, 365, 405, 375, 395, 410},
		{400, 365, 385, 345, 375, 330, 365, 350, 390, 340, 375, 400},
		{410, 385, 350, 375, 330, 365, 315, 350, 385, 340, 370, 405},
		{385, 350, 395, 360, 330, 380, 345, 400, 355, 325, 370, 405},
		{405, 365, 335, 385, 345, 315, 370, 330, 395, 350, 380, 405},
		{370, 335, 305, 350, 315, 285, 340, 300, 365, 325, 355, 395},
	}

	if worldNumber < 1 {
		worldNumber = 1
	}
	if worldNumber > len(patterns) {
		worldNumber = len(patterns)
	}

	pattern := patterns[worldNumber-1]

	platforms := []Platform{
		{
			X: 0,
			Y: 438,
			W: 330,
			H: 102,
		},
	}

	count := 9 + (worldNumber % 3) + ((levelInWorld + 1) % 3)

	x := 390.0
	previousY := 438.0

	for i := 0; i < count; i++ {
		step := i + (levelInWorld-1)*3

		y := float64(pattern[step%len(pattern)])

		// كل مستوى يغيّر مسار المنصات بدل تكرار نفس الترتيب.
		levelShift := (levelInWorld - 1) * 5
		if (i+levelShift)%4 == 2 {
			y += 12
		}
		if (i+worldNumber+levelInWorld)%5 == 0 {
			y -= 10
		}

		// منع القفزات الرأسية المبالغ فيها بين منصتين متتاليتين.
		if y < previousY-48 {
			y = previousY - 48
		}
		if y > previousY+42 {
			y = previousY + 42
		}

		if y < 285 {
			y = 285
		}
		if y > 425 {
			y = 425
		}

		width := 150.0 - float64((i*3+levelInWorld+worldNumber)%4)*13

		// منصات واسعة تمنح اللاعب فرصة للراحة وتغيير الإيقاع.
		if (i+levelInWorld)%4 == 0 {
			width += 35
		}

		// منصات ضيقة في مواضع مختلفة لزيادة التحدي تدريجيًا.
		if (i+levelInWorld+worldNumber)%5 == 2 {
			width -= 18
		}

		if levelInWorld >= 3 && (i+worldNumber)%4 == 1 {
			width -= 10
		}

		if width < 112 {
			width = 112
		}
		if width > 190 {
			width = 190
		}

		gap := 52.0 + float64((i*2+levelInWorld+worldNumber)%4)*9

		// المستويات المتقدمة تزيد المسافة بين بعض المنصات.
		if levelInWorld >= 3 && (i+worldNumber)%3 == 0 {
			gap += 7
		}

		// تخفيف المسافة بعد بعض المنصات الضيقة.
		if (i+levelInWorld)%4 == 3 {
			gap -= 8
		}

		if gap < 48 {
			gap = 48
		}
		if gap > 85 {
			gap = 85
		}

		// المنصة الأخيرة عريضة وثابتة لتكون ساحة النهاية أو الزعيم.
		if i == count-1 {
			width = 300
			y = 438
			gap = 0
		}

		platforms = append(platforms, Platform{
			X: x,
			Y: y,
			W: width,
			H: 540 - y,
		})

		x += width + gap
		previousY = y
	}

	return platforms
}

func FirstLevel() World {
	return NewLevel(1)
}

func (w *World) ResolvePlayer(p *player.Player) {
	p.OnGround = false
	previousY := p.Y - p.VY

	for _, platform := range w.Platforms {
		overlapsX := p.X+p.Width > platform.X &&
			p.X < platform.X+platform.W

		previousBottom := previousY + p.Height
		currentBottom := p.Y + p.Height

		if overlapsX &&
			p.VY >= 0 &&
			previousBottom <= platform.Y &&
			currentBottom >= platform.Y {
			p.Y = platform.Y - p.Height
			p.VY = 0
			p.OnGround = true
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
		if e.Boss && e.Active {
			return false
		}
	}

	return true
}
