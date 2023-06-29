package system

import (
	"github.com/elemir/herd"
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/elemir/herd/examples/bunnymark/component"
)

func Render(screen *ebiten.Image, query herd.Query3[component.Position, component.Hue, component.Sprite]) {
	query.ForEach(func(_ herd.EntityID, pos *component.Position, hue *component.Hue, sprite *component.Sprite) {
		sw, sh := float64(screen.Bounds().Dx()), float64(screen.Bounds().Dy())
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(pos.X*sw, pos.Y*sh)
		if *hue.Colorful {
			op.ColorM.RotateHue(hue.Value)
		}
		screen.DrawImage(sprite.Image, op)
	})
}
