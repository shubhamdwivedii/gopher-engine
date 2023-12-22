package gopher

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/math/f64"

	scr "github.com/shubhamdwivedii/gopher-engine/scene/screen"
)

type Gopher struct {
	Img *ebiten.Image
	X   float64
	Y   float64
	CX  float64
	CY  float64
	W   int
	H   int
	V   float64
	OP  *ebiten.DrawImageOptions
}

func New(cx, cy, v float64) *Gopher {
	img, _, err := ebitenutil.NewImageFromFile("./examples/assets/gopher.png")
	if err != nil {
		log.Fatal(err)
	}

	w, h := img.Bounds().Dx(), img.Bounds().Dy()
	x, y := cx-float64(w/2), cy-float64(h/2)

	return &Gopher{
		Img: img,
		X:   x,
		Y:   y,
		CX:  cx,
		CY:  cy,
		W:   w,
		H:   h,
		V:   v,
		OP:  &ebiten.DrawImageOptions{},
	}
}

func (g *Gopher) GetPosition() f64.Vec2 {
	return f64.Vec2{g.CX, g.CY}
}

func (g *Gopher) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.X -= g.V
		g.CX -= g.V
	}

	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.X += g.V
		g.CX += g.V
	}

	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		g.Y += g.V
		g.CY += g.V
	}

	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		g.Y -= g.V
		g.CY -= g.V
	}
	return nil
}

func (g *Gopher) Draw(gameScreen scr.Screen) {
	g.OP.GeoM.Reset()
	g.OP.GeoM.Translate(g.X, g.Y)
	gameScreen.DrawRect(g.X, g.Y, float64(g.W), float64(g.H), false, color.RGBA{255, 0, 0, 64})
	gameScreen.DrawImage(g.Img, g.OP)
}
