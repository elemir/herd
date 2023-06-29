package internal

import "reflect"

type Storage[ID comparable] struct {
	indices map[reflect.Type](map[ID]int)
	// TODO(@elemir90): make add component faster via self-written unsafe slice
	storages map[reflect.Type]reflect.Value

	fullIndex map[ID]struct{}
}

func NewStorage[ID comparable]() Storage[ID] {
	return Storage[ID]{
		indices:   make(map[reflect.Type](map[ID]int)),
		storages:  make(map[reflect.Type]reflect.Value),
		fullIndex: make(map[ID]struct{}),
	}
}

func (g *Storage[ID]) Indice(typ reflect.Type) map[ID]int {
	if _, ok := g.indices[typ]; !ok {
		g.indices[typ] = make(map[ID]int)
	}

	return g.indices[typ]
}

func (g *Storage[ID]) Storage(typ reflect.Type) reflect.Value {
	if _, ok := g.storages[typ]; !ok {
		sliceType := reflect.SliceOf(typ)
		g.storages[typ] = reflect.New(sliceType)
		g.storages[typ].Elem().Set(reflect.MakeSlice(sliceType, 0, 32))
	}

	return g.storages[typ]
}

func (g *Storage[ID]) Add(id ID, content any) {
	val := reflect.ValueOf(content)
	typ := val.Type()

	storage := g.Storage(typ)
	pos := storage.Elem().Len()
	storage.Elem().Set(reflect.Append(storage.Elem(), val))

	indice := g.Indice(typ)
	indice[id] = pos

	g.fullIndex[id] = struct{}{}
}

func (g *Storage[ID]) Count() int {
	return len(g.fullIndex)
}
