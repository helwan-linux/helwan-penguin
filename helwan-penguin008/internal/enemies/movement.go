package enemies

func (e *Enemy) Update() {
	if !e.Active {
		return
	}

	if e.VX == 0 {
		e.VX = 1
	}

	e.X += e.VX

	if e.X < e.MinX {
		e.X = e.MinX
		if e.VX < 0 {
			e.VX = -e.VX
		}
	}

	if e.X+e.Width > e.MaxX {
		e.X = e.MaxX - e.Width
		if e.VX > 0 {
			e.VX = -e.VX
		}
	}
}
