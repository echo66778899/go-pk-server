package ui

import (
	"image"
	"log"
	"sync"

	tb "github.com/nsf/termbox-go"
)

type MainPlayer struct {
	CurrentPlayerPossition int
	sync.Mutex
}

var CurrentPlayer = MainPlayer{CurrentPlayerPossition: 0}

type Drawable interface {
	GetRect() image.Rectangle
	SetRect(int, int, int, int)
	Draw(*Buffer)
	sync.Locker
}

// Initialize termbox
func Init() {
	err := tb.Init()
	if err != nil {
		log.Fatalf("Failed to initialize termbox: %v", err)
	}
}

func Deinit() {
	tb.Close()
}

func Render(items ...Drawable) {
	tb.Clear(tb.ColorDefault, tb.ColorDefault)

	for _, item := range items {
		buf := NewBuffer(item.GetRect())
		item.Lock()
		item.Draw(buf)
		item.Unlock()
		for point, cell := range buf.CellMap {
			if point.In(buf.Rectangle) {
				tb.SetCell(
					point.X, point.Y,
					cell.Rune,
					tb.Attribute(cell.Style.Fg+1)|tb.Attribute(cell.Style.Modifier), tb.Attribute(cell.Style.Bg+1),
				)
			}
		}
	}
	tb.Flush()
}
