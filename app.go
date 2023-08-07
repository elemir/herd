package herd

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/elemir/herd/internal"
)

type EntityID int

// Startup is a initilize system that updates after App started. Startup system should return (true, nil) when initialisation finished successul
type Startup func() (bool, error)

// System is a basic system that updates every tick after App initialization is finished
type System func() error

// Renderer is a system that draws something every frame
type Renderer func(*ebiten.Image)

type startupInfo struct {
	system   Startup
	finished bool
}

// An App incapsulated all game logic and rendering. An App provides methods to adding systems and renderers and also implements ebiten.Game interface
type App struct {
	alreadyUpdated bool
	lastEntity     EntityID

	storage   internal.Storage[EntityID]
	systems   []System
	renderers []Renderer

	startups    []startupInfo
	initialized bool

	Manager    *Manager
	SystemInfo *SystemInfo
}

// NewApp returns a new App instance
func NewApp() *App {
	app := &App{
		storage:    internal.NewStorage[EntityID](),
		Manager:    newManager(),
		SystemInfo: &SystemInfo{},
	}

	return app
}

func (app *App) AddSystems(systems ...System) error {
	for _, system := range systems {
		app.systems = append(app.systems, system)
	}

	return nil
}

func (app *App) AddStartups(startups ...Startup) error {
	for _, startup := range startups {
		app.startups = append(app.startups, startupInfo{
			system: startup,
		})
	}

	return nil
}

func (app *App) AddRenderers(renderers ...Renderer) error {
	for _, renderer := range renderers {
		app.renderers = append(app.renderers, renderer)
	}

	return nil
}

func (app *App) Update() error {
	if !app.initialized {
		initialized := true

		for i, startup := range app.startups {
			if startup.finished {
				continue
			}

			finished, err := startup.system()
			if err != nil {
				return err
			}
			app.startups[i].finished = finished
			initialized = initialized && finished
		}

		app.initialized = initialized
		if !initialized {
			return nil
		}
	}

	for _, system := range app.systems {
		if err := system(); err != nil {
			return err
		}
	}

	for _, bundle := range app.Manager.spawnQueue {
		id := app.newEntity()
		if err := app.storage.Add(id, bundle); err != nil {
			return nil
		}
	}

	app.Manager.clear()
	app.SystemInfo.Entities = app.storage.Count()

	return nil
}

func (app *App) Draw(screen *ebiten.Image) {
	for _, renderer := range app.renderers {
		renderer(screen)
	}
}

func (app *App) Layout(w, h int) (int, int) {
	app.SystemInfo.Bounds = image.Rect(0, 0, w, h)

	return w, h
}

func (app *App) newEntity() EntityID {
	app.lastEntity++

	return app.lastEntity
}
