package app

import (
	editorApi "github.com/antoni-ostrowski/gvim/internal/editor-api"
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

func (c *CommandPrompt) HandleKey(event *tcell.EventKey, editorApi editorApi.EditorApi) bool {
	switch event.Key() {
	case tcell.KeyEnter:
		if len(c.Input.Buffer) == 0 {
			return true
		}

		if string(c.Input.Buffer[0]) == "q" {
			editorApi.SendQuitSignal()
			return true
		}
		return true
	}
	return c.Input.HandleKey(event, editorApi)
}

func DrawStatusLine(screen tcell.Screen, appState *App) {
	_, h := screen.Size()

	screen.PutStrStyled(0, h-2, appState.Machine.GetMode().GetMode(), tcell.StyleDefault)
}
