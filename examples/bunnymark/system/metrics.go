package system

import (
	"fmt"
	"image/color"

	"github.com/elemir/herd"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"

	"github.com/elemir/herd/examples/bunnymark/component"
)

type Metrics struct {
	Settings *component.Settings
	System   *herd.SystemInfo
}

func NewMetrics(app *herd.App, settings *component.Settings) (Metrics, error) {
	return Metrics{
		Settings: settings,
		System:   app.SystemInfo,
	}, nil
}

func (m Metrics) Update() error {
	select {
	case <-m.Settings.Ticker.C:
		m.Settings.Objects.Update(float64(m.System.Entities))
		m.Settings.Tps.Update(ebiten.CurrentTPS())
		m.Settings.Fps.Update(ebiten.CurrentFPS())
	default:
	}

	return nil
}

func (m Metrics) Draw(screen *ebiten.Image) {
	str := fmt.Sprintf(
		"GPU: %s\nTPS: %.2f, FPS: %.2f, Objects: %.f\nBatching: %t, Amount: %d\nResolution: %dx%d",
		m.Settings.Gpu, m.Settings.Tps.Last(), m.Settings.Fps.Last(), m.Settings.Objects.Last(),
		!m.Settings.Colorful, m.Settings.Amount,
		m.System.Bounds.Dx(), m.System.Bounds.Dy(),
	)

	rect := text.BoundString(basicfont.Face7x13, str)
	width, height := float64(rect.Dx()), float64(rect.Dy())

	padding := 20.0
	rectW, rectH := width+padding, height+padding
	plotW, plotH := 100.0, 40.0

	ebitenutil.DrawRect(screen, 0, 0, rectW, rectH, color.RGBA{A: 128})
	text.Draw(screen, str, basicfont.Face7x13, int(padding)/2, 10+int(padding)/2, colornames.White)

	m.Settings.Tps.Draw(screen, 0, padding+rectH, plotW, plotH)
	m.Settings.Fps.Draw(screen, 0, padding+rectH*2, plotW, plotH)
	m.Settings.Objects.Draw(screen, 0, padding+rectH*3, plotW, plotH)
}
