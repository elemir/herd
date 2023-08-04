package herd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type Bundle struct {
	X int
	Y string
}

type SimpleX struct {
	X int
}

func TestQuery(t *testing.T) {
	app := NewApp()

	expected := []Bundle{{10, "10"}, {23, "23"}, {45, "45"}}
	for _, in := range expected {
		app.Manager.Spawn(in)
	}
	app.Manager.Spawn(SimpleX{42})

	query, err := NewQuery[Bundle](app)
	require.NoError(t, err)

	var output []Bundle
	err = app.AddSystems(func() error {
		query.ForEach(func(b *Bundle) {
			output = append(output, *b)
		})

		return nil
	})
	require.NoError(t, err)

	err = app.Update()
	require.NoError(t, err)
	require.Nil(t, output)

	err = app.Update()
	require.NoError(t, err)
	require.ElementsMatch(t, expected, output)
}

func TestOneComponent(t *testing.T) {
	app := NewApp()

	expected := []SimpleX{{10}, {23}, {45}, {42}}
	for _, in := range []Bundle{{10, "10"}, {23, "23"}, {45, "45"}} {
		app.Manager.Spawn(in)
	}
	app.Manager.Spawn(SimpleX{42})

	query, err := NewQuery[SimpleX](app)
	require.NoError(t, err)

	var output []SimpleX

	err = app.AddSystems(func() error {
		query.ForEach(func(x *SimpleX) {
			output = append(output, *x)
		})

		return nil
	})
	require.NoError(t, err)

	err = app.Update()
	require.NoError(t, err)
	require.Nil(t, output)

	err = app.Update()
	require.NoError(t, err)
	require.ElementsMatch(t, expected, output)
}

func TestAnonymouseInline(t *testing.T) {
	app := NewApp()

	input := []struct {
		int
		string
	}{{10, "22"}}

	for _, in := range input {
		app.Manager.Spawn(in)
	}
	app.Manager.Spawn(struct{ int }{42})

	query, err := NewQuery[struct{ int }](app)
	require.NoError(t, err)

	var output []struct{ int }
	err = app.AddSystems(func() error {
		query.ForEach(func(b *struct{ int }) {
			output = append(output, *b)
		})

		return nil
	})
	require.NoError(t, err)

	err = app.Update()
	require.NoError(t, err)
	require.Nil(t, output)

	err = app.Update()
	require.NoError(t, err)
	require.ElementsMatch(t, []struct{ int }{{42}, {10}}, output)
}
