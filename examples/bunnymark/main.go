package main

import (
	"flag"
	"log"
	"math/rand"
	"os"
	"runtime/pprof"
	"time"

	"github.com/elemir/herd"
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/elemir/herd/examples/bunnymark/assets"
	"github.com/elemir/herd/examples/bunnymark/component"
	"github.com/elemir/herd/examples/bunnymark/helper"
	"github.com/elemir/herd/examples/bunnymark/system"
)

func Settings(settings herd.Res[component.Settings]) {
	*settings = component.Settings{
		Ticker:   time.NewTicker(500 * time.Millisecond),
		Gpu:      helper.GpuInfo(),
		Tps:      helper.NewPlot(20, 60),
		Fps:      helper.NewPlot(20, 60),
		Objects:  helper.NewPlot(20, 60000),
		Sprite:   assets.Bunny,
		Colorful: false,
		Amount:   100,
	}
}

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowSizeLimits(300, 200, -1, -1)
	ebiten.SetFPSMode(ebiten.FPSModeVsyncOffMaximum)
	ebiten.SetWindowResizable(true)
	rand.Seed(time.Now().UTC().UnixNano())

	game := herd.NewGame()

	if err := game.AddStartupSystems(Settings); err != nil {
		log.Fatal(err)
	}

	if err := game.AddSystems(
		system.Velocity, system.Gravity, system.Bounce,
		system.CalculateMetrics, system.Spawn,
	); err != nil {
		log.Fatal(err)
	}

	if err := game.AddRenderers(system.Background, system.Render, system.DrawMetrics); err != nil {
		log.Fatal(err)
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
