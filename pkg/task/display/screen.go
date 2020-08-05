package display

import (
	"fmt"

	"github.com/jroimartin/gocui"
	"github.com/thejasbabu/track-it/domain"
	"github.com/thejasbabu/track-it/pkg/task"
	"github.com/thejasbabu/track-it/util"
)

const (
	GreeterPaneID     = "Greeter"
	TaskListPaneID    = "TaskList"
	ConfirmBoxPaneID  = "ConfirmBox"
	TaskFormPaneID    = "TaskForm"
	HelperPaneID      = "Helper"
	TaskDefnPaneID    = "TaskDefn"
	TaskTrackerPaneID = "TaskStatus"
)

type Event struct {
	Message string
	Error   error
}

// Screen holds all the widgets that is shown in the screen
type Screen struct {
	panes         map[string]Pane
	popups        map[string]Popup
	positions     map[string]Position
	EventListener chan Event
}

// NewScreen returns a new screen obj
func NewScreen(operator task.Operator, gui *gocui.Gui, client *util.RedditClient) Screen {
	panes := make(map[string]Pane)
	positions := make(map[string]Position)
	popups := make(map[string]Popup)
	maxX, maxY := gui.Size()
	relX1, relX2 := partitionX(gui)
	relY1, relY2 := partitionY(gui)
	eventChan := make(chan Event)

	panes[GreeterPaneID] = NewGreeterPane(util.SystemClock{}, GreeterPaneID, client)
	positions[GreeterPaneID] = Position{X0: 0, Y0: 0, X1: relX1 - 1, Y1: relY1 - 1}
	panes[HelperPaneID] = NewHelperPane(HelperPaneID)
	positions[HelperPaneID] = Position{X0: 0, Y0: relY1, X1: relX1 - 1, Y1: maxY - 1}
	panes[TaskListPaneID] = NewListPane(operator, TaskListPaneID)
	positions[TaskListPaneID] = Position{X0: relX1, Y0: 0, X1: relX2 - 1, Y1: relY2 - 1}
	panes[TaskDefnPaneID] = NewDefinerPane(operator, TaskDefnPaneID)
	positions[TaskDefnPaneID] = Position{X0: relX2, Y0: 0, X1: maxX - 1, Y1: relY2 - 1}
	panes[TaskTrackerPaneID] = NewTrackerPane(operator, TaskTrackerPaneID)
	positions[TaskTrackerPaneID] = Position{X0: relX1, Y0: relY2, X1: maxX - 1, Y1: maxY - 1}

	var width, height = maxX / 4, maxY / 5
	popups[ConfirmBoxPaneID] = NewConfirmBox(ConfirmBoxPaneID)
	positions[ConfirmBoxPaneID] = Position{X0: maxX/2 - width/2 - 1, Y0: maxY/2 - height/2 - 1, X1: maxX/2 + width/2 + 1, Y1: maxY/2 + height/2 + 1}
	popups[TaskFormPaneID] = NewTaskForm(TaskFormPaneID, operator, eventChan)
	positions[TaskFormPaneID] = Position{X0: maxX/2 - width - 1, Y0: maxY/2 - height - 1, X1: maxX/2 + width, Y1: maxY/2 + height}
	return Screen{panes: panes, positions: positions, EventListener: eventChan, popups: popups}
}

// Position represents a 2-D position in the screen
type Position struct {
	X0 int
	Y0 int
	X1 int
	Y1 int
}

// Display displays all the widgets in the screen
func (s Screen) Display(gui *gocui.Gui) error {
	for paneID, pane := range s.panes {
		err := pane.Layout(gui, s.positions[paneID])
		if err != nil {
			x, y := gui.Size()
			return fmt.Errorf("Error %s, Pane Map: %#v, X: %d, Y: %d", err.Error(), s.positions, x, y)
		}
	}
	return nil
}

func (s Screen) SetUp(gui *gocui.Gui) {
	gui.SelFgColor = gocui.ColorGreen | gocui.AttrBold
	gui.BgColor = gocui.ColorDefault
	gui.Highlight = true
	gui.Mouse = true
	gui.InputEsc = false

	if err := gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, s.quit); err != nil {
		s.EventListener <- Event{Error: err, Message: "quit event"}
	}
	if err := gui.SetKeybinding(TaskListPaneID, gocui.KeyArrowDown, gocui.ModNone, s.down); err != nil {
		s.EventListener <- Event{Error: err, Message: "arrow down event"}
	}
	if err := gui.SetKeybinding(TaskListPaneID, gocui.KeyArrowUp, gocui.ModNone, s.up); err != nil {
		s.EventListener <- Event{Error: err, Message: "arrow up event"}
	}
	if err := gui.SetKeybinding(TaskListPaneID, gocui.KeyCtrlD, gocui.ModNone, s.done); err != nil {
		s.EventListener <- Event{Error: err, Message: "ctrl-D event"}
	}
	if err := gui.SetKeybinding(TaskListPaneID, gocui.KeyCtrlU, gocui.ModNone, s.undone); err != nil {
		s.EventListener <- Event{Error: err, Message: "ctrl-U event"}
	}
	if err := gui.SetKeybinding(TaskListPaneID, gocui.KeyCtrlR, gocui.ModNone, s.remove); err != nil {
		s.EventListener <- Event{Error: err, Message: "ctrl-R event"}
	}
	if err := gui.SetKeybinding("", gocui.KeyCtrlA, gocui.ModNone, s.add); err != nil {
		s.EventListener <- Event{Error: err, Message: "ctrl-A event"}
	}
	if err := gui.SetKeybinding("", gocui.KeyCtrlQ, gocui.ModNone, s.closeAnyPopUp); err != nil {
		s.EventListener <- Event{Error: err, Message: "ctrl-A event"}
	}

	go func() {
		for event := range s.EventListener {
			if event.Error != nil {
				panic(event.Error)
			} else if event.Message == "Refresh" {
				s.refresh(gui)
			}
			// TODO: Handle logging events
		}
	}()
}

