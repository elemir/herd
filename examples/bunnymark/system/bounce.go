package system

import (
	"github.com/elemir/herd"

	"github.com/elemir/herd/examples/bunnymark/component"
	"github.com/elemir/herd/examples/bunnymark/helper"
)

func Bounce(system herd.Res[herd.SystemInfo], query herd.Query3[component.Position, component.Velocity, component.Sprite]) {
	sw, sh := float64(system.Bounds.Dx()), float64(system.Bounds.Dy())

	query.ForEach(func(_ herd.EntityID, pos *component.Position, vel *component.Velocity, sprite *component.Sprite) {
		iw, ih := float64(sprite.Image.Bounds().Dx()), float64(sprite.Image.Bounds().Dy())
		relW, relH := iw/sw, ih/sh
		if pos.X+relW > 1 {
			vel.X *= -1
			pos.X = 1 - relW
		}
		if pos.X < 0 {
			vel.X *= -1
			pos.X = 0
		}
		if pos.Y+relH > 1 {
			vel.Y *= -0.85
			pos.Y = 1 - relH
			if helper.Chance(0.5) {
				vel.Y -= helper.RangeFloat(0, 0.009)
			}
		}
		if pos.Y < 0 {
			vel.Y = 0
			pos.Y = 0
		}
	})
}
