package main

import (
	"bytes"
	"embed"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

//go:embed assets
var assetsFS embed.FS

func mustLoadImage(filePath string) *ebiten.Image {
	imgSource, err := assetsFS.ReadFile(filePath)
	if err != nil {
		log.Panic(err)
	}
	image, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(imgSource))
	if err != nil {
		log.Panic(err)
	}
	return image
}
