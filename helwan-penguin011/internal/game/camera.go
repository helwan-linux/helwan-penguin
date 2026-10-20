package game

const cameraDeadZone = 360.0

type Camera struct{ X float64 }

func (c *Camera) Follow(playerX, playerY float64) {
	_ = playerY
	target := playerX - cameraDeadZone
	if target < 0 {
		target = 0
	}
	c.X = target
}
