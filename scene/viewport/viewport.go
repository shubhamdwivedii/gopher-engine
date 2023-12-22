package viewport

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	. "github.com/shubhamdwivedii/gopher-engine/constants"
	cam "github.com/shubhamdwivedii/gopher-engine/scene/viewport/camera"
	"golang.org/x/image/math/f64"
)

type Viewport struct {
	WorldSize        f64.Vec2 // Dimensions of Actual World (Viewport is just a portion of World)
	ViewSize         f64.Vec2 // Dimensions of the Viewport ie: Screen/Window
	Position         f64.Vec2 // TopLeft of Viewport (Not Middle/Center) w.r.t. WorldSize
	InitialPosition  f64.Vec2 // Initial Position of the Viewport
	WorldCenter      f64.Vec2 // Centre Coordinates of the World
	Margin           float64  // Will maintain a Margin between edge of Viewport and edge of World
	ZoomFactor       int      // Used to Zoom in and out of World
	Rotation         int      // Used to Rotate the Viewport
	AllowOutOfBounds bool     // Viewport can go outside of the World
	Camera           cam.Camera
}

func New(screenWidth, screenHeight, worldWidth, worldHeight int, centreX, centreY float64) *Viewport {
	posX, posY := centreX-float64(screenWidth)/2, centreY-float64(screenHeight)/2
	viewport := &Viewport{
		ViewSize:        f64.Vec2{float64(screenWidth), float64(screenHeight)},
		WorldSize:       f64.Vec2{float64(worldWidth), float64(worldHeight)},
		WorldCenter:     f64.Vec2{float64(worldWidth) / 2, float64(worldHeight) / 2},
		Position:        f64.Vec2{posX, posY},
		InitialPosition: f64.Vec2{posX, posY},
	}
	viewport.Camera = cam.New(worldWidth, worldHeight, 60, 60, 160, 120, viewport.MoveBy)
	// Rest are zero valued
	return viewport
}

func (v *Viewport) SetMargin(margin float64) {
	v.Margin = margin
}

// Get Center point Of Viewport in the World
func (v *Viewport) GetCenter() (cx, cy float64) {
	return v.Position[0] + v.ViewSize[0]/2, v.Position[1] + v.ViewSize[1]/2
}

func (v *Viewport) String() string {
	return fmt.Sprintf(
		"T: %.1f, R: %d, S: %d",
		v.Position, v.Rotation, v.ZoomFactor,
	)
}

// Checks if Viewport is OutOfBounds (Outside WorldView)
// Retuns dx, dy to adjust In-Bound
func (v *Viewport) OutOfBounds() (dx, dy float64) {
	x1, y1 := v.Position[0], v.Position[1]
	x2, y2 := x1+v.ViewSize[0], y1+v.ViewSize[1]

	if x1 < v.Margin {
		dx = v.Margin - x1
	}

	if x2 > v.WorldSize[0]-v.Margin {
		dx = v.WorldSize[0] - v.Margin - x2
	}

	if y1 < v.Margin {
		dy = v.Margin - y1
	}

	if y2 > v.WorldSize[1]-v.Margin {
		dy = v.WorldSize[1] - v.Margin - y2
	}

	if dx > 0 || dy > 0 {
		fmt.Println("OUT OF BOUNDS")
	}

	return dx, dy
}

func (v *Viewport) GetPosition() f64.Vec2 {
	return v.Position // TopLeft Coordinates in World
}

func (v *Viewport) GetSize() f64.Vec2 {
	return v.ViewSize
}

func (v *Viewport) GetScale() f64.Vec2 {
	return f64.Vec2{float64(v.ZoomFactor), float64(v.ZoomFactor)}
}
func (v *Viewport) GetRotation() float64 {
	return float64(v.Rotation)
}

