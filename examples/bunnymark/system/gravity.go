package system

import (
	"github.com/elemir/herd"

	"github.com/elemir/herd/examples/bunnymark/component"
)

func Gravity(query herd.Query2[component.Velocity, component.Gravity]) {
	query.ForEach(func(_ herd.EntityID, vel *component.Velocity, grav *component.Gravity) {
		vel.Y += grav.Value
	})
}
