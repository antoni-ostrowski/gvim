package app

import (
	utils "github.com/antoni-ostrowski/gvim/internal"
	editorApi "github.com/antoni-ostrowski/gvim/internal/editor-api"
	"github.com/antoni-ostrowski/gvim/internal/machine"
	"github.com/antoni-ostrowski/gvim/internal/rendering"
	"github.com/gdamore/tcell/v3"
)

type CommandPrompt struct {
	Input rendering.TextInput
}

var _ editorApi.UiElement = (*CommandPrompt)(nil)

func (c *CommandPrompt) Draw(screen tcell.Screen) {
	c.Input.Draw(screen)
}

func (c *CommandPrompt) HandleKey(event *tcell.EventKey, api editorApi.EditorApi) bool {
	switch event.Key() {
	case tcell.KeyEnter:
		if len(c.Input.Buffer) == 0 {
			return true
		}

		if string(c.Input.Buffer[0]) == "q" {
			api.SendQuitSignal()
			return true
		}

		if string(c.Input.Buffer[0]) == "w" {
			err := api.WriteFile()
			if err != nil {
				utils.Debuglog("err writing file %v", err)
			}
			api.ToggleCommandPrompt(false)
			return true
		}
		return true
	}
	return c.Input.HandleKey(event, api)
}

func DrawStatusLine(screen tcell.Screen, appState *App) {
	_, h := screen.Size()

	screen.PutStrStyled(0, h-2, GetCurrentEditorModeName(appState), tcell.StyleDefault)
}

func GetCurrentEditorModeName(appState *App) string {
	switch appState.Machine.GetMode().(type) {
	case *machine.NormalMode:
		return "NORMAL"
	case *machine.InsertMode:
		return "INSERT"
	case *machine.VisualMode:
		return "VISUAL"
	}
	return "UNKNOWN"
}
