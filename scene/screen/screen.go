package screen

import (
	"errors"
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	. "github.com/shubhamdwivedii/gopher-engine/constants"
	shk "github.com/shubhamdwivedii/gopher-engine/scene/screen/shaker"
	vpt "github.com/shubhamdwivedii/gopher-engine/scene/viewport"
	"github.com/shubhamdwivedii/gopher-engine/utils"
	"golang.org/x/image/font"
	"golang.org/x/image/math/f64"
)

type Screen interface {
	SetDebug(debugOn bool)
	Update() error
	Render(target *ebiten.Image)
	GetImage() (screenImage *ebiten.Image)
	GetViewport() (viewport *vpt.Viewport)
	GetShaker() (shaker shk.ScreenShaker)

	DrawImage(image *ebiten.Image, op *ebiten.DrawImageOptions)
	DrawLine(x1, y1, x2, y2 float64, col color.Color)
	DrawRect(x, y, width, height float64, fill bool, col color.Color)
	Fill(col color.Color)
	DebugPrint(text string)
	DebugPrintAt(text string, x, y int)
	DrawText(text string, fnt font.Face, x, y int, clr color.Color)
}

type CustomScreen struct {
	ScreenSize     f64.Vec2
	WorldSize      f64.Vec2
	Image          *ebiten.Image
	Viewport       *vpt.Viewport
	Shaker         shk.ScreenShaker
	DrawOP         *ebiten.DrawImageOptions
	Debug          bool
	AutoScaling    bool
	AutoPadding    bool
	StaticViewport bool
}

func New(screenWidth, screenHeight, worldWidth, worldHeight int, viewport *vpt.Viewport) (Screen, error) {
	autoPadding := false
	if screenWidth == worldWidth && screenHeight == worldHeight {
		if viewport != nil {
			return nil, errors.New("viewport should be nil when screen-size equals world-size")
		}
		autoPadding = true
		worldWidth += AUTO_PADDING * 2
		worldHeight += AUTO_PADDING * 2
	} else {
		if viewport == nil {
			return nil, errors.New("viewport cannot be nil if screen-size and world-size are different")
		}
	}
	screenImg := ebiten.NewImage(worldWidth, worldHeight)

	shaker := shk.New()

	return &CustomScreen{
		Image:          screenImg,
		ScreenSize:     f64.Vec2{float64(screenWidth), float64(screenHeight)},
		WorldSize:      f64.Vec2{float64(worldWidth), float64(worldHeight)},
		Viewport:       viewport,
		Shaker:         shaker,
		DrawOP:         &ebiten.DrawImageOptions{},
		AutoScaling:    true,
		StaticViewport: viewport == nil,
		AutoPadding:    autoPadding,
	}, nil

}

func (s *CustomScreen) SetDebug(debugOn bool) {
	s.Debug = debugOn
}

func (s *CustomScreen) Update() error {
	return s.Shaker.Update()
}

// Draws CustomScreen to RenderScreen (target)

func (s *CustomScreen) Render(targetScreen *ebiten.Image) {
	s.DrawOP.GeoM.Reset()

	sdx, sdy := s.Shaker.GetOffsets()
	s.DrawOP.GeoM.Translate(-sdx, -sdy)

	if s.AutoPadding && s.Viewport == nil {
		// Need To Render CustomScreen slightly off left/top (on RenderScreen) to adjust for AutoPadding
		s.DrawOP.GeoM.Translate(-AUTO_PADDING, -AUTO_PADDING)
	} else {
		transformMatrix := s.Viewport.GetMatrix()
		s.DrawOP.GeoM.Concat(transformMatrix)
	}

	// Scaling Screen Image to Render Resolution
	if s.AutoScaling {
		resX, resY := targetScreen.Bounds().Dx(), targetScreen.Bounds().Dy()
		if resX != int(s.ScreenSize[0]) || resY != int(s.ScreenSize[1]) {
			scaleX, scaleY := float64(resX)/s.ScreenSize[0], float64(resY)/s.ScreenSize[1]
			s.DrawOP.GeoM.Scale(scaleX, scaleY)
		}
	}

	s.drawCameraFocusArea()

	// Render Screen Image to Real Render Screen
	targetScreen.DrawImage(s.Image, s.DrawOP)
}

func (s *CustomScreen) GetImage() *ebiten.Image {
	return s.Image
}

func (s *CustomScreen) GetViewport() *vpt.Viewport {
	return s.Viewport
}

func (s *CustomScreen) GetShaker() shk.ScreenShaker {
	return s.Shaker
}

