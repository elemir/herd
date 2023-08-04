package herd

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/elemir/herd/internal"
)

type EntityID int

type System func() error

type Renderer func(*ebiten.Image)

type App struct {
	alreadyUpdated bool
	lastEntity     EntityID

	storage   internal.Storage[EntityID]
	systems   []System
	renderers []Renderer

	Manager    *Manager
	SystemInfo *SystemInfo
}

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

func (app *App) AddRenderers(renderers ...Renderer) error {
	for _, renderer := range renderers {
		app.renderers = append(app.renderers, renderer)
	}

	return nil
}

func (app *App) Update() error {
	for _, system := range app.systems {
		if err := system(); err != nil {
			return nil
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