func (v *Viewport) GetMatrix() ebiten.GeoM {
	matrix := ebiten.GeoM{}
	// matrix.Translate(-v.Position[0], -v.Position[1])

	/* IMPORTANT :-
	Since Viewport's position(x,y) is always supposed to be at top-left of Render-Window,
	We need to translate the whole world by -vpX, -vpY,
	where (vpX,vpY) is position of the Viewport in the World */
	cx, cy := v.GetCenter()
	matrix.Translate(-cx, -cy)

	// FIX Scaling
	if v.ZoomFactor != 0 {
		matrix.Scale(
			math.Pow(1.01, float64(v.ZoomFactor)),
			math.Pow(1.01, float64(v.ZoomFactor)),
		)
	} else {
		matrix.Scale(
			math.Pow(1.0, float64(v.ZoomFactor)),
			math.Pow(1.0, float64(v.ZoomFactor)),
		)
	}

	matrix.Rotate(float64(v.Rotation) * 2 * math.Pi / 360)
	matrix.Translate(cx, cy)

	/* NOTE :-
	This is not the matrix of just the Viewport, but the matrix of entire World,
	Any Scaling, Rotations done to Viewport also need to be done to the World
	*/

	return matrix
}

/*
Converts Screen Coordinates to World Coordinates
Can be used when you want OVERLAY elements on Screen that stays on the Viewport
Example: UI elements, Text, Health Bar, FPS etc.
*/
func (v *Viewport) ScreenToWorld(posX, posY int) f64.Vec2 {
	inverseMatrix := v.GetMatrix()
	if inverseMatrix.IsInvertible() {
		inverseMatrix.Invert()
		wPosX, wPosY := inverseMatrix.Apply(float64(posX), float64(posY))
		return f64.Vec2{wPosX, wPosY}
	}
	// when scaling its possible that matrix is not invertible
	return f64.Vec2{math.NaN(), math.NaN()}
}

func (v *Viewport) Reset() {
	v.Position[0] = v.InitialPosition[0]
	v.Position[1] = v.InitialPosition[1]
	v.Rotation = 0
	v.ZoomFactor = 0
}

// (0,0) is default origin
func (v *Viewport) MoveTo(x, y float64) {
	v.Position[0] = x
	v.Position[1] = y

	if !v.AllowOutOfBounds {
		dx, dy := v.OutOfBounds()
		v.Position[0] += dx
		v.Position[1] += dy
	}
}

func (v *Viewport) MoveBy(dx, dy float64) {
	v.Position[0] += dx
	v.Position[1] += dy

	if !v.AllowOutOfBounds {
		odx, ody := v.OutOfBounds()
		v.Position[0] += odx
		v.Position[1] += ody
	}
}

func (v *Viewport) ZoomBy(dz int) {
	if (v.ZoomFactor+dz) > MIN_VIEWPORT_ZOOM && (v.ZoomFactor+dz) < MAX_VIEWPORT_ZOOM {
		v.ZoomFactor += dz
	}
}

// default z = 0
func (v *Viewport) SetZoom(z int) {
	if z > MIN_VIEWPORT_ZOOM && z < MAX_VIEWPORT_ZOOM {
		v.ZoomFactor = z
	}
}

// default r = 0
func (v *Viewport) SetRotation(r int) {
	v.Rotation = r
}

func (v *Viewport) RoatateBy(dr int) {
	v.Rotation += dr
}

// offset is (0,0) when camera position is (ww/2, wh/2)
func (v *Viewport) GetOffsets() (float64, float64) {
	// Right +ve, Left -ve, Up -ve, Down +ve

	cx, cy := v.Position[0]+v.ViewSize[0]/2, v.Position[1]+v.ViewSize[1]/2 // Center point of the Viewport
	// dx, dy := cx-v.WorldSize[0]/2, cy-v.WorldSize[1]/2
	dx, dy := cx-v.ViewSize[0]/2, cy-v.ViewSize[1]/2

	// if camera moves towards right ??? Why this works ?
	return -dx, -dy

	// return -v.Position[0], -v.Position[1]

}

// Concat this matrix for adjust for camera position,
func (v *Viewport) GetOffsetMatrix() ebiten.GeoM {
	matrix := ebiten.GeoM{}
	dx, dy := v.GetOffsets()
	matrix.Translate(dx, dy)
	return matrix
}

// Both GetOffsets and GetOffsetMatrix are used by the respective functions in Screen
// See Screen's GetOffsets and GetOffsetMatrix to understand how these are used
