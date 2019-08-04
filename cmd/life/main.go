package main

import (
	"fmt"
	"log"

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
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("cells", 0, 0, maxX-1, maxY-4); err != nil {
		v.Title = "Cells"
		fmt.Fprintln(v, grid.grid.String("O", ".", 2))
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

	go grid.Loop(g)

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
