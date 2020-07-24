package main

import (
	"os/user"
	"path/filepath"
	"time"

	"github.com/jroimartin/gocui"
	"github.com/thejasbabu/track-it/pkg/task"
	"github.com/thejasbabu/track-it/pkg/task/display"
	"github.com/thejasbabu/track-it/util"
)

const (
	// ScreenRefreshInterval determines the screen refresh interval
	ScreenRefreshInterval = 1 * time.Second
	// RootContainerID is the containerID of the root
	RootContainerID = "root"
)

func main() {
	dataPath := dataPath()
	db, err := util.Open(dataPath)
	panicIfErr(err)

	repo := task.NewBadgerRepository(db)
	taskOperator := task.NewOperator(repo, util.UUIDIdentifier{}, util.SystemClock{})

	gui, err := gocui.NewGui(gocui.Output256)
	defer gui.Close()
	panicIfErr(err)

	screen := display.NewScreen(taskOperator, gui)
	gui.SetManagerFunc(screen.Display)
	screen.SetUp(gui)
	if err := gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		panic(err)
	}
}

//TODO: Handle when user's HOME changes
func dataPath() string {
	usr, err := user.Current()
	if err != nil {
		panicIfErr(err)
	}
	return filepath.Join(usr.HomeDir, ".track-it")
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