// Includes Padding-Offset if Viewport is Nil
func (s *CustomScreen) GetOffsets() (dx, dy float64) {
	if s.Viewport != nil {
		dx, dy = s.Viewport.GetOffsets()
	}
	if s.AutoPadding && s.Viewport == nil {
		dx += AUTO_PADDING
		dy += AUTO_PADDING
	}
	return
}

func (s *CustomScreen) GetOffsetMatrix() (offsetMatrix ebiten.GeoM) {
	if s.Viewport != nil {
		offsetMatrix = s.Viewport.GetOffsetMatrix()
	}
	if s.AutoPadding && s.Viewport == nil {
		offsetMatrix.Translate(AUTO_PADDING, AUTO_PADDING)
	}
	return
}

// NOTE :--
// Offsets are used to render relative to screenOrigin (instead of worldOrigin)
// offx, offy := float64(worldWidth-screenWidth)/2, float64(worldHeight-screenHeight)/2

// offsetMatrix := ebiten.GeoM{}
// offsetMatrix.Translate(offx, offy)

/***************** DRAWING FUNCTIONS *********************/

// Takes coordinates based on Screen and Adjusts automatically for World (Screen x1,y1 are 0,0)
func (s *CustomScreen) DrawImage(image *ebiten.Image, op *ebiten.DrawImageOptions) {
	cameraMatrix := s.GetOffsetMatrix()
	op.GeoM.Concat(cameraMatrix)
	utils.DrawImage(image, s.Image, op)
}

func (s *CustomScreen) Fill(col color.Color) {
	utils.Fill(s.Image, col)
}

func (s *CustomScreen) DrawLine(x1, y1, x2, y2 float64, col color.Color) {
	offx, offy := s.GetOffsets()
	utils.DrawLine(s.Image, x1+offx, y1+offy, x2+offx, y2+offy, col)
}

func (s *CustomScreen) DrawRect(x, y, width, height float64, solid bool, clr color.Color) {
	offx, offy := s.GetOffsets()
	utils.DrawRect(s.Image, x+offx, y+offy, width, height, solid, clr)
}

func (s *CustomScreen) DebugPrint(text string) {
	utils.DebugPrint(s.Image, text)
}

func (s *CustomScreen) DebugPrintAt(text string, x, y int) {
	offx, offy := s.GetOffsets()
	utils.DebugPrintAt(s.Image, text, x+int(offx), y+int(offy))
}

func (s *CustomScreen) DrawText(txt string, fnt font.Face, x, y int, clr color.Color) {
	offx, offy := s.GetOffsets()
	utils.DrawText(s.Image, txt, fnt, x+int(offx), y+int(offy), clr)
}

// TRASH FUNC

func (s *CustomScreen) drawCameraFocusArea() {
	cPosition := s.Viewport.Position
	cFocusCenter, cFocusView := s.Viewport.Camera.GetFocus()
	s.DebugPrintAt(fmt.Sprintf("Viewport-TLX: %0.2f Viewport-TLY: %0.2f", cPosition[0], cPosition[1]), 0, 32)
	s.DebugPrintAt(fmt.Sprintf("Camera-CX: %0.2f Camera-CY: %0.2f", cFocusCenter[0], cFocusCenter[1]), 0, 64)

	x1 := cFocusCenter[0] - cFocusView[0]/2
	x2 := x1 + cFocusView[0]
	y1 := cFocusCenter[1] - cFocusView[1]/2
	y2 := y1 + cFocusView[1]

	// Camera Offset is adjusted in CustomScreen.DebugPrintAt()

	s.DrawLine(x1, y1, x2, y1, color.RGBA{0, 0, 255, 255})
	s.DrawLine(x1+1, y1, x1+1, y2, color.RGBA{0, 0, 255, 255})
	s.DrawLine(x2, y1, x2, y2, color.RGBA{0, 0, 255, 255})
	s.DrawLine(x1, y2-1, x2, y2-1, color.RGBA{0, 0, 255, 255})

	worldXY := s.Viewport.ScreenToWorld(ebiten.CursorPosition())

	utils.DebugPrintAt(
		s.Image,
		fmt.Sprintf("Cursor World Pos: %.2f,%.2f",
			worldXY[0], worldXY[1]),
		0, 92,
	)

	x, y := ebiten.CursorPosition()
	utils.DebugPrintAt(
		s.Image,
		fmt.Sprintf("Cursor Viewport Pos: %v,%v",
			x, y),
		0, 128,
	)

	utils.DebugPrintAt(
		s.Image,
		fmt.Sprintf("Focus Top Left: %.2f,%.2f",
			x1, y1),
		0, 164,
	)

	s.DebugPrintAt(
		fmt.Sprintf("Viewport Position: %.2f, %.2f",
			s.Viewport.Position[0], s.Viewport.Position[1]),
		0, 192,
	)

}
