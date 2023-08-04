package system

import (
	"math"

	"github.com/elemir/herd"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/elemir/herd/examples/bunnymark/component"
	"github.com/elemir/herd/examples/bunnymark/helper"
)

type Spawn struct {
	Manager  *herd.Manager
	Settings *component.Settings
	System   *herd.SystemInfo
}

func NewSpawn(app *herd.App, settings *component.Settings) (Spawn, error) {
	return Spawn{
		Manager:  app.Manager,
		Settings: settings,
		System:   app.SystemInfo,
	}, nil
}

func (s Spawn) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		s.addBunnies()
	}

	if ids := ebiten.AppendTouchIDs(nil); len(ids) > 0 {
		s.addBunnies()
	}

	if _, offset := ebiten.Wheel(); offset != 0 {
		s.Settings.Amount += int(offset * 10)
		if s.Settings.Amount < 0 {
			s.Settings.Amount = 0
		}
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		s.Settings.Colorful = !s.Settings.Colorful
	}

	return nil
}

type Bunnie struct {
	Pos     component.Position
	Vel     component.Velocity
	Hue     component.Hue
	Gravity component.Gravity
	Sprite  component.Sprite
}

func (s Spawn) addBunnies() {
	// Spawns specific amount of bunnies at the edges of the screen
	// It will alternately add bunnies to the left and right corners of the screen
	for i := 0; i < s.Settings.Amount; i++ {
		s.Manager.Spawn(Bunnie{component.Position{
			X: float64(s.System.Entities % 2), // Alternate screen edges
		}, component.Velocity{
			X: helper.RangeFloat(0, 0.005),
			Y: helper.RangeFloat(0.0025, 0.005),
		}, component.Hue{
			Colorful: &s.Settings.Colorful,
			Value:    helper.RangeFloat(0, 2*math.Pi),
		}, component.Gravity{
			Value: 0.00095,
		}, component.Sprite{
			Image: s.Settings.Sprite,
		}})
	}
}
