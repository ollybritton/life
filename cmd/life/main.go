package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/jroimartin/gocui"
	"gitlab.com/ollybritton/life"
)

var grid GridWrapper

func init() {
	var str string
	str += "..O..\n"
	str += "...O.\n"
	str += ".OOO."

	grid.grid = life.NewGridFromString(str, "O", ".")
	grid.grid.Extend(100, 100)
	grid.SetState("animation")

	grid.loopInterval = 100 * time.Millisecond
	grid.active = "O"
	grid.inactive = "."
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("cells", 0, 0, maxX-1, maxY-4); err != nil {
		v.Title = "Cells"
		fmt.Fprintln(v, grid.grid.String(grid.active, grid.inactive, 2))
	}

	if v, err := g.SetView("commands", 0, maxY-3, maxX-1, maxY-1); err != nil {
		v.Title = "Commands"
		v.Editor = &grid
		v.Editable = true

		v.FgColor = gocui.ColorGreen
		grid.editorstate = "message"
		grid.message = "Type 'quit' to exit."
		grid.Update(v)

		if _, err := g.SetCurrentView("commands"); err != nil {
			return err
		}
	}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	grid.gui = g

	var filename, active, inactive string

	switch len(os.Args) {
	case 2:
		filename = os.Args[1]
		active = "O"
		inactive = "."

	case 3:
		filename = os.Args[1]
		active = "O"
		inactive = "."

	case 4:
		filename = os.Args[1]
		active = os.Args[2]
		inactive = os.Args[3]
	}

	if len(os.Args) >= 2 {
		filebytes, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Fatalln(err)
		}

		grid.grid = life.NewGridFromString(string(filebytes), active, inactive)
		grid.grid.Extend(100, 100)
	}

	go grid.Loop(g)

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
