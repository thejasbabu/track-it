package display

import "github.com/jroimartin/gocui"

type Pane interface {
	Layout(gui *gocui.Gui, pos Position) error
	KeyUp(g *gocui.Gui, v *gocui.View) (interface{}, error)
	KeyDown(g *gocui.Gui, v *gocui.View) (interface{}, error)
	Update(g *gocui.Gui, op Operation, data interface{}) (interface{}, error)
}

type Operation string

const (
	TaskSelected Operation = "TaskSelected"
	TaskDone     Operation = "TaskDone"
	TaskUnDone   Operation = "TaskUnDone"
	TaskDelete   Operation = "TaskDelete"
	TaskRefresh  Operation = "TaskRefresh"
)
