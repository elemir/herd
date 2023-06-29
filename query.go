package herd

// TODO(@elemir90): understand how use private fields here
type Query[T any] struct {
	Index   map[EntityID]int
	Storage *[]T
}

func (q Query[T]) ForEach(f func(id EntityID, t *T)) {
	for id, i := range q.Index {
		f(id, &(*q.Storage)[i])
	}
}

type Query2[T, U any] struct {
	QueryT Query[T]
	QueryU Query[U]
}

func (q Query2[T, U]) ForEach(f func(id EntityID, t *T, u *U)) {
	for id, tidx := range q.QueryT.Index {
		uidx, ok := q.QueryU.Index[id]
		if !ok {
			continue
		}
		f(id, &(*q.QueryT.Storage)[tidx], &(*q.QueryU.Storage)[uidx])
	}
}

type Query3[T, U, V any] struct {
	QueryT Query[T]
	QueryU Query[U]
	QueryV Query[V]
}

func (q Query3[T, U, V]) ForEach(f func(id EntityID, t *T, u *U, v *V)) {
	for id, tidx := range q.QueryT.Index {
		uidx, ok := q.QueryU.Index[id]
		if !ok {
			continue
		}
		vidx, ok := q.QueryV.Index[id]
		if !ok {
			continue
		}
		f(id, &(*q.QueryT.Storage)[tidx], &(*q.QueryU.Storage)[uidx], &(*q.QueryV.Storage)[vidx])
	}
}
