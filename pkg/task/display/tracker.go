package display

import (
	"fmt"
	"strings"
	"time"

	"github.com/jroimartin/gocui"
	"github.com/thejasbabu/track-it/domain"
	"github.com/thejasbabu/track-it/pkg/task"
	"github.com/thejasbabu/track-it/util"
)

// Tracker is used to see the results of the selected task
type Tracker struct {
	operator task.Operator
	id       string
	clock    util.Clock
}

// NewTrackerPane returns a Tracker Pane
func NewTrackerPane(operator task.Operator, id string) Pane {
	return &Tracker{operator: operator, id: id, clock: util.SystemClock{}}
}

// Layout lays out the Task Tracker pane in gui
func (t *Tracker) Layout(gui *gocui.Gui, pos Position) error {
	if v, err := gui.SetView(t.id, pos.X0, pos.Y0, pos.X1, pos.Y1); err != nil {
		if err != nil && err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " Task Status "
		v.Wrap = true
		v.Frame = true
		fmt.Fprintf(v, "%s", Quote)
	}
	return nil
}

// KeyUp action
func (t *Tracker) KeyUp(g *gocui.Gui, v *gocui.View) (interface{}, error) {
	return nil, nil
}

// KeyDown action
func (t *Tracker) KeyDown(g *gocui.Gui, v *gocui.View) (interface{}, error) {
	return nil, nil
}

func (t *Tracker) Update(gui *gocui.Gui, op Operation, data interface{}) (interface{}, error) {
	gui.Update(func(g *gocui.Gui) error {
		v, err := g.View(t.id)
		if err != nil {
			return err
		}
		v.Clear()
		v.Title = " Task Status "
		v.Wrap = true
		v.Frame = true
		task := data.(domain.Task)
		var heading strings.Builder
		var status strings.Builder
		switch task.RepeatInterval {
		case domain.DAILY:
			heading, status = daily(t.clock.CurrentTime(), task)
		case domain.WEEKLY:
			heading, status = weekly(t.clock.CurrentTime(), task)
		case domain.MONTHLY:
			heading, status = monthly(t.clock.CurrentTime(), task)
		case domain.NONE:
			heading, status = singleTimeTask(task)
		}
		fmt.Fprintf(v, fmt.Sprintf("Date  \t\t%s\n", heading.String()))
		fmt.Fprintf(v, fmt.Sprintf("Status\t\t%s\n", status.String()))
		return nil
	})
	return data, nil
}

func daily(currentTime time.Time, task domain.Task) (strings.Builder, strings.Builder) {
	var heading strings.Builder
	var status strings.Builder
	for i := 10; i >= 0; i-- {
		time := currentTime.AddDate(0, 0, -i)
		dateFormat := fmt.Sprintf("%.3s/%d\t\t", time.Month().String(), time.Day())
		heading.WriteString(dateFormat)
		if task.IsComplete(time) {
			fmt.Fprintf(&status, " \u2705\t\t\t\t\t\t")
		} else {
			fmt.Fprintf(&status, " \u2757\t\t\t\t\t\t")
		}
	}
	return heading, status
}

func weekly(currentTime time.Time, task domain.Task) (strings.Builder, strings.Builder) {
	var heading strings.Builder
	var status strings.Builder
	for i := 10; i >= 0; i-- {
		time := currentTime.AddDate(0, 0, -i*7)
		year, week := time.ISOWeek()
		dateFormat := fmt.Sprintf("%d/%d\t\t", year, week)
		heading.WriteString(dateFormat)
		if task.IsComplete(time) {
			fmt.Fprintf(&status, "\u2705\t\t\t\t\t\t\t\t")
		} else {
			fmt.Fprintf(&status, "\u2757\t\t\t\t\t\t\t\t")
		}
	}
	return heading, status
}

func monthly(currentTime time.Time, task domain.Task) (strings.Builder, strings.Builder) {
	var heading strings.Builder
	var status strings.Builder
	for i := 5; i >= 0; i-- {
		time := currentTime.AddDate(0, -i, 0)
		dateFormat := fmt.Sprintf("%.3s\t\t", time.Month().String())
		heading.WriteString(dateFormat)
		if task.IsComplete(time) {
			fmt.Fprintf(&status, "\u2705\t\t\t\t")
		} else {
			fmt.Fprintf(&status, "\u2757\t\t\t\t")
		}
	}
	return heading, status
}

func singleTimeTask(task domain.Task) (strings.Builder, strings.Builder) {
	var heading strings.Builder
	var status strings.Builder
	if len(task.Records) > 0 {
		lastUpdatedTime := task.Records[0].LastUpdatedTime
		dateFormat := fmt.Sprintf("%.3s/%d\t\t", lastUpdatedTime.Month().String(), lastUpdatedTime.Day())
		heading.WriteString(dateFormat)
		if task.IsComplete(lastUpdatedTime) {
			fmt.Fprintf(&status, "\u2705\t\t\t\t")
		} else {
			fmt.Fprintf(&status, "\u2757\t\t\t\t")
		}
	}
	return heading, status
}
