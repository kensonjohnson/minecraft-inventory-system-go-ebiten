package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	Inventory  *Inventory
	AtlasImage *ebiten.Image
}

func NewGame() *Game {

	img := mustLoadImage("assets/atlas.png")

	inventory := NewInventory()
	inventory.Hand.ItemId = 1
	inventory.Hand.Amount = 10
	inventory.PlaceItems(1)
	inventory.Hand.ItemId = 2
	inventory.Hand.Amount = 20
	inventory.PlaceItems(4)

	return &Game{
		Inventory:  inventory,
		AtlasImage: img,
	}
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		return ebiten.Termination
	}
	g.Inventory.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{100, 149, 237, 255})
	g.Inventory.Draw(screen, g.AtlasImage)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("Inventory Example")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
