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

type Game struct {
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

func NewGame() *Game {
	game := &Game{
		manager:   NewManager(),
		storage:   internal.NewStorage[EntityID](),
		resources: make(map[reflect.Type]reflect.Value),
		info:      &SystemInfo{},
	}

	game.resources[systemInfoType] = reflect.ValueOf(game.info)

	return game
}

func (g *Game) AddStartupSystems(systems ...any) error {
	for _, system := range systems {
		f, err := g.wrapSystem(system)
		if err != nil {
			return fmt.Errorf("unable to add startup system %v: %w", getFuncName(system), err)
		}

		g.startupSystems = append(g.startupSystems, f)
	}

	return nil
}

func (g *Game) AddSystems(systems ...any) error {
	for _, system := range systems {
		f, err := g.wrapSystem(system)
		if err != nil {
			return fmt.Errorf("unable to add system %v: %w", getFuncName(system), err)
		}

		g.systems = append(g.systems, f)
	}

	return nil
}

func (g *Game) AddRenderers(renderers ...any) error {
	for _, renderer := range renderers {
		f, err := g.wrapRenderer(renderer)
		if err != nil {
			return fmt.Errorf("unable to add renderer %v: %w", getFuncName(renderer), err)
		}

		g.renderers = append(g.renderers, f)
	}

	return nil
}

func getFuncName(fnc interface{}) string {
	strs := strings.Split((runtime.FuncForPC(reflect.ValueOf(fnc).Pointer()).Name()), "/")

	return strs[len(strs)-1]
}

func (g *Game) Update() error {
	if !g.alreadyUpdated {
		for _, system := range g.startupSystems {
			system()
		}
		g.alreadyUpdated = true
	}

	for _, system := range g.systems {
		system()
	}

	for _, spawn := range *g.manager.spawnQueue {
		id := g.newEntity()
		for _, comp := range spawn {
			g.storage.Add(id, comp)
		}
	}
	g.manager.clear()

	g.info.Entities = g.storage.Count()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, renderer := range g.renderers {
		renderer(screen)
	}
}

func (g *Game) Layout(w, h int) (int, int) {
	g.info.Bounds = image.Rect(0, 0, w, h)

	return w, h
}

func (g *Game) wrapSystem(system any) (func(), error) {
	val := reflect.ValueOf(system)
	typ := val.Type()

	if typ.Kind() != reflect.Func {
		return nil, fmt.Errorf("system should be a function, got %s", typ.Kind())
	}

	if typ.NumOut() != 0 {
		return nil, fmt.Errorf("system should not return any value, got %d out parameters", typ.NumOut())
	}

	pkg := reflect.TypeOf(Game{}).PkgPath()

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
			args[i] = g.createQuery(argType)
		case strings.HasPrefix(name, "Query2["), strings.HasPrefix(name, "Query3["):
			args[i] = g.createMultiQuery(argType)
		case strings.HasPrefix(name, "Res["):
			args[i] = g.createRes(argType)
		case name == "Manager":
			args[i] = reflect.ValueOf(g.manager)
		default:
			return nil, fmt.Errorf("system should takes one of Manager, Query[T] or Res[T] argument types, got %s", name)
		}
	}

	return func() {
		val.Call(args)
	}, nil

}

func (g *Game) wrapRenderer(renderer any) (func(*ebiten.Image), error) {
	val := reflect.ValueOf(renderer)
	typ := val.Type()

	if typ.Kind() != reflect.Func {
		return nil, fmt.Errorf("system should be a function, got %s", typ.Kind())
	}

	if typ.NumOut() != 0 {
		return nil, fmt.Errorf("system should not return any value, got %d out parameters", typ.NumOut())
	}

	pkg := reflect.TypeOf(Game{}).PkgPath()

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
			args[i] = g.createQuery(argType)
		case strings.HasPrefix(name, "Query2["), strings.HasPrefix(name, "Query3["):
			args[i] = g.createMultiQuery(argType)
		case strings.HasPrefix(name, "Res["):
			args[i] = g.createRes(argType)
		case name == "Manager":
			args[i] = reflect.ValueOf(g.manager)
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

func (g *Game) createRes(typ reflect.Type) reflect.Value {
	typeParam := typ.Elem()

	if _, ok := g.resources[typeParam]; !ok {
		g.resources[typeParam] = reflect.New(typeParam)
	}

	return g.resources[typeParam]
}

func (g *Game) createMultiQuery(typ reflect.Type) reflect.Value {
	val := reflect.New(typ).Elem()

	for i := 0; i < val.NumField(); i++ {
		queryField := val.Field(i)

		queryField.Set(g.createQuery(queryField.Type()))

	}

	return val
}

func (g *Game) createQuery(typ reflect.Type) reflect.Value {
	val := reflect.New(typ).Elem()
	typeParam := extractQueueTypeParameter(typ)

	val.Field(0).Set(reflect.ValueOf(g.storage.Indice(typeParam)))
	val.Field(1).Set(g.storage.Storage(typeParam))

	return val
}

func extractQueueTypeParameter(queryTyp reflect.Type) reflect.Type {
	return queryTyp.Field(1).Type.Elem().Elem()
}

func (g *Game) newEntity() EntityID {
	g.lastEntity++

	return g.lastEntity
}
