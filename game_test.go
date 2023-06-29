package herd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQuery(t *testing.T) {
	game := NewGame()

	var input, output []int

	input = []int{10, 23, 45, 12}

	err := game.AddStartupSystems(func(manager Manager) {
		for _, in := range input {
			manager.Spawn(in)
		}
	})
	require.NoError(t, err)

	err = game.AddSystems(func(query Query[int]) {
		query.ForEach(func(_ EntityID, val *int) {
			output = append(output, *val)
		})
	})
	require.NoError(t, err)

	game.Update()
	require.Nil(t, output)

	game.Update()
	require.ElementsMatch(t, input, output)
}

type testPair struct {
	num int
	str string
}

func TestQuery2(t *testing.T) {
	game := NewGame()

	var input, output []testPair

	input = []testPair{
		{10, "one"},
		{1, "two"},
		{30, "three"},
	}

	err := game.AddStartupSystems(func(manager Manager) {
		for _, in := range input {
			manager.Spawn(in.num, in.str)
		}
	})
	require.NoError(t, err)

	err = game.AddSystems(func(query Query2[int, string]) {
		query.ForEach(func(_ EntityID, num *int, str *string) {
			output = append(output, testPair{
				num: *num, str: *str,
			})
		})
	})
	require.NoError(t, err)

	game.Update()
	require.Nil(t, output)

	game.Update()
	require.ElementsMatch(t, input, output)
}

func TestRes(t *testing.T) {
	game := NewGame()

	var input, output int

	input = 10

	err := game.AddStartupSystems(func(res Res[int]) {
		*res = input
	})
	require.NoError(t, err)

	err = game.AddSystems(func(res Res[int]) {
		output = *res
	})
	require.NoError(t, err)

	game.Update()
	require.Equal(t, input, output)
}
