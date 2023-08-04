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

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func CreateApp() (*herd.App, error) {
	app := herd.NewApp()

	settings := component.Settings{
		Ticker:   time.NewTicker(500 * time.Millisecond),
		Gpu:      helper.GpuInfo(),
		Tps:      helper.NewPlot(20, 60),
		Fps:      helper.NewPlot(20, 60),
		Objects:  helper.NewPlot(20, 60000),
		Sprite:   assets.Bunny,
		Colorful: false,
		Amount:   1000,
	}

	velocity, err := system.NewVelocity(app)
	if err != nil {
		return nil, err
	}

	gravity, err := system.NewGravity(app)
	if err != nil {
		return nil, err
	}

	bounce, err := system.NewBounce(app)
	if err != nil {
		return nil, err
	}

	metrics, err := system.NewMetrics(app, &settings)
	if err != nil {
		return nil, err
	}

	spawn, err := system.NewSpawn(app, &settings)
	if err != nil {
		return nil, err
	}

	render, err := system.NewRender(app)
	if err != nil {
		return nil, err
	}

	if err := app.AddSystems(
		velocity.Update, gravity.Update, bounce.Update,
		bounce.Update, metrics.Update, spawn.Update,
	); err != nil {
		log.Fatal(err)
	}

	if err := app.AddRenderers(system.Background, render.Draw, metrics.Draw); err != nil {
		log.Fatal(err)
	}

	return app, nil
}

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

	app, err := CreateApp()
	if err != nil {
		log.Fatal(err)
	}

	if err := ebiten.RunGame(app); err != nil {
		log.Fatal(err)
	}
}
