package utils

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/math/f64"
)

type entity interface {
	GetPosition() f64.Vec2 // TopLeft Coordinates in World
	GetSize() f64.Vec2     // Size of the entity
	GetScale() f64.Vec2
	GetRotation() float64
}

func Fill(image *ebiten.Image, col color.Color) {
	image.Fill(col)
}

func DrawImage(image *ebiten.Image, targetImage *ebiten.Image, op *ebiten.DrawImageOptions) {
	targetImage.DrawImage(image, op)
}

func DrawLine(image *ebiten.Image, x1, y1, x2, y2 float64, col color.Color) {
	ebitenutil.DrawLine(image, x1, y1, x2, y2, col)
}

func DrawRect(image *ebiten.Image, x, y, width, height float64, solid bool, clr color.Color) {
	if solid {
		ebitenutil.DrawRect(image, x, y, width, height, clr)
	} else {
		x2 := x + width
		y2 := y + height
		DrawLine(image, x, y, x2, y, clr)
		DrawLine(image, x+1, y, x+1, y2, clr)
		DrawLine(image, x2, y, x2, y2, clr)
		DrawLine(image, x, y2-1, x2, y2-1, clr)
	}
}

func DebugPrint(image *ebiten.Image, text string) {
	DebugPrintAt(image, text, 0, 0)
}

func DebugPrintAt(image *ebiten.Image, text string, x, y int) {
	ebitenutil.DebugPrintAt(image, text, x, y)
}

func DrawText(image *ebiten.Image, txt string, fnt font.Face, x, y int, clr color.Color) {
	text.Draw(image, txt, fnt, x, y, clr)
}
