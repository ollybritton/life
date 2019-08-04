package life_test

import (
	"testing"

	"gitlab.com/ollybritton/life"
)

func TestGridDimensions(t *testing.T) {
	var grid = life.NewGrid(10, 10)

	if grid.Width != 10 || grid.Height != 10 {
		t.Errorf(
			"grid height and width should be 10, 10, got %d and %d", grid.Width, grid.Height,
		)
	}
}

func TestGridString(t *testing.T) {
	var str string
	str += "..O\n"
	str += ".O.\n"
	str += "O.."

	var grid = life.NewGridFromString(str, "O", ".")

	if grid.Width != 3 || grid.Height != 3 {
		t.Errorf(
			"grid height and width should be 3, 3, got %d and %d", grid.Width, grid.Height,
		)
	}

	if grid.Get(0, 0) != 1 {
		t.Errorf(
			"(0,0) should be active",
		)
	}
}

func TestGridCoord(t *testing.T) {
	var str string
	str += ".O..O\n"
	str += ".O...\n"
	str += "O..O.\n"
	str += "O....\n"
	str += "O.O.."

	var grid = life.NewGridFromString(str, "O", ".")

	var coordtests = []struct {
		x, y           int
		xIndex, yIndex int
	}{
		{0, 0, 0, 4},
		{0, 4, 0, 0},
		{2, 1, 2, 3},
		{4, 4, 4, 0},
		{2, 2, 2, 2},
	}

	for _, test := range coordtests {
		xIndex, yIndex := grid.CoordToIndex(test.x, test.y)
		if xIndex != test.xIndex || yIndex != test.yIndex {
			t.Errorf("test %+v failed", test)
		}
	}
}

func TestGridIndex(t *testing.T) {
	var str string
	str += ".O..O\n"
	str += ".O...\n"
	str += "O..O.\n"
	str += "O....\n"
	str += "O.O.."

	var grid = life.NewGridFromString(str, "O", ".")

	var indextests = []struct {
		xIndex, yIndex int
		x, y           int
	}{
		{0, 4, 0, 0},
		{0, 0, 0, 4},
		{2, 3, 2, 1},
		{4, 0, 4, 4},
		{2, 2, 2, 2},
	}

	for _, test := range indextests {
		x, y := grid.CoordToIndex(test.xIndex, test.yIndex)
		if x != test.x || y != test.y {
			t.Errorf("test %+v failed", test)
		}
	}
}

func TestGridGet(t *testing.T) {
	var str string
	str += ".O..O\n"
	str += ".O...\n"
	str += "O..O.\n"
	str += "O....\n"
	str += "O.O.."

	var grid = life.NewGridFromString(str, "O", ".")

	var gettests = []struct {
		x, y int
		out  int
	}{
		{0, 0, 1},
		{4, 4, 1},
		{1, 1, 0},
		{3, 3, 0},
		{2, 2, 0},
		{1, 4, 1},
	}

	for _, test := range gettests {
		val := grid.Get(test.x, test.y)
		if val != test.out {
			t.Errorf("test %+v failed, expected %d but got %d instead", test, test.out, val)
		}
	}
}

func TestGridGetDefault(t *testing.T) {
	var str string
	str += ".O..O\n"
	str += ".O...\n"
	str += "O..O.\n"
	str += "O....\n"
	str += "O.O.."

	var grid = life.NewGridFromString(str, "O", ".")

	var defaulttests = []struct {
		x, y int
		out  int
	}{
		{0, 0, 1},
		{4, 4, 1},
		{1, 1, 0},
		{3, 3, 0},
		{2, 2, 0},
		{1, 4, 1},
		{5, 5, 0},
		{10, 10, 0},
		{-5, 20, 0},
		{0, -1, 0},
	}

	for _, test := range defaulttests {
		val := grid.GetDefault(test.x, test.y)
		if val != test.out {
			t.Errorf("test %+v failed, expected %d but got %d instead", test, test.out, val)
		}
	}
}

func TestGridGetModulo(t *testing.T) {
	var str string
	str += ".O..O\n"
	str += ".O...\n"
	str += "O..O.\n"
	str += "O....\n"
	str += "O.O.."

	var grid = life.NewGridFromString(str, "O", ".")

	var modulotest = []struct {
		x, y int
		out  int
	}{
		{0, 0, 1},
		{4, 4, 1},
		{1, 1, 0},
		{3, 3, 0},
		{2, 2, 0},
		{1, 4, 1},
		{5, 5, 0},
		{10, 10, 0},
		{-5, 20, 0},
		{0, -1, 0},
		{0, 1, 1},
	}

	for _, test := range modulotest {
		val := grid.GetDefault(test.x, test.y)
		if val != test.out {
			t.Errorf("test %+v failed, expected %d but got %d instead", test, test.out, val)
		}
	}
}

func TestGridSet(t *testing.T) {
	var str string
	str += ".O..O\n"
	str += ".O...\n"
	str += "O..O.\n"
	str += "O....\n"
	str += "O.O.."

	var grid = life.NewGridFromString(str, "O", ".")

	if grid.Set(0, 4, 1); grid.Get(0, 4) != 1 {
		t.Errorf("could not successfully set (0,4) to 1")
	}
}

func TestGridNeighbours(t *testing.T) {
	var str string
	str += ".O..O\n"
	str += ".O...\n"
	str += "O..O.\n"
	str += "O....\n"
	str += "O.O.."

	var grid = life.NewGridFromString(str, "O", ".")

	if grid.Neighbours(0, 0, grid.GetDefault) != 1 {
		t.Errorf("test failed, expected 1 neighbours but got %d", grid.Neighbours(0, 0, grid.GetDefault))
	}

	if grid.Neighbours(0, 0, grid.GetModulo) != 1 {
		t.Errorf("test failed, expected 1 neighbours but got %d", grid.Neighbours(0, 0, grid.GetModulo))
	}
}
