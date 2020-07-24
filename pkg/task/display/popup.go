package display

import (
	"fmt"
	"strings"

	"github.com/jroimartin/gocui"
	"github.com/thejasbabu/track-it/domain"
	"github.com/thejasbabu/track-it/pkg/task"
	"github.com/thejasbabu/track-it/util"
)

type Popup interface {
	Show(gui *gocui.Gui, pos Position, data interface{}) error
	Close(gui *gocui.Gui) error
}

type Confirm struct {
	id string
}

func NewConfirmBox(id string) Popup {
	return &Confirm{id: id}
}

func (c *Confirm) Show(gui *gocui.Gui, pos Position, data interface{}) error {
	if v, err := gui.SetView(c.id, pos.X0, pos.Y0, pos.X1, pos.Y1); err != nil {
		if err != nil && err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " Confirm "
		v.Wrap = true
		v.Frame = true
		message := data.(string)
		fmt.Fprintf(v, "%s", message)
		gui.SetViewOnTop(c.id)
		gui.SetCurrentView(c.id)
	}
	return nil
}

func (c *Confirm) Close(gui *gocui.Gui) error {
	gui.DeleteView(c.id)
	gui.DeleteKeybindings(c.id)
	return nil
}

type TaskForm struct {
	id            string
	operator      task.Operator
	isInputValid  bool
	eventListener chan Event
	clock         util.Clock
}

const (
	TaskDescription = "TaskDescription"
	TaskInterval    = "TaskInterval"
	TaskTags        = "TaskTags"
	TaskSubmit      = "TaskSubmit"
)

func NewTaskForm(id string, operator task.Operator, eventListener chan Event) Popup {
	return &TaskForm{id: id, operator: operator, isInputValid: false, eventListener: eventListener, clock: util.SystemClock{}}
}

func (f *TaskForm) Show(gui *gocui.Gui, pos Position, data interface{}) error {
	if v, err := gui.SetView(f.id, pos.X0, pos.Y0, pos.X1, pos.Y1); err != nil {
		v.Title = "Enter on submit button to save, ctrl-q to quit, tab to switch inputs"
		gui.SelFgColor = gocui.ColorGreen | gocui.AttrBold
		gui.SetViewOnTop(f.id)
		gui.SetCurrentView(f.id)
		task := domain.NewTask(f.clock.CurrentTime())
		if err := gui.SetKeybinding(f.id, gocui.KeyTab, gocui.ModNone, f.switchInputs(&task)); err != nil {
			return err
		}
		if err := gui.SetKeybinding(TaskDescription, gocui.KeyTab, gocui.ModNone, f.switchInputs(&task)); err != nil {
			return err
		}
		if err := gui.SetKeybinding(TaskInterval, gocui.KeyTab, gocui.ModNone, f.switchInputs(&task)); err != nil {
			return err
		}
		if err := gui.SetKeybinding(TaskTags, gocui.KeyTab, gocui.ModNone, f.switchInputs(&task)); err != nil {
			return err
		}
		if err := gui.SetKeybinding(TaskSubmit, gocui.KeyTab, gocui.ModNone, f.switchInputs(&task)); err != nil {
			return err
		}
		if err := gui.SetKeybinding(TaskSubmit, gocui.KeyEnter, gocui.ModNone, f.saveTask(&task)); err != nil {
			return err
		}

		if v, err := gui.SetView(TaskDescription, pos.X0+1, pos.Y0+1, pos.X0+55, pos.Y0+3); err != nil {
			if err != nil && err != gocui.ErrUnknownView {
				return err
			}
			v.Title = " Description "
			v.Wrap = true
			v.Editable = true
			v.Frame = true
		}
		if v, err := gui.SetView(TaskInterval, pos.X0+1, pos.Y0+5, pos.X0+55, pos.Y0+7); err != nil {
			if err != nil && err != gocui.ErrUnknownView {
				return err
			}
			v.Title = "Interval ((N)one, (D)aily, (W)eekly or (M)onthly)"
			v.Wrap = true
			v.Editable = true
			v.Frame = true
		}
		if v, err := gui.SetView(TaskTags, pos.X0+1, pos.Y0+9, pos.X0+55, pos.Y0+11); err != nil {
			if err != nil && err != gocui.ErrUnknownView {
				return err
			}
			v.Title = "Tags (comma seperated)"
			v.Wrap = true
			v.Editable = true
			v.Frame = true
		}
		if v, err := gui.SetView(TaskSubmit, pos.X0+15, pos.Y0+13, pos.X0+25, pos.Y0+15); err != nil {
			if err != nil && err != gocui.ErrUnknownView {
				return err
			}
			v.Wrap = true
			v.Frame = true
			fmt.Fprint(v, "Submit")
		}
	}
	return nil
}

func (f *TaskForm) Close(gui *gocui.Gui) error {
	_, err := gui.View(f.id)
	if err == nil {
		gui.DeleteView(TaskSubmit)
		gui.DeleteKeybindings(TaskSubmit)

		v, _ := gui.View(TaskTags)
		v.Clear()
		gui.DeleteView(TaskTags)

		v, _ = gui.View(TaskInterval)
		v.Clear()
		gui.DeleteView(TaskInterval)

		v, _ = gui.View(TaskDescription)
		v.Clear()
		gui.DeleteView(TaskDescription)

		gui.DeleteView(f.id)
		gui.DeleteKeybindings(f.id)
	}
	return nil
}

func (f *TaskForm) switchInputs(task *domain.Task) func(gui *gocui.Gui, v *gocui.View) error {
	return func(gui *gocui.Gui, v *gocui.View) error {
		if v == nil || v.Name() == f.id {
			_, err := gui.SetCurrentView(TaskDescription)
			return err
		} else if v.Name() == TaskDescription {
			if strings.TrimSpace(v.Buffer()) != "" {
				f.isInputValid = true
				task.Description = strings.TrimSpace(v.Buffer())
			} else {
				f.isInputValid = false
			}
			_, err := gui.SetCurrentView(TaskInterval)
			return err
		} else if v.Name() == TaskInterval {
			task.SetInterval(strings.ToUpper(strings.TrimSpace(v.Buffer())))
			_, err := gui.SetCurrentView(TaskTags)
			return err
		} else if v.Name() == TaskTags {
			task.Tags = strings.Split(v.Buffer(), ",")
			_, err := gui.SetCurrentView(TaskSubmit)
			return err
		} else {
			_, err := gui.SetCurrentView(f.id)
			return err
		}
	}
}

func (f *TaskForm) saveTask(task *domain.Task) func(gui *gocui.Gui, v *gocui.View) error {
	return func(gui *gocui.Gui, v *gocui.View) error {
		if f.isInputValid {
			err := f.operator.Add(*task)
			if err != nil {
				return err
			}
			f.eventListener <- Event{Message: "Refresh"}
			return f.Close(gui)
		} else {
			v, _ := gui.View(TaskSubmit)
			v.BgColor = gocui.ColorRed
			v.SelBgColor = gocui.ColorRed
			return nil
		}
	}
}

func getInterval(interval string) domain.Interval {
	switch {
	case strings.HasPrefix(interval, "D"):
		return domain.DAILY
	case strings.HasPrefix(interval, "W"):
		return domain.WEEKLY
	case strings.HasPrefix(interval, "M"):
		return domain.MONTHLY
	default:
		return domain.NONE
	}
}
