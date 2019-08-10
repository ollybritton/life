# Life
Life is a Go package for implementing Conway's Game of Life. 

![Example of `life` command](demo/demo.gif)

- [Life](#life)
  - [Download](#download)
  - [Basic Usage](#basic-usage)
    - [Creating a Grid](#creating-a-grid)
      - [`NewGrid` Example](#newgrid-example)
      - [`NewGridFromString` Example](#newgridfromstring-example)
    - [Printing](#printing)
    - [Setting/Changing values](#settingchanging-values)
      - [Converting between coordinates and indicies](#converting-between-coordinates-and-indicies)
      - [Setting cell values](#setting-cell-values)
      - [Getting cell values](#getting-cell-values)
    - [Getting cell information](#getting-cell-information)
    - [Stepping](#stepping)
  - [This documentation is shit, give me examples](#this-documentation-is-shit-give-me-examples)

## Download
```bash
# For the package and cli
go get github.com/ollybritton/life/...

# For just the package
go get github.com/ollybritton/life
```

## Basic Usage
For more documentation, see [GoDoc](https://godoc.org/github.com/ollybritton/life).

The basic usage of the API follows these basic steps:

1. Create a `Grid`.
2. Set values inside the grid.
3. Advance a generation using `.Step()`

### Creating a Grid
There are two ways you can create a grid. You can either use the simple `life.NewGrid(width, height)` function, which creates an empty grid of the given dimensions, or `life.NewGridFromString(str, active, inactive)` which will create a grid from a given string..

#### `NewGrid` Example
```go
grid := life.NewGrid(20, 20)
```

#### `NewGridFromString` Example
```go
var input string
input += "...O.\n"
input += "....O\n"
input += "..OOO"

grid := life.NewGridFromString(input, "O", ".")
```

It is common that you have a grid made from an input string, but it is the wrong dimensions. For example, you may input a glider that is only 3x3, but want the grid to be 20x20. To do this, use `grid.Extend`:

```go
grid.Extend(20,20)
```

### Printing
You can get the string representation of a grid using the `grid.String(active, inactive, ration)` function.

`active` and `inactive` control which characters should be used to print active and inactive cells. `ratio` is the amount of times to print each character. For example, `ratio=2` will mean that each character is printed twice. This is used because in most cases it will appear squashed with a ratio of 1.

Consider this:

```go
var input string
input += "...O.\n"
input += "....O\n"
input += "..OOO"

grid := life.NewGridFromString(input, "O", ".")

fmt.Println(grid.String("O", ".", 2))
// Outputs
// ......OO..
// ........OO
// ....OOOOOO

// Note the difference between
// ...O.
// ....O
// ..OOO
```

### Setting/Changing values
Values are set using a coordinate system and not array indicies like you would expect. For example, consider this configuration.

```
[. . . . . . . . . .]
[. . O . . . . . . .]
[. . . O . . . . . .]
[. O O O . . . . . .]
[. . . . . . . . . .]
```

If the cells were represented as a multidimensional array, you would be able to access the top-left element using `[0][0]`. However, when setting and modifing cell values, you have to use an **x** and **y** coordinate, starting at `(0,0)` from the bottom left.

This means that the initial example of `[0][0]` becomes `(0, 4)`.

#### Converting between coordinates and indicies
Two helper functions are defined in order to facilitate this:

```go
grid.IndexToCoord(xIndex, yIndex) // => (x, y)
grid.IndexToCoord(x, y) // => (xIndex, yIndex)
```

#### Setting cell values
```go
grid.Set(0, 0, 1) // Sets the bottom left cell to active.
grid.Set(5, 5, 0) // Sets the cell located at (5, 5) to inactive.
```

#### Getting cell values
Three functions exist for getting cell values, which all handle out-of-bound lookups differently.

* `grid.Get(x, y)`
  Panic if the lookup is out of bounds.

* `grid.GetDefault(x, y)`
  Default to inactive (0) if the lookup is out of bounds.

* `grid.GetModulo(x, y)`
  **Currently Broken**
  Wrap around to the opposite side

For example:
```go
var input string
input += "...O.\n"
input += "....O\n"
input += "..OOO"

grid := life.NewGridFromString(input, "O", ".")
grid.Get(2, 0) // Returns 1.
```

### Getting cell information
To find out the number of neighbours around a cell, use `grid.Neighbours(x, y, getter)`, where `getter` is one of the above functions for the lookup of a cell.

```go
var input string
input += "...O.\n"
input += "....O\n"
input += "..OOO"

grid := life.NewGridFromString(input, "O", ".")
grid.Neighbours(2, 1, grid.GetDefault) // Returns 3
```

To find out if the cell will be alive in the next generation, use `grid.Eval(x, y)`

```go 
var input string
input += "...O.\n"
input += "....O\n"
input += "..OOO"

grid := life.NewGridFromString(input, "O", ".")
grid.Eval(2, 1) // Returns 0 (inactive)
```

### Stepping
To move the entire grid one generation forward, use `grid.Step(getter)`, where `getter` is one of the `life.Get, life.GetDefault` or `life.GetModulo` functions.

```go
var input string
input += "...O...\n"
input += "....O..\n"
input += "..OOO..\n"
input += "......."

grid := life.NewGridFromString(input, "O", ".")
grid.Step(grid.GetDefault) // Moves the above grid one generation forward.
```

## This documentation is shit, give me examples
I know it is. What makes it even worse is that the only example I can give is the command line tool, located in `cmd/life`. What is even worse is it's a terrible use of the library and follows bascially zero Go conventions. Good luck. You're probably better off implementing it yourself.