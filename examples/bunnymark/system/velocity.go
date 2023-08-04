package system

import (
	"github.com/elemir/herd"
	"github.com/elemir/herd/examples/bunnymark/component"
)

type MoveBundle struct {
	Pos component.Position
	Vel component.Velocity
}

type Velocity struct {
	Query herd.Query[MoveBundle]
}

func NewVelocity(app *herd.App) (Velocity, error) {
	query, err := herd.NewQuery[MoveBundle](app)
	if err != nil {
		return Velocity{}, err
	}

	return Velocity{
		Query: query,
	}, nil
}

func (v Velocity) Update() error {
	v.Query.ForEach(func(mv *MoveBundle) {
		mv.Pos.X += mv.Vel.X
		mv.Pos.Y += mv.Vel.Y
	})

	return nil
}
