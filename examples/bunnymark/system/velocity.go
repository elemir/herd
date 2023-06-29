package system

import (
	"github.com/elemir/herd"
	"github.com/elemir/herd/examples/bunnymark/component"
)

func Velocity(query herd.Query2[component.Position, component.Velocity]) {
	query.ForEach(func(_ herd.EntityID, pos *component.Position, vel *component.Velocity) {
		pos.X += vel.X
		pos.Y += vel.Y
	})
}
