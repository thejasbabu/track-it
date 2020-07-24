package display

import (
	"fmt"
	"strings"

	"github.com/jroimartin/gocui"
	"github.com/thejasbabu/track-it/domain"
	"github.com/thejasbabu/track-it/pkg/task"
)

// Definer is used to define the selected task
type Definer struct {
	operator task.Operator
	id       string
}

// NewDefinerPane returns a Task Definer Pane
func NewDefinerPane(operator task.Operator, id string) Pane {
	return &Definer{operator: operator, id: id}
}

// Layout lays out the Task Definer pane in gui
func (d *Definer) Layout(gui *gocui.Gui, pos Position) error {
	if v, err := gui.SetView(d.id, pos.X0, pos.Y0, pos.X1, pos.Y1); err != nil {
		if err != nil && err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " Select task to view more "
		v.Wrap = true
		v.Frame = true
		fmt.Fprintf(v, "%s", Quote)
	}
	return nil
}

// KeyUp action
func (d *Definer) KeyUp(gui *gocui.Gui, v *gocui.View) (interface{}, error) {
	return nil, nil
}

// KeyDown action
func (d *Definer) KeyDown(gui *gocui.Gui, v *gocui.View) (interface{}, error) {
	return nil, nil
}

func (d *Definer) Update(gui *gocui.Gui, op Operation, data interface{}) (interface{}, error) {
	gui.Update(func(g *gocui.Gui) error {
		v, err := g.View(d.id)
		if err != nil {
			return err
		}
		v.Clear()
		v.Title = " Task Statement "
		v.Wrap = true
		v.Frame = true
		task := data.(domain.Task)
		fmt.Fprintf(v, "Desc: %s\n", task.Description)
		fmt.Fprint(v, "\n")
		fmt.Fprintf(v, "Interval: %s\n", task.RepeatInterval)
		fmt.Fprint(v, "\n")
		fmt.Fprintf(v, "Tags: %s\n", strings.Join(task.Tags, ","))
		fmt.Fprint(v, "\n")
		fmt.Fprintf(v, "Streak: %d\n", task.Streak)
		return nil
	})
	return data, nil
}
