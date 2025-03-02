package main

import (
	"fmt"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	EmptyCell = 0
	RowSize   = 9
	CellSize  = 64
	GapSize   = 2
	OffsetX   = 200
	OffsetY   = 200
)

type Cell struct {
	ItemId int
	Amount uint
	Hitbox image.Rectangle
}

type Inventory struct {
	Cells []*Cell
	Hand  Cell
}

func GetHitbox(index int) image.Rectangle {
	gridX := index % RowSize
	gridY := index / RowSize
	pixelX := gridX*(CellSize+GapSize) + OffsetX
	pixelY := gridY*(CellSize+GapSize) + OffsetY

	return image.Rect(pixelX, pixelY, pixelX+CellSize, pixelY+CellSize)
}

func NewInventory() *Inventory {
	cells := make([]*Cell, 27)
	for i := range cells {
		cells[i] = &Cell{
			Hitbox: GetHitbox(i),
		}
	}

	return &Inventory{
		Cells: cells,
		Hand:  Cell{},
	}
}

func (i *Inventory) Update() error {
	mouseX, mouseY := ebiten.CursorPosition()
	i.Hand.Hitbox = image.Rect(mouseX, mouseY, mouseX+CellSize, mouseY+CellSize)
	mouseClicked := inpututil.IsMouseButtonJustPressed((ebiten.MouseButton0))
	rightMouseClicked := inpututil.IsMouseButtonJustPressed((ebiten.MouseButton2))

	if mouseClicked {
		// check if clicked cell
		mouseP := image.Point{mouseX, mouseY}
		for index, cell := range i.Cells {
			if mouseP.In(cell.Hitbox) {
				i.PlaceItems(index)
			}
		}
	}

	if rightMouseClicked {
		// check if clicked cell
		mouseP := image.Point{mouseX, mouseY}
		for index, cell := range i.Cells {
			if mouseP.In(cell.Hitbox) {
				i.PlaceOneItem(index)
			}
		}
	}

	return nil
}

func (i *Inventory) DrawHand(screen *ebiten.Image, atlas *ebiten.Image) {
	if i.Hand.ItemId != EmptyCell {
		id := i.Hand.ItemId - 1
		numItemsPerRow := atlas.Bounds().Dx() / 8
		gridX := id % numItemsPerRow
		gridY := id / numItemsPerRow

		sourceRect := image.Rect(
			gridX*8.0,
			gridY*8.0,
			(gridX+1)*8.0,
			(gridY+1)*8.0,
		)

		options := ebiten.DrawImageOptions{}
		options.GeoM.Scale(8, 8)
		options.GeoM.Translate(float64(i.Hand.Hitbox.Min.X-(CellSize/2)), float64(i.Hand.Hitbox.Min.Y-(CellSize/2)))

		screen.DrawImage(
			atlas.SubImage(sourceRect).(*ebiten.Image),
			&options,
		)

		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d", i.Hand.Amount), i.Hand.Hitbox.Min.X-(CellSize/2), i.Hand.Hitbox.Min.Y-(CellSize/2))
	}
}

func (i *Inventory) Draw(screen *ebiten.Image, atlas *ebiten.Image) {

	for _, cell := range i.Cells {
		vector.DrawFilledRect(
			screen,
			float32(cell.Hitbox.Min.X),
			float32(cell.Hitbox.Min.Y),
			float32(cell.Hitbox.Dx()),
			float32(cell.Hitbox.Dy()),
			color.RGBA{20, 20, 20, 255},
			true,
		)

		if cell.ItemId != EmptyCell {
			id := cell.ItemId - 1
			numItemsPerRow := atlas.Bounds().Dx() / 8
			gridX := id % numItemsPerRow
			gridY := id / numItemsPerRow

			sourceRect := image.Rect(
				gridX*8.0,
				gridY*8.0,
				(gridX+1)*8.0,
				(gridY+1)*8.0,
			)

			options := ebiten.DrawImageOptions{}
			options.GeoM.Scale(8, 8)
			options.GeoM.Translate(float64(cell.Hitbox.Min.X), float64(cell.Hitbox.Min.Y))

			screen.DrawImage(
				atlas.SubImage(sourceRect).(*ebiten.Image),
				&options,
			)

			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d", cell.Amount), cell.Hitbox.Min.X, cell.Hitbox.Min.Y)
		}
	}

	i.DrawHand(screen, atlas)
}

func (i *Inventory) PrintInventory() {
	for _, cell := range i.Cells {
		fmt.Printf("Cell: %d %d\n", cell.ItemId, cell.Amount)
	}
	fmt.Printf("Hand: %d %d\n", i.Hand.ItemId, i.Hand.Amount)
}

func (i *Inventory) PlaceOneItem(index int) bool {
	// validate the index
	if index < 0 || index > len(i.Cells)-1 {
		return false
	}

	cell := i.Cells[index]

	// check if we can merge
	if cell.ItemId == i.Hand.ItemId {
		cell.Amount++
		i.Hand.Amount--
		if i.Hand.Amount < 1 {
			i.Hand.ItemId = EmptyCell
		}
		return true
	}

	// check if cell is empty
	if cell.ItemId == EmptyCell {
		cell.ItemId = i.Hand.ItemId
		cell.Amount++
		i.Hand.Amount--
		if i.Hand.Amount < 1 {
			i.Hand.ItemId = EmptyCell
		}
		return true
	}

	if i.Hand.ItemId == EmptyCell && cell.Amount > 1 {
		half := cell.Amount / 2
		i.Hand.ItemId = cell.ItemId
		i.Hand.Amount = half
		cell.Amount -= half
		return true
	}

	// they don't match and call isn't empty, so swap
	temp := cell.ItemId
	cell.ItemId = i.Hand.ItemId
	i.Hand.ItemId = temp

	tempAmount := cell.Amount
	cell.Amount = i.Hand.Amount
	i.Hand.Amount = tempAmount

	return true
}

func (i *Inventory) PlaceItems(index int) bool {
	// validate the index
	if index < 0 || index > len(i.Cells)-1 {
		return false
	}

	cell := i.Cells[index]

	// check if we can merge
	if cell.ItemId == i.Hand.ItemId {
		cell.Amount += i.Hand.Amount
		i.Hand.Amount = 0
		i.Hand.ItemId = EmptyCell
		return true
	}

	// check if cell is empty
	if cell.ItemId == EmptyCell {
		cell.ItemId = i.Hand.ItemId
		cell.Amount = i.Hand.Amount
		i.Hand.ItemId = EmptyCell
		i.Hand.Amount = 0
		return true
	}

	// they don't match and call isn't empty, so swap
	temp := cell.ItemId
	cell.ItemId = i.Hand.ItemId
	i.Hand.ItemId = temp

	tempAmount := cell.Amount
	cell.Amount = i.Hand.Amount
	i.Hand.Amount = tempAmount

	return true
}
