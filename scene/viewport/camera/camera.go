package camera

import (
	"golang.org/x/image/math/f64"
)

type Camera interface {
	Update(x, y float64) error
	OverflowAllowed(allowed bool)
	GetFocus() (focusPoint f64.Vec2, focusSize f64.Vec2)
}

type FocusBoxCam struct {
	WorldSize        f64.Vec2 // Dimensions of Actual World (Viewport is just a portion of World)
	FocusSize        f64.Vec2 // Dimensions of the FocusArea ie: Screen/Window
	FocusPoint       f64.Vec2 // Point of focus, initially ww/2, wh/2 (center point of FocusBox)
	PositionCallback func(dx, dy float64)
	Debug            bool
	AllowOutOfBounds bool
}

// worldWidth, worldHeight is the width/height of the World (including out-of-screen area)
// focusWidth, focusHeight is the width/height of focus area (used to check if x,y is out of focus)
// focusX, focusY is point of focus within the the World (center of FocusArea)
func New(worldWidth, worldHeight, focusWidth, focusHeight int, focusX, focusY float64, updatePosition func(dx, dy float64)) Camera {
	return &FocusBoxCam{
		WorldSize:        f64.Vec2{float64(worldWidth), float64(worldHeight)},
		FocusSize:        f64.Vec2{float64(focusWidth), float64(focusHeight)},
		FocusPoint:       f64.Vec2{focusX, focusY},
		PositionCallback: updatePosition,
	}
	// ORIGIN is (0,0), FocusedEntity is nil
}

func (c *FocusBoxCam) GetFocus() (focusPoint f64.Vec2, focusSize f64.Vec2) {
	return c.FocusPoint, c.FocusSize
}

func (c *FocusBoxCam) OverflowAllowed(allowed bool) {
	c.AllowOutOfBounds = allowed
}

func (c *FocusBoxCam) CheckWorldOverflow(x, y float64) (dx, dy float64) {
	x1, y1 := x-c.FocusSize[0]/2, y-c.FocusSize[1]/2
	x2, y2 := x+c.FocusSize[0]/2, y+c.FocusSize[1]/2

	if x1 < 0 {
		dx = -x1
	}

	if x2 > c.WorldSize[0] {
		dx = c.WorldSize[0] - x2
	}

	if y1 < 0 {
		dy = -y1
	}

	if y2 > c.WorldSize[1] {
		dy = c.WorldSize[1] - y2
	}

	return
}

func (c *FocusBoxCam) Update(x, y float64) error {
	dx, dy := 0.0, 0.0
	hw, hh := c.FocusSize[0]/2, c.FocusSize[1]/2     // half width, half height
	lx, rx := c.FocusPoint[0]-hw, c.FocusPoint[0]+hw // left x, right x
	ty, by := c.FocusPoint[1]-hh, c.FocusPoint[1]+hh // top y, bottom y

	if x < lx {
		dx = x - lx
	}

	if x > rx {
		dx = x - rx
	}

	if y < ty {
		dy = y - ty
	}

	if y > by {
		dy = y - by
	}

	if !c.AllowOutOfBounds {
		rdx, rdy := c.CheckWorldOverflow(c.FocusPoint[0]+dx, c.FocusPoint[1]+dy)
		c.FocusPoint[0] += dx + rdx
		c.FocusPoint[1] += dy + rdy
		c.PositionCallback(dx+rdx, dy+rdy)
		return nil
	}

	c.FocusPoint[0] += dx
	c.FocusPoint[1] += dy

	c.PositionCallback(dx, dy)

	return nil
}
