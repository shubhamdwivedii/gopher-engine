package main

import (
	"fmt"
	"image/color"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	gop "github.com/shubhamdwivedii/gopher-engine/examples/gopher"
	ovr "github.com/shubhamdwivedii/gopher-engine/scene/overlay"
	scr "github.com/shubhamdwivedii/gopher-engine/scene/screen"
	vpt "github.com/shubhamdwivedii/gopher-engine/scene/viewport"
)

type Game struct{}

const (
	WORLD_W, WORLD_H = 420, 420
	VIEW_W, VIEW_H   = 320, 240
)

var gameScreen scr.Screen
var overlayScreen ovr.Overlay
var viewport *vpt.Viewport

var gopher *gop.Gopher
var healthbars *ebiten.Image
var worldbg *ebiten.Image
var gophers *ebiten.Image

func init() {
	var err error

	healthbars, _, err = ebitenutil.NewImageFromFile("./examples/assets/overlay_320x240.png")
	if err != nil {
		log.Fatal(err)
	}

	worldbg, _, err = ebitenutil.NewImageFromFile("./examples/assets/world_420x420.png")
	if err != nil {
		log.Fatal(err)
	}

	gophers, _, err = ebitenutil.NewImageFromFile("./examples/assets/gopher.png")
	if err != nil {
		log.Fatal(err)
	}

	gopher = gop.New(0, 0, 7)
	viewport = vpt.New(VIEW_W, VIEW_H, WORLD_W, WORLD_H, 160, 120)

	gameScreen, err = scr.New(VIEW_W, VIEW_H, WORLD_W, WORLD_H, viewport)

	if err != nil {
		log.Fatal(err)
	}
	overlayScreen = ovr.New(VIEW_W, VIEW_H)
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		gameScreen.GetShaker().Shake()
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) {
		viewport.MoveBy(-1, 0)
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		viewport.MoveBy(1, 0)
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		viewport.MoveBy(0, -1)
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		viewport.MoveBy(0, 1)
	}

	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		viewport.ZoomBy(-1)
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		viewport.ZoomBy(1)
	}

	if ebiten.IsKeyPressed(ebiten.KeyR) {
		viewport.RoatateBy(1)
	}

	if ebiten.IsKeyPressed(ebiten.KeyZ) {
		viewport.Reset()
	}

	gopher.Update()
	fmt.Println("GOPHER POSISION", gopher.CX, gopher.CY)
	// Update Camera After FocusEntity has been updated. (Or else you'll see jitter)
	gopherPos := gopher.GetPosition()
	viewport.Camera.Update(gopherPos[0], gopherPos[1])
	gameScreen.Update()
	return nil
}

func (g *Game) Draw(renderScreen *ebiten.Image) {
	// Draw to game screen first
	gameScreen.Fill(color.RGBA{202, 244, 244, 0xff})
	world := &ebiten.DrawImageOptions{}
	world.GeoM.Scale(float64(1), float64(1))
	gameScreen.DrawImage(worldbg, world)
	// gameScreen.GetImage().DrawImage(worldbg, &ebiten.DrawImageOptions{})

	gopher.Draw(gameScreen)

	g2OP := &ebiten.DrawImageOptions{}
	g3OP := &ebiten.DrawImageOptions{}
	g4OP := &ebiten.DrawImageOptions{}
	g1OP := &ebiten.DrawImageOptions{}

	g1OP.GeoM.Reset()
	g1OP.GeoM.Translate(0, 0)
	g2OP.GeoM.Reset()
	g2OP.GeoM.Translate(0, 420)
	g3OP.GeoM.Reset()
	g3OP.GeoM.Translate(420, 0)
	g4OP.GeoM.Reset()
	g4OP.GeoM.Translate(420, 420)

	gameScreen.DrawImage(gophers, g1OP)
	gameScreen.DrawImage(gophers, g2OP)
	gameScreen.DrawImage(gophers, g3OP)
	gameScreen.DrawImage(gophers, g4OP)

	gameScreen.Render(renderScreen)

	// Render Overlay Over the GameScreen
	healthbarsOP := &ebiten.DrawImageOptions{}
	// Transparency Doesn't work without this
	healthbarsOP.CompositeMode = ebiten.CompositeModeCopy
	healthbarsOP.ColorM.Scale(1, 1, 1, 0.25)
	overlayScreen.DrawImage(healthbars, healthbarsOP)
	overlayScreen.Render(renderScreen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	// return 1024, 768 // To Test Resolution Independent Scaling
	return VIEW_W, VIEW_H // Ideally Return Internal Resolution Here.
	// return 1024, 768 // Scaling of OVerlay broken
}

func main() {
	ebiten.SetWindowSize(640, 480)
	gameScreen.GetShaker().SetShakeIntensity(7.5)
	gameScreen.SetDebug(true)
	// gameScreen.GetViewport().SetMargin(10)
	viewport.AllowOutOfBounds = false
	viewport.Camera.OverflowAllowed(viewport.AllowOutOfBounds)
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
