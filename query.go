package herd

import (
	"fmt"
	"unsafe"

	"github.com/elemir/herd/internal"
)

type Query[T any] struct {
	iterator internal.Iterator[EntityID]
}

func NewQuery[T any](app *App) (Query[T], error) {
	iterator, err := app.storage.Iterator(internal.TypeOf[T]())
	if err != nil {
		return Query[T]{}, fmt.Errorf("create iterator: %w", err)
	}

	return Query[T]{
		iterator: iterator,
	}, nil
}

func (q Query[T]) ForEach(f func(t *T)) {
	q.iterator.ForEach(func(_ EntityID, ptr unsafe.Pointer) bool {
		f((*T)(ptr))

		return true
	})
}

func (q Query[T]) Iterate(f func(id EntityID, t *T) bool) {
	q.iterator.ForEach(func(id EntityID, ptr unsafe.Pointer) bool {
		return f(id, (*T)(ptr))
	})
}
