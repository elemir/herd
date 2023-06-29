package herd

import (
	"image"
	"reflect"
)

type Res[T any] *T

var (
	systemInfoType = reflect.TypeOf(SystemInfo{})
)

type SystemInfo struct {
	Entities int
	Bounds   image.Rectangle
}
