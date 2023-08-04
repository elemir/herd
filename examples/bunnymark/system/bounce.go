package system

import (
	"github.com/elemir/herd"

	"github.com/elemir/herd/examples/bunnymark/component"
	"github.com/elemir/herd/examples/bunnymark/helper"
)

type Texture struct {
	Pos    component.Position
	Vel    component.Velocity
	Sprite component.Sprite
}

type Bounce struct {
	System *herd.SystemInfo
	Query  herd.Query[Texture]
}

func NewBounce(app *herd.App) (Bounce, error) {
	query, err := herd.NewQuery[Texture](app)
	if err != nil {
		return Bounce{}, err
	}

	return Bounce{
		System: app.SystemInfo,
		Query:  query,
	}, nil
}

func (b Bounce) Update() error {
	sw, sh := float64(b.System.Bounds.Dx()), float64(b.System.Bounds.Dy())

	b.Query.ForEach(func(texture *Texture) {
		iw, ih := float64(texture.Sprite.Image.Bounds().Dx()), float64(texture.Sprite.Image.Bounds().Dy())
		relW, relH := iw/sw, ih/sh
		if texture.Pos.X+relW > 1 {
			texture.Vel.X *= -1
			texture.Pos.X = 1 - relW
		}
		if texture.Pos.X < 0 {
			texture.Vel.X *= -1
			texture.Pos.X = 0
		}
		if texture.Pos.Y+relH > 1 {
			texture.Vel.Y *= -0.85
			texture.Pos.Y = 1 - relH
			if helper.Chance(0.5) {
				texture.Vel.Y -= helper.RangeFloat(0, 0.009)
			}
		}
		if texture.Pos.Y < 0 {
			texture.Vel.Y = 0
			texture.Pos.Y = 0
		}
	})

	return nil
}
