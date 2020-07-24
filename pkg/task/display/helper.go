package display

import (
	"fmt"

	"github.com/alexeyco/simpletable"
	"github.com/jroimartin/gocui"
)

type KeyBinding struct {
	Key    string
	Action string
}

var (
	KeyBindings = []KeyBinding{
		{Key: "ctrl-a", Action: "Add new task"},
		{Key: "ctrl-r", Action: "Delete task"},
		{Key: "ctrl-d", Action: "Check task"},
		{Key: "ctrl-u", Action: "Un-Check task"},
		{Key: "ctrl-q", Action: "Close popups"},
		{Key: "tab", Action: "Switch tabs"},
		{Key: "arrow-up", Action: "Up"},
		{Key: "arrow-down", Action: "Down"},
		{Key: "ctrl-c", Action: "Quit"},
	}
)

// Helper Pane
type Helper struct {
	id string
}

// NewHelperPane returns a new Helper pane
func NewHelperPane(id string) Pane {
	return &Helper{id: id}
}

// Layout lays out the Helper pane in gui
func (h *Helper) Layout(gui *gocui.Gui, pos Position) error {
	if v, err := gui.SetView(h.id, pos.X0, pos.Y0, pos.X1, pos.Y1); err != nil {
		if err != nil && err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " Shortcuts "
		v.Wrap = true
		v.Frame = true
		helperTable := simpletable.New()
		helperTable.Header = &simpletable.Header{
			Cells: []*simpletable.Cell{
				{Align: simpletable.AlignCenter, Text: "Key"},
				{Align: simpletable.AlignCenter, Text: "Action"},
			},
		}
		for _, keyBinding := range KeyBindings {
			r := []*simpletable.Cell{
				{Align: simpletable.AlignLeft, Text: keyBinding.Key},
				{Align: simpletable.AlignLeft, Text: keyBinding.Action},
			}
			helperTable.Body.Cells = append(helperTable.Body.Cells, r)
		}

		helperTable.SetStyle(simpletable.StyleCompactLite)
		fmt.Fprintf(v, "%s", helperTable.String())
	}
	return nil
}

// KeyUp action
func (h *Helper) KeyUp(gui *gocui.Gui, v *gocui.View) (interface{}, error) {
	return nil, nil
}

// KeyUp action
func (h *Helper) KeyDown(gui *gocui.Gui, v *gocui.View) (interface{}, error) {
	return nil, nil
}

func (h *Helper) Update(gui *gocui.Gui, op Operation, data interface{}) (interface{}, error) {
	return nil, nil
}
