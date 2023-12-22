package overlay

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/shubhamdwivedii/gopher-engine/utils"
	"golang.org/x/image/font"
	"golang.org/x/image/math/f64"
)

// Overlay is like a Static Screen over Game Screen, used for UI elements like health-bar etc
type Overlay interface {
	Render(screen *ebiten.Image)
	DrawImage(image *ebiten.Image, op *ebiten.DrawImageOptions)
	DrawLine(x1, y1, x2, y2 float64, col color.Color)
	DrawRect(x, y, width, height float64, fill bool, col color.Color)
	Fill(col color.Color)
	DebugPrint(text string)
	DebugPrintAt(text string, x, y int)
	DrawText(text string, fnt font.Face, x, y int, clr color.Color)
}

// Can be used for Overlay, Effects or Transisions
type StaticScreen struct {
	ScreenSize  f64.Vec2
	Image       *ebiten.Image
	DrawOP      *ebiten.DrawImageOptions
	Debug       bool
	AutoScaling bool
}

func New(width, height int) Overlay {
	screenImg := ebiten.NewImage(width, height)

	// Move to separate func debug
	screenImg.Fill(color.RGBA{64, 220, 14, 64})

	return &StaticScreen{
		Image:       screenImg,
		ScreenSize:  f64.Vec2{float64(width), float64(height)},
		DrawOP:      &ebiten.DrawImageOptions{},
		AutoScaling: true,
	}

}

// Renders Overlay on RenderScreen (target)
func (s *StaticScreen) Render(targetScreen *ebiten.Image) {
	s.DrawOP.GeoM.Reset()

	// Scaling Screen Image to Render Resolution
	if !s.AutoScaling {
		resX, resY := targetScreen.Bounds().Dx(), targetScreen.Bounds().Dy()
		if resX != int(s.ScreenSize[0]) || resY != int(s.ScreenSize[1]) {
			scaleX, scaleY := float64(resX)/float64(s.ScreenSize[0]), float64(resY)/float64(s.ScreenSize[1])
			s.DrawOP.GeoM.Scale(scaleX, scaleY)
		}
	}

	utils.DrawImage(s.Image, targetScreen, s.DrawOP)
	// screen.DrawImage(s.Image, s.DrawOP)
}

func (s *StaticScreen) Fill(col color.Color) {
	utils.Fill(s.Image, col)
}

func (s *StaticScreen) DrawImage(image *ebiten.Image, op *ebiten.DrawImageOptions) {
	utils.DrawImage(image, s.Image, op)
}

func (s *StaticScreen) DrawLine(x1, y1, x2, y2 float64, col color.Color) {
	utils.DrawLine(s.Image, x1, y1, x2, y2, col)
}

func (s *StaticScreen) DrawRect(x, y, width, height float64, solid bool, clr color.Color) {
	utils.DrawRect(s.Image, x, y, width, height, solid, clr)
}

func (s *StaticScreen) DebugPrint(text string) {
	utils.DebugPrint(s.Image, text)
}

func (s *StaticScreen) DebugPrintAt(text string, x, y int) {
	utils.DebugPrintAt(s.Image, text, x, y)
}

func (s *StaticScreen) DrawText(txt string, fnt font.Face, x, y int, clr color.Color) {
	utils.DrawText(s.Image, txt, fnt, x, y, clr)
}
