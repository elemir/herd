package internal

import (
	"fmt"
	"reflect"
	"unsafe"
)

type FieldType struct {
	Name string
	Type reflect.Type
}

type Storage[ID comparable] struct {
	arrays map[FieldType]*SparseArray[ID]

	fullIndex map[ID]struct{}
}

func NewStorage[ID comparable]() Storage[ID] {
	return Storage[ID]{
		arrays:    make(map[FieldType]*SparseArray[ID]),
		fullIndex: make(map[ID]struct{}),
	}
}

func (s *Storage[ID]) Add(id ID, bundle any) error {
	typ := reflect.TypeOf(bundle)

	if typ.Kind() != reflect.Struct {
		return fmt.Errorf("bundle type should be a struct, got %s", typ.Kind())
	}

	for i := 0; i < typ.NumField(); i++ {
		fieldType := typ.Field(i)

		array := s.sparseArray(fieldType.Name, fieldType.Type)
		fieldPtr := unsafe.Add(toPointer(bundle), fieldType.Offset)
		array.Add(id, fieldPtr)
	}

	s.fullIndex[id] = struct{}{}

	return nil
}

type iface struct {
	Type, Data unsafe.Pointer
}

func toPointer(val any) unsafe.Pointer {
	return (*iface)(unsafe.Pointer(&val)).Data
}

func (s *Storage[ID]) sparseArray(name string, typ reflect.Type) *SparseArray[ID] {
	field := FieldType{
		Name: name,
		Type: typ,
	}

	if s.arrays[field] == nil {
		s.arrays[field] = NewSparseArray[ID]()
	}

	return s.arrays[field]
}

func (s *Storage[ID]) Count() int {
	return len(s.fullIndex)
}

func (s *Storage[ID]) Iterator(typ reflect.Type) (Iterator[ID], error) {
	if typ.Kind() != reflect.Struct {
		return Iterator[ID]{}, fmt.Errorf("bundle type should be a struct, got %s", typ.Kind())
	}

	val := reflect.New(typ).Elem()

	arrays := make([]*SparseArray[ID], val.NumField())
	fields := make([]FieldValue, val.NumField())

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		fields[i] = FieldValue{
			pointer: unsafe.Add(val.Addr().UnsafePointer(), fieldType.Offset),
			size:    fieldType.Type.Size(),
		}

		arrays[i] = s.sparseArray(typ.Field(i).Name, field.Type())
	}

	return Iterator[ID]{
		elem:      val.Addr().UnsafePointer(),
		fields:    fields,
		arrays:    arrays,
		fullIndex: s.fullIndex,
	}, nil
}

type FieldValue struct {
	pointer unsafe.Pointer
	size    uintptr
}

type Iterator[ID comparable] struct {
	elem unsafe.Pointer

	arrays []*SparseArray[ID]
	fields []FieldValue

	fullIndex map[ID]struct{}
}

func (iter Iterator[ID]) ForEach(f func(ID, unsafe.Pointer) bool) {
	if len(iter.arrays) == 0 {
		for id := range iter.fullIndex {
			if cont := f(id, iter.elem); !cont {
				break
			}
		}
		return
	}

	positions := make([]int, len(iter.arrays))

EntityLoop:
	for id, pos := range iter.arrays[0].index {
		positions[0] = pos
		for i, array := range iter.arrays {
			pos, ok := array.index[id]
			if !ok {
				continue EntityLoop
			}
			positions[i] = pos
		}

		for i, pos := range positions {
			copyPointer(iter.fields[i].pointer, iter.arrays[i].slice[pos], iter.fields[i].size)
		}

		if cont := f(id, iter.elem); !cont {
			break
		}

		for i, pos := range positions {
			copyPointer(iter.arrays[i].slice[pos], iter.fields[i].pointer, iter.fields[i].size)
		}
	}
}

func copyPointer(dst unsafe.Pointer, src unsafe.Pointer, size uintptr) {
	srcSlice := unsafe.Slice((*byte)(src), size)
	dstSlice := unsafe.Slice((*byte)(dst), size)

	copy(dstSlice, srcSlice)
}
