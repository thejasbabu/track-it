package display

import (
	"fmt"

	"github.com/jroimartin/gocui"
	"github.com/thejasbabu/track-it/util"
)

const (
	Quote = "Track-It!\nTrack your tasks, habits and more!\n\nBy Thejas"
)

// Greeter is used to display greeting quotes
type Greeter struct {
	clock util.Clock
	id    string
}

// NewGreeterPane returns a Greeter Pane
func NewGreeterPane(clock util.Clock, id string) Pane {
	return &Greeter{clock: clock, id: id}
}

// Layout lays out the Greeter pane in gui
func (g *Greeter) Layout(gui *gocui.Gui, pos Position) error {
	if v, err := gui.SetView(g.id, pos.X0, pos.Y0, pos.X1, pos.Y1); err != nil {
		if err != nil && err != gocui.ErrUnknownView {
			return err
		}
		v.Title = title(g.clock)
		v.Wrap = true
		v.Frame = true
		fmt.Fprintf(v, "%s", Quote)
	}
	return nil
}

// KeyUp action
func (g *Greeter) KeyUp(_ *gocui.Gui, v *gocui.View) (interface{}, error) {
	return nil, nil
}

// KeyDown action
func (g *Greeter) KeyDown(_ *gocui.Gui, v *gocui.View) (interface{}, error) {
	return nil, nil
}

func (g *Greeter) Update(gui *gocui.Gui, op Operation, data interface{}) (interface{}, error) {
	return nil, nil
}

func title(clock util.Clock) string {
	t := clock.CurrentTime()
	switch {
	case t.Hour() < 12:
		return " Good morning! "
	case t.Hour() < 17:
		return " Good afternoon! "
	default:
		return " Good evening! "
	}
}
