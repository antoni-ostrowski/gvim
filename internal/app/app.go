package app

import (
	"os"

	"github.com/gdamore/tcell/v3"
)

type App struct {
	ScreenEventChan chan tcell.Event
	Mode            EditorMode
	CommandPrompt   bool
}

type EditorMode interface {
	GetMode() string
	KeyHandler(screen tcell.Screen, app *App, event *tcell.EventKey)
}

type NormalMode struct {
}

func (m *NormalMode) KeyHandler(screen tcell.Screen, app *App, event *tcell.EventKey) {
	switch event.Key() {
	case tcell.KeyCtrlC:
		screen.Fini()
		os.Exit(0)
	case tcell.KeyCtrlV:
		app.Mode = &VisualBlockMode{}
	case tcell.KeyRune:
		str := event.Str()
		DrawDebugStr(screen, str)
		if str == ":" {
			app.Mode = &CommandPromptMode{}
			return
		}
		if str == "i" {
			app.Mode = &InsertMode{}
			return
		}
		if str == "a" {
			app.Mode = &InsertMode{}
			return
		}

		if str == "v" {
			app.Mode = &VisualMode{}
			return
		}
	}
}

func (m *NormalMode) GetMode() string {
	return "normal"
}

type InsertMode struct {
}

func (m *InsertMode) KeyHandler(screen tcell.Screen, app *App, event *tcell.EventKey) {
	switch event.Key() {
	case tcell.KeyCtrlC:
		screen.Fini()
		os.Exit(0)
	case tcell.KeyEscape:
		app.Mode = &NormalMode{}
	}
}

func (m *InsertMode) GetMode() string {
	return "insert"
}

type VisualMode struct {
}

func (m *VisualMode) KeyHandler(screen tcell.Screen, app *App, event *tcell.EventKey) {
	switch event.Key() {
	case tcell.KeyCtrlC:
		screen.Fini()
		os.Exit(0)
	case tcell.KeyEscape:
		app.Mode = &NormalMode{}
	}
}

func (m *VisualMode) GetMode() string {
	return "visual"
}

type VisualBlockMode struct {
}

func (m *VisualBlockMode) KeyHandler(screen tcell.Screen, app *App, event *tcell.EventKey) {
	switch event.Key() {
	case tcell.KeyCtrlC:
		screen.Fini()
		os.Exit(0)
	case tcell.KeyEscape:
		app.Mode = &NormalMode{}
	}
}

func (m *VisualBlockMode) GetMode() string {
	return "visualBlock"
}

type CommandPromptMode struct {
}

func (m *CommandPromptMode) KeyHandler(screen tcell.Screen, app *App, event *tcell.EventKey) {
	switch event.Key() {
	case tcell.KeyCtrlC:
		screen.Fini()
		os.Exit(0)
	case tcell.KeyEscape:
		app.Mode = &NormalMode{}
	}
}

func (m *CommandPromptMode) GetMode() string {
	return "commandPromptMode"
}

func DrawAppState(screen tcell.Screen, app *App) {
	screen.Clear()

	screen.PutStrStyled(0, 0, "test from renderer func - "+app.Mode.GetMode(), tcell.StyleDefault)

	if app.Mode.GetMode() == "commandPromptMode" {
		screen.PutStrStyled(1, 10, "should show command prompt", tcell.StyleDefault)
	}
	screen.Show()
}

func DrawDebugStr(screen tcell.Screen, msg string) {
	w, h := screen.Size()
	screen.PutStrStyled(w-len(msg)-1, h-1, msg, tcell.StyleDefault)
}
