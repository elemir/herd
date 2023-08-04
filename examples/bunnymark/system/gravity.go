package system

import (
	"github.com/elemir/herd"

	"github.com/elemir/herd/examples/bunnymark/component"
)

type GravityBundle struct {
	Vel     component.Velocity
	Gravity component.Gravity
}

type Gravity struct {
	Query herd.Query[GravityBundle]
}

func NewGravity(app *herd.App) (Gravity, error) {
	query, err := herd.NewQuery[GravityBundle](app)
	if err != nil {
		return Gravity{}, err
	}

	return Gravity{
		Query: query,
	}, nil
}

func (g Gravity) Update() error {
	g.Query.ForEach(func(gb *GravityBundle) {
		gb.Vel.Y += gb.Gravity.Value
	})

	return nil
}
