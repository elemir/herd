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

func CalculateMetrics(settings herd.Res[component.Settings], system herd.Res[herd.SystemInfo]) {
	select {
	case <-settings.Ticker.C:
		settings.Objects.Update(float64(system.Entities))
		settings.Tps.Update(ebiten.CurrentTPS())
		settings.Fps.Update(ebiten.CurrentFPS())
	default:
	}
}

func DrawMetrics(screen *ebiten.Image, settings herd.Res[component.Settings], system herd.Res[herd.SystemInfo]) {
	str := fmt.Sprintf(
		"GPU: %s\nTPS: %.2f, FPS: %.2f, Objects: %.f\nBatching: %t, Amount: %d\nResolution: %dx%d",
		settings.Gpu, settings.Tps.Last(), settings.Fps.Last(), settings.Objects.Last(),
		!settings.Colorful, settings.Amount,
		system.Bounds.Dx(), system.Bounds.Dy(),
	)

	rect := text.BoundString(basicfont.Face7x13, str)
	width, height := float64(rect.Dx()), float64(rect.Dy())

	padding := 20.0
	rectW, rectH := width+padding, height+padding
	plotW, plotH := 100.0, 40.0

	ebitenutil.DrawRect(screen, 0, 0, rectW, rectH, color.RGBA{A: 128})
	text.Draw(screen, str, basicfont.Face7x13, int(padding)/2, 10+int(padding)/2, colornames.White)

	settings.Tps.Draw(screen, 0, padding+rectH, plotW, plotH)
	settings.Fps.Draw(screen, 0, padding+rectH*2, plotW, plotH)
	settings.Objects.Draw(screen, 0, padding+rectH*3, plotW, plotH)
}
