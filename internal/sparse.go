package internal

import (
	"unsafe"
)

type SparseArray[ID comparable] struct {
	index map[ID]int
	slice []unsafe.Pointer
}

func NewSparseArray[ID comparable]() *SparseArray[ID] {
	return &SparseArray[ID]{
		index: make(map[ID]int),
	}
}

func (arr *SparseArray[ID]) Add(id ID, ptr unsafe.Pointer) {
	pos := len(arr.slice)
	arr.slice = append(arr.slice, ptr)
	arr.index[id] = pos
}
