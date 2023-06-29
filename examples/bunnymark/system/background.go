package system

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/elemir/herd/examples/bunnymark/assets"
)

func Background(screen *ebiten.Image) {
	screen.Fill(assets.Background)
}