func (s Screen) up(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		switch v.Name() {
		case TaskListPaneID:
			item, err := s.panes[TaskListPaneID].KeyUp(g, v)
			if err != nil {
				return err
			}
			task := item.(domain.Task)
			s.panes[TaskDefnPaneID].Update(g, TaskSelected, task)
			s.panes[TaskTrackerPaneID].Update(g, TaskSelected, task)
		}
	}
	return nil
}

func (s Screen) down(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		switch v.Name() {
		case TaskListPaneID:
			item, err := s.panes[TaskListPaneID].KeyDown(g, v)
			if err != nil {
				return err
			}
			task := item.(domain.Task)
			s.panes[TaskDefnPaneID].Update(g, TaskSelected, task)
			s.panes[TaskTrackerPaneID].Update(g, TaskSelected, task)
		}
	}
	return nil
}

func (s Screen) done(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		item, err := s.panes[TaskListPaneID].Update(g, TaskDone, nil)
		if err != nil {
			return err
		}
		task := item.(domain.Task)
		s.panes[TaskDefnPaneID].Update(g, TaskSelected, task)
		s.panes[TaskTrackerPaneID].Update(g, TaskSelected, task)
	}
	return nil
}

func (s Screen) undone(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		item, err := s.panes[TaskListPaneID].Update(g, TaskUnDone, nil)
		if err != nil {
			return err
		}
		task := item.(domain.Task)
		s.panes[TaskDefnPaneID].Update(g, TaskSelected, task)
		s.panes[TaskTrackerPaneID].Update(g, TaskSelected, task)
	}
	return nil
}

func (s Screen) add(gui *gocui.Gui, v *gocui.View) error {
	if v != nil {
		s.popups[TaskFormPaneID].Show(gui, s.positions[TaskFormPaneID], nil)
		if err := gui.SetKeybinding(TaskFormPaneID, gocui.KeyCtrlQ, gocui.ModNone, s.closePopUp(TaskFormPaneID)); err != nil {
			s.EventListener <- Event{Error: err, Message: "delete task event"}
		}
	}
	return nil
}

func (s Screen) remove(gui *gocui.Gui, v *gocui.View) error {
	if v != nil {
		s.popups[ConfirmBoxPaneID].Show(gui, s.positions[ConfirmBoxPaneID], "Press ENTER to confirm deletion of selected task, ctrl-q to close")
		if err := gui.SetKeybinding(ConfirmBoxPaneID, gocui.KeyEnter, gocui.ModNone, s.deleteTask); err != nil {
			s.EventListener <- Event{Error: err, Message: "delete task event"}
		}
		if err := gui.SetKeybinding(ConfirmBoxPaneID, gocui.KeyCtrlQ, gocui.ModNone, s.closePopUp(ConfirmBoxPaneID)); err != nil {
			s.EventListener <- Event{Error: err, Message: "delete task event"}
		}
	}
	return nil
}

func (s Screen) deleteTask(gui *gocui.Gui, v *gocui.View) error {
	item, _ := s.panes[TaskListPaneID].Update(gui, TaskDelete, nil)
	task := item.(domain.Task)
	s.panes[TaskDefnPaneID].Update(gui, TaskSelected, task)
	s.panes[TaskTrackerPaneID].Update(gui, TaskSelected, task)
	s.popups[ConfirmBoxPaneID].Close(gui)
	gui.SetCurrentView(TaskListPaneID)
	gui.Cursor = true
	return nil
}

func (s Screen) refresh(gui *gocui.Gui) {
	item, _ := s.panes[TaskListPaneID].Update(gui, TaskRefresh, nil)
	task := item.(domain.Task)
	s.panes[TaskDefnPaneID].Update(gui, TaskSelected, task)
	s.panes[TaskTrackerPaneID].Update(gui, TaskSelected, task)
	for _, popup := range s.popups {
		popup.Close(gui)
	}
	gui.Cursor = true
	gui.SetCurrentView(TaskListPaneID)
}

func (s Screen) closeAnyPopUp(gui *gocui.Gui, v *gocui.View) error {
	for _, popUp := range s.popups {
		popUp.Close(gui)
	}

	gui.SetCurrentView(TaskListPaneID)
	gui.Cursor = true
	return nil
}

func (s Screen) closePopUp(viewName string) func(gui *gocui.Gui, v *gocui.View) error {
	return func(gui *gocui.Gui, v *gocui.View) error {
		s.popups[viewName].Close(gui)
		gui.SetCurrentView(TaskListPaneID)
		gui.Cursor = true
		return nil
	}
}

func (s Screen) quit(g *gocui.Gui, v *gocui.View) error {
	s.EventListener <- Event{Message: "quiting"}
	return gocui.ErrQuit
}

func partitionX(gui *gocui.Gui) (int, int) {
	x, _ := gui.Size()
	return (x * 2) / 10, (x * 80) / 100
}

func partitionY(gui *gocui.Gui) (int, int) {
	_, y := gui.Size()
	return (y * 5) / 10, (y * 70) / 100
}
