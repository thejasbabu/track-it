package display

import (
	"fmt"

	"github.com/jroimartin/gocui"
	"github.com/thejasbabu/track-it/domain"
	"github.com/thejasbabu/track-it/pkg/task"
	"github.com/thejasbabu/track-it/util"
)

// Lister lists all the task added
type Lister struct {
	operator     task.Operator
	id           string
	currentTasks []domain.Task
	clock        util.Clock
}

// NewListPane returns a task pane
func NewListPane(op task.Operator, id string) Pane {
	return &Lister{operator: op, id: id, clock: util.SystemClock{}}
}

// Layout lays out TaskList pane
func (l *Lister) Layout(gui *gocui.Gui, pos Position) error {
	tasks, err := l.operator.GetTasks()
	if err != nil {
		return fmt.Errorf("Error from the task with len:%d:%w", len(tasks), err)
	}
	l.currentTasks = tasks
	if v, err := gui.SetView(l.id, pos.X0, pos.Y0, pos.X1, pos.Y1); err != nil {
		gui.SetCurrentView(l.id)
		renderList(v, tasks, l.clock)
	}

	return nil
}

// KeyUp action
func (l *Lister) KeyUp(gui *gocui.Gui, v *gocui.View) (interface{}, error) {
	cx, cy := v.Cursor()
	switch {
	case cy == 0:
		return domain.Task{}, nil
	default:
		moveUp(v, cx, cy)
		return l.currentTasks[cy-1], nil
	}
}

func (l *Lister) KeyDown(gui *gocui.Gui, v *gocui.View) (interface{}, error) {
	cx, cy := v.Cursor()
	totalTasks := len(l.currentTasks)
	switch {
	case cy == (totalTasks - 1):
		return l.currentTasks[cy], nil
	default:
		moveDown(v, cx, cy)
		return l.currentTasks[cy+1], nil
	}
}

func (l *Lister) Update(gui *gocui.Gui, op Operation, data interface{}) (interface{}, error) {
	v, _ := gui.View(l.id)
	_, cy := v.Cursor()
	switch {
	case op == TaskRefresh:
		err := l.refreshTaskList(gui)
		if err != nil {
			return nil, err
		}
		if len(l.currentTasks) > 0 {
			v.SetCursor(0, 1)
			return l.currentTasks[0], nil
		}

	case cy < 0 || cy >= len(l.currentTasks):
		return domain.Task{}, nil

	case op == TaskDone:
		err := l.operator.MarkAsDone(l.currentTasks[cy])
		if err != nil {
			return nil, err
		}
		err = l.refreshTaskList(gui)
		if err != nil {
			return nil, err
		}
		return l.currentTasks[cy], nil

	case op == TaskUnDone:
		err := l.operator.MarkAsUnDone(l.currentTasks[cy])
		if err != nil {
			return nil, err
		}
		err = l.refreshTaskList(gui)
		if err != nil {
			return nil, err
		}
		return l.currentTasks[cy], nil

	case op == TaskDelete:
		err := l.operator.DeleteTask(l.currentTasks[cy])
		if err != nil {
			return nil, err
		}
		err = l.refreshTaskList(gui)
		if err != nil {
			return nil, err
		}
		if len(l.currentTasks) > 0 {
			v.SetCursor(0, 1)
			return l.currentTasks[0], nil
		}
	}
	return domain.Task{}, nil
}

func (l *Lister) refreshTaskList(gui *gocui.Gui) error {
	tasks, err := l.operator.GetTasks()
	if err != nil {
		return fmt.Errorf("Error from the task with len:%d:%w", len(tasks), err)
	}
	l.currentTasks = tasks
	gui.Update(func(g *gocui.Gui) error {
		v, err := g.View(l.id)
		if err != nil {
			return err
		}
		v.Clear()
		renderList(v, l.currentTasks, l.clock)
		return nil
	})
	return nil
}

func renderList(v *gocui.View, taskList []domain.Task, clock util.Clock) {
	v.Highlight = true
	v.Frame = true
	v.SelBgColor = gocui.ColorGreen
	v.SelFgColor = gocui.ColorBlack
	var completed int
	for _, task := range taskList {
		if task.IsComplete(clock.CurrentTime()) {
			completed++
			fmt.Fprintf(v, "\u2705   %s\n", task.Description)
		} else {
			fmt.Fprintf(v, "\u2757   %s\n", task.Description)
		}
	}
	v.Title = fmt.Sprintf(" Tasks [Total: %d, completed: %d %% ]", len(taskList), pct(completed, len(taskList)))
}

func moveDown(v *gocui.View, x, y int) error {
	if err := v.SetCursor(x, y+1); err != nil {
		ox, oy := v.Origin()
		if err := v.SetOrigin(ox, oy+1); err != nil {
			return err
		}
	}
	return nil
}

func moveUp(v *gocui.View, x, y int) error {
	if err := v.SetCursor(x, y-1); err != nil {
		ox, oy := v.Origin()
		if err := v.SetOrigin(ox, oy-1); err != nil {
			return err
		}
	}
	return nil
}

func checkForEmpty(input string) string {
	if len(input) == 0 || input == "" {
		return "NA"
	}
	return input
}

func pct(a int, b int) int {
	hundred := 100
	if b == 0 {
		return hundred
	}
	return (hundred * a) / b
}
