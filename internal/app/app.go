package app

import (
	"github.com/antoni-ostrowski/gvim/internal/machine"
	"github.com/antoni-ostrowski/gvim/internal/rendering"
	"github.com/gdamore/tcell/v3"
)

type App struct {
	ScreenEventChan chan tcell.Event
	CommandPrompt   CommandPrompt
	Machine         machine.VimMachine
	UiElements      []UiElement
}

type UiElement interface {
	HandleKey(ev *tcell.EventKey)
	Draw(screen tcell.Screen)
}

type CommandPrompt struct {
	Input rendering.TextInput
}

func DrawAppState(screen tcell.Screen, app *App) {
	screen.Clear()
	screen.PutStrStyled(0, 0, "test from renderer func - "+app.Machine.Mode.GetMode(), tcell.StyleDefault)

	if app.Machine.Mode.GetMode() == "commandPromptMode" {
		app.CommandPrompt.Input.Draw(screen)
	}
	screen.Show()
}

func DrawDebugStr(screen tcell.Screen, msg string) {
	w, h := screen.Size()
	screen.PutStrStyled(w-len(msg)-1, h-1, msg, tcell.StyleDefault)
}
