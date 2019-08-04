package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/jroimartin/gocui"
	"gitlab.com/ollybritton/life"
)

// CurrentDirectory gets the path to the directory the current file is inside.
func CurrentDirectory() string {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		panic("Cannot find current directory")
	}

	dir := path.Dir(filename)
	return dir
}

// GetTemplateList gets a list of all the templates.
func GetTemplateList() (templates []string) {
	current := CurrentDirectory()
	examples := path.Join(current, "examples")

	files, err := ioutil.ReadDir(examples)
	if err != nil {
		log.Fatalln(err)
	}

	for _, f := range files {
		templates = append(templates, f.Name())
	}

	return templates
}

// GetTemplate gets the string of the template with a given name.
func GetTemplate(name string) string {
	current := CurrentDirectory()
	examples := path.Join(current, "examples", name)

	fileBytes, err := ioutil.ReadFile(examples)
	if err != nil {
		log.Fatalln(err)
	}

	return string(fileBytes)
}

// GridWrapper is a wrapper for grid which allows it to print and be controlled.
type GridWrapper struct {
	grid  life.Grid
	state string
	mu    sync.Mutex
	gui   *gocui.Gui

	command     string
	editorstate string
	message     string
}

// SetState updates the state of the grid.
func (grid *GridWrapper) SetState(s string) {
	grid.mu.Lock()
	grid.state = s
	grid.mu.Unlock()
}

// GetState gets the state of the grid.
func (grid *GridWrapper) GetState() string {
	return grid.state
}

// Loop handles printing the grid to a writer.
func (grid *GridWrapper) Loop(g *gocui.Gui) {
	for {
		switch grid.GetState() {
		case "animation":
			time.Sleep(time.Millisecond * 200)
			g.Update(func(g *gocui.Gui) error {
				v, err := g.View("cells")
				if err != nil {
					return err
				}

				v.Clear()

				grid.grid.Step(grid.grid.GetDefault)
				fmt.Fprintln(v, grid.grid.String("O", ".", 2))

				return nil
			})

		case "step":
			g.Update(func(g *gocui.Gui) error {
				v, err := g.View("cells")
				if err != nil {
					return err
				}

				v.Clear()

				grid.grid.Step(grid.grid.GetDefault)
				fmt.Fprintln(v, grid.grid.String("O", ".", 2))

				return nil
			})

			grid.SetState("pause")

		default:
			time.Sleep(time.Millisecond * 200)

		}
	}
}

// Update updates the view.
func (grid *GridWrapper) Update(v *gocui.View) {
	switch grid.editorstate {
	case "message":
		v.FgColor = gocui.ColorGreen
		fmt.Fprintf(v, " %v", grid.message)
	}
}

// Edit handles the state of the editor
func (grid *GridWrapper) Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	switch grid.editorstate {
	case "message":
		v.Clear()
		v.FgColor = gocui.ColorWhite
		grid.editorstate = "command"
		v.SetCursor(1, 0)

		grid.Edit(v, key, ch, mod)

	case "command":
		grid.CommandMode(v, key, ch, mod)

	}
}

// CommandMode handles command input into the program.
func (grid *GridWrapper) CommandMode(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	switch {
	case ch != 0 && mod == 0:
		v.EditWrite(ch)

	case key == gocui.KeySpace:
		v.EditWrite(' ')

	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		v.EditDelete(true)

	case key == gocui.KeyDelete:
		v.EditDelete(false)

	case key == gocui.KeyEnter:
		buf := v.Buffer()
		buf = strings.TrimSpace(buf)
		buf = strings.Trim(buf, "/")
		grid.command = buf

		v.Clear()
		v.SetCursor(0, 0)

		message := grid.handleCommand(v)

		grid.message = message
		grid.editorstate = "message"
		grid.Update(v)
	}
}

func (grid *GridWrapper) handleCommand(v *gocui.View) string {
	parts := strings.Split(grid.command, " ")

	if parts[0] == "play" && grid.GetState() != "animation" {
		grid.SetState("animation")
		return "Started animation."
	}

	if parts[0] == "pause" {
		grid.SetState("pause")
		return "Paused animation."
	}

	if parts[0] == "step" || (parts[0] == "" && grid.GetState() != "animation") {
		grid.SetState("step")
		return "Moving forward one generation."
	}

	if parts[0] == "quit" {
		grid.gui.Close()
		os.Exit(0)
	}

	if parts[0] == "clear" || parts[0] == "reset" {
		grid.grid = life.NewGrid(100, 100)
		grid.Update(v)
	}

	if parts[0] == "set" {
		if len(parts) == 1 {
			return "No template specified. Usage: set [template name]"
		}

		template := parts[1]
		template = strings.ToLower(template)

		var templateExists bool

		for _, templateName := range GetTemplateList() {
			if template == templateName {
				templateExists = true
				break
			}
		}

		if !templateExists {
			return "Template with that name doesn't exist. Try 'templates' to see a list."
		}

		grid.grid = life.NewGridFromString(GetTemplate(template), "O", ".")
		grid.grid.Extend(100, 100)
		grid.Update(v)
		return "Template " + template + " set."

	}

	if parts[0] == "templates" {
		return strings.Join(GetTemplateList(), ", ")
	}

	return ""
}
