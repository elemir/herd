package system

import (
	"math"

	"github.com/elemir/herd"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/elemir/herd/examples/bunnymark/component"
	"github.com/elemir/herd/examples/bunnymark/helper"
)

func Spawn(mngr herd.Manager, settings herd.Res[component.Settings], system herd.Res[herd.SystemInfo]) {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		addBunnies(mngr, settings, system)
	}

	if ids := ebiten.AppendTouchIDs(nil); len(ids) > 0 {
		addBunnies(mngr, settings, system) // not accurate, cause no input manager for this
	}

	if _, offset := ebiten.Wheel(); offset != 0 {
		settings.Amount += int(offset * 10)
		if settings.Amount < 0 {
			settings.Amount = 0
		}
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		settings.Colorful = !settings.Colorful
	}
}

func addBunnies(mngr herd.Manager, settings herd.Res[component.Settings], system herd.Res[herd.SystemInfo]) {
	// Spawns specific amount of bunnies at the edges of the screen
	// It will alternately add bunnies to the left and right corners of the screen
	for i := 0; i < settings.Amount; i++ {
		mngr.Spawn(component.Position{
			X: float64(system.Entities % 2), // Alternate screen edges
		}, component.Velocity{
			X: helper.RangeFloat(0, 0.005),
			Y: helper.RangeFloat(0.0025, 0.005),
		}, component.Hue{
			Colorful: &settings.Colorful,
			Value:    helper.RangeFloat(0, 2*math.Pi),
		}, component.Gravity{
			Value: 0.00095,
		}, component.Sprite{
			Image: settings.Sprite,
		})
	}
}
