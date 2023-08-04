package system

import (
	"github.com/elemir/herd"
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/elemir/herd/examples/bunnymark/component"
)

type Tile struct {
	Pos    component.Position
	Hue    component.Hue
	Sprite component.Sprite
}

type Render struct {
	Query herd.Query[Tile]
}

func NewRender(app *herd.App) (Render, error) {
	query, err := herd.NewQuery[Tile](app)
	if err != nil {
		return Render{}, err
	}

	return Render{
		Query: query,
	}, nil
}

func (r Render) Draw(screen *ebiten.Image) {
	r.Query.ForEach(func(tile *Tile) {
		sw, sh := float64(screen.Bounds().Dx()), float64(screen.Bounds().Dy())
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(tile.Pos.X*sw, tile.Pos.Y*sh)
		if *tile.Hue.Colorful {
			op.ColorM.RotateHue(tile.Hue.Value)
		}
		screen.DrawImage(tile.Sprite.Image, op)
	})
}
