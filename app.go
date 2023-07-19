package herd

import (
	"fmt"
	"image"
	"reflect"
	"runtime"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/elemir/herd/internal"
)

type EntityID int

type App struct {
	alreadyUpdated bool
	lastEntity     EntityID

	storage        internal.Storage[EntityID]
	systems        []func()
	startupSystems []func()
	renderers      []func(*ebiten.Image)

	manager Manager

	info *SystemInfo

	resources map[reflect.Type]reflect.Value
}

func NewApp() *App {
	app := &App{
		manager:   NewManager(),
		storage:   internal.NewStorage[EntityID](),
		resources: make(map[reflect.Type]reflect.Value),
		info:      &SystemInfo{},
	}

	app.resources[systemInfoType] = reflect.ValueOf(app.info)

	return app
}

func (app *App) AddStartupSystems(systems ...any) error {
	for _, system := range systems {
		f, err := app.wrapSystem(system)
		if err != nil {
			return fmt.Errorf("unable to add startup system %v: %w", getFuncName(system), err)
		}

		app.startupSystems = append(app.startupSystems, f)
	}

	return nil
}

func (app *App) AddSystems(systems ...any) error {
	for _, system := range systems {
		f, err := app.wrapSystem(system)
		if err != nil {
			return fmt.Errorf("unable to add system %v: %w", getFuncName(system), err)
		}

		app.systems = append(app.systems, f)
	}

	return nil
}

func (app *App) AddRenderers(renderers ...any) error {
	for _, renderer := range renderers {
		f, err := app.wrapRenderer(renderer)
		if err != nil {
			return fmt.Errorf("unable to add renderer %v: %w", getFuncName(renderer), err)
		}

		app.renderers = append(app.renderers, f)
	}

	return nil
}

func getFuncName(fnc interface{}) string {
	strs := strings.Split((runtime.FuncForPC(reflect.ValueOf(fnc).Pointer()).Name()), "/")

	return strs[len(strs)-1]
}

func (app *App) Update() error {
	if !app.alreadyUpdated {
		for _, system := range app.startupSystems {
			system()
		}
		app.alreadyUpdated = true
	}

	for _, system := range app.systems {
		system()
	}

	for _, spawn := range *app.manager.spawnQueue {
		id := app.newEntity()
		for _, comp := range spawn {
			app.storage.Add(id, comp)
		}
	}
	app.manager.clear()

	app.info.Entities = app.storage.Count()

	return nil
}

func (app *App) Draw(screen *ebiten.Image) {
	for _, renderer := range app.renderers {
		renderer(screen)
	}
}

func (app *App) Layout(w, h int) (int, int) {
	app.info.Bounds = image.Rect(0, 0, w, h)

	return w, h
}

func (app *App) wrapSystem(system any) (func(), error) {
	val := reflect.ValueOf(system)
	typ := val.Type()

	if typ.Kind() != reflect.Func {
		return nil, fmt.Errorf("system should be a function, got %s", typ.Kind())
	}

	if typ.NumOut() != 0 {
		return nil, fmt.Errorf("system should not return any value, got %d out parameters", typ.NumOut())
	}

	pkg := reflect.TypeOf(App{}).PkgPath()

	args := make([]reflect.Value, typ.NumIn())

	// TODO(@elemir90): avoid pass one argument several types
	for i := 0; i < typ.NumIn(); i++ {
		argType := typ.In(i)
		name := argType.Name()

		if argType.PkgPath() != pkg {
			return nil, fmt.Errorf("system should takes arguments of Manager, Query[T] or Res[T] types, got %s", name)
		}

		switch {
		case strings.HasPrefix(name, "Query["):
			args[i] = app.createQuery(argType)
		case strings.HasPrefix(name, "Query2["), strings.HasPrefix(name, "Query3["):
			args[i] = app.createMultiQuery(argType)
		case strings.HasPrefix(name, "Res["):
			args[i] = app.createRes(argType)
		case name == "Manager":
			args[i] = reflect.ValueOf(app.manager)
		default:
			return nil, fmt.Errorf("system should takes one of Manager, Query[T] or Res[T] argument types, got %s", name)
		}
	}

	return func() {
		val.Call(args)
	}, nil

}

func (app *App) wrapRenderer(renderer any) (func(*ebiten.Image), error) {
	val := reflect.ValueOf(renderer)
	typ := val.Type()

	if typ.Kind() != reflect.Func {
		return nil, fmt.Errorf("system should be a function, got %s", typ.Kind())
	}

	if typ.NumOut() != 0 {
		return nil, fmt.Errorf("system should not return any value, got %d out parameters", typ.NumOut())
	}

	pkg := reflect.TypeOf(App{}).PkgPath()

	args := make([]reflect.Value, typ.NumIn())

	screenIdx := -1
	// TODO(@elemir90): avoid pass one argument several types
	for i := 0; i < typ.NumIn(); i++ {
		argType := typ.In(i)
		name := argType.Name()

		if argType.String() == "*ebiten.Image" {
			screenIdx = i
			args[i] = reflect.New(argType).Elem()

			continue
		}

		if argType.PkgPath() != pkg {
			return nil, fmt.Errorf("system should takes arguments of Manager, Query[T] or Res[T] types, got %s", name)
		}

		switch {
		case strings.HasPrefix(name, "Query["):
			args[i] = app.createQuery(argType)
		case strings.HasPrefix(name, "Query2["), strings.HasPrefix(name, "Query3["):
			args[i] = app.createMultiQuery(argType)
		case strings.HasPrefix(name, "Res["):
			args[i] = app.createRes(argType)
		case name == "Manager":
			args[i] = reflect.ValueOf(app.manager)
		default:
			return nil, fmt.Errorf("system should take one of Manager, Query[T] or Res[T] argument types, got %s", name)
		}
	}

	if screenIdx == -1 {
		return nil, fmt.Errorf("renderer should take exactly one *screen.Image argument")
	}

	return func(screen *ebiten.Image) {
		args[screenIdx].Set(reflect.ValueOf(screen))

		val.Call(args)
	}, nil

}

func (app *App) createRes(typ reflect.Type) reflect.Value {
	typeParam := typ.Elem()

	if _, ok := app.resources[typeParam]; !ok {
		app.resources[typeParam] = reflect.New(typeParam)
	}

	return app.resources[typeParam]
}

func (app *App) createMultiQuery(typ reflect.Type) reflect.Value {
	val := reflect.New(typ).Elem()

	for i := 0; i < val.NumField(); i++ {
		queryField := val.Field(i)

		queryField.Set(app.createQuery(queryField.Type()))

	}

	return val
}

func (app *App) createQuery(typ reflect.Type) reflect.Value {
	val := reflect.New(typ).Elem()
	typeParam := extractQueueTypeParameter(typ)

	val.Field(0).Set(reflect.ValueOf(app.storage.Indice(typeParam)))
	val.Field(1).Set(app.storage.Storage(typeParam))

	return val
}

func extractQueueTypeParameter(queryTyp reflect.Type) reflect.Type {
	return queryTyp.Field(1).Type.Elem().Elem()
}

func (app *App) newEntity() EntityID {
	app.lastEntity++

	return app.lastEntity
}
