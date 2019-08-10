package life

import (
	"strings"
)

// Grid represents the set of cells used to make the simulation.
// By default, all cells are the number 0, meaning dead.
// If a cell is set to 1, it is marked as alive.
type Grid struct {
	Cells     [][]int
	Width     int
	Height    int
	tempCells [][]int
}

// NewGrid creates a new grid with the width and height specified.
func NewGrid(width, height int) Grid {
	var grid = Grid{
		Width:  width,
		Height: height,
	}

	for i := 0; i < height; i++ {
		var layer []int

		for j := 0; j < width; j++ {
			layer = append(layer, 0)
		}
		grid.Cells = append(grid.Cells, layer)
	}

	return grid
}

// NewGridFromString creates a new grid using the initial state specified in a string.
func NewGridFromString(str, active, inactive string) Grid {
	var grid Grid

	layers := strings.Split(str, "\n")
	grid = NewGrid(len(layers[0]), len(layers))

	for i, layer := range layers {
		for j, c := range layer {
			if string(c) == active {
				grid.Cells[i][j] = 1
			}
		}
	}

	return grid
}

// String will represent the grid in a human readable way.
func (g *Grid) String(active, inactive string, ratio int) string {
	var str string
	for _, y := range g.Cells {
		for _, x := range y {
			if x == 1 {
				str += strings.Repeat(active, ratio)
			} else {
				str += strings.Repeat(inactive, ratio)
			}
		}

		str += "\n"
	}

	return str
}

// Extend will expand the size of a grid so that it meet the size specified.
// If the size is smaller than the existing one, it will have no effect.
func (g *Grid) Extend(width, height int) {
	addWidth := width - g.Width
	addHeight := height - g.Height

	for y := range g.Cells {
		for i := 0; i < addWidth; i++ {
			g.Cells[y] = append(g.Cells[y], 0)
		}
	}

	for i := 0; i < addHeight; i++ {
		var row []int

		for j := 0; j < width; j++ {
			row = append(row, 0)
		}

		g.Cells = append(g.Cells, row)
	}

	g.Width = width
	g.Height = height
}

// CoordToIndex converts a coordinate, such as (4,4), into the indicies needed to access
// the element in the cell slice.
func (g *Grid) CoordToIndex(x, y int) (int, int) {
	return x, -y + g.Height - 1
}

// IndexToCord coverts a index into the coordinate form.
func (g *Grid) IndexToCord(xIndex, yIndex int) (int, int) {
	return xIndex, -yIndex + g.Height - 1
}

// Get will get the value of a coordinate in the grid.
func (g *Grid) Get(x, y int) int {
	xIndex, yIndex := g.CoordToIndex(x, y)
	return g.Cells[yIndex][xIndex]
}

// GetDefault will get the value of a coordinate in the grid, defaulting to 0
// if it is out of the bounds of the grid.
func (g *Grid) GetDefault(x, y int) int {
	xIndex, yIndex := g.CoordToIndex(x, y)

	if xIndex < 0 || xIndex >= g.Width {
		return 0
	}

	if yIndex < 0 || yIndex >= g.Height {
		return 0
	}

	return g.Cells[yIndex][xIndex]
}

// GetModulo will get the value of a coordinate in the grid, wrapping around the grid if
// the cell is out of bounds.
func (g *Grid) GetModulo(x, y int) int {
	xIndex, yIndex := g.CoordToIndex(x, y)

	xIndex = xIndex % g.Width
	yIndex = yIndex % g.Height

	if xIndex < 0 {
		xIndex = -xIndex
	}

	if yIndex < 0 {
		yIndex = -yIndex
	}

	return g.Cells[yIndex][xIndex]
}

// Set will set the value inside of a coordinate to a value specified.
func (g *Grid) Set(x, y, val int) {
	xIndex, yIndex := g.CoordToIndex(x, y)
	g.Cells[yIndex][xIndex] = val
}

func (g *Grid) setTemp(x, y, val int) {
	xIndex, yIndex := g.CoordToIndex(x, y)
	g.tempCells[yIndex][xIndex] = val
}

// Neighbours will return the number of neighbours around a cell specified by coordinates.
func (g *Grid) Neighbours(x, y int, getter func(x, y int) int) int {
	var count int
	var diffs = [3]int{-1, 0, 1}

	for _, diffX := range diffs {
		for _, diffY := range diffs {
			if diffX == 0 && diffY == 0 {
				continue
			}

			if getter(x+diffX, y+diffY) == 1 {
				count++
			}
		}
	}

	return count
}

// Eval will evaluate a cell at a given coordinate and returns the value that it should be in the next
// generation.
func (g *Grid) Eval(x, y int, getter func(x, y int) int) int {
	neighbours := g.Neighbours(x, y, getter)

	switch {
	case neighbours < 2:
		return 0

	case neighbours > 3:
		return 0

	case neighbours == 3:
		return 1
	}

	return getter(x, y)
}

// Step wil progress the simulation by one generation.
func (g *Grid) Step(getter func(x, y int) int) {
	var temp = NewGrid(g.Width, g.Height)

	for yIndex := range temp.Cells {
		for xIndex := range temp.Cells[yIndex] {

			x, y := g.IndexToCord(xIndex, yIndex)
			temp.Set(x, y, g.Eval(x, y, getter))

		}
	}

	g.Cells = temp.Cells
}
