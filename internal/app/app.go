package app

import (
	editorApi "github.com/antoni-ostrowski/gvim/internal/editor-api"
	"github.com/antoni-ostrowski/gvim/internal/machine"
	"github.com/antoni-ostrowski/gvim/internal/rendering"
	"github.com/gdamore/tcell/v3"
)

type App struct {
	Machine      machine.VimMachine
	UiElements   []rendering.Drawable
	QuitChn      chan struct{}
	Screen       tcell.Screen
	EditorBuffer EditorBuffer
}

var _ editorApi.EditorApi = (*App)(nil)

func (a *App) SendQuitSignal() {
	a.QuitChn <- struct{}{}
}
func (a *App) CurrentMode() editorApi.EditorMode {
	return a.Machine.Mode
}

func (a *App) MoveEditorBuffCursor(amount int, direction editorApi.Direction) {
	currX := a.EditorBuffer.CursorX
	currY := a.EditorBuffer.CursorY
	switch direction.(type) {
	case editorApi.Left:
		if currX != 0 {
			a.EditorBuffer.CursorX = currX - amount
		}
	case editorApi.Right:
		a.EditorBuffer.CursorX = currX + amount
	case editorApi.Up:
		if currY != 0 {
			a.EditorBuffer.CursorY = currY - 1
		}
	case editorApi.Down:
		a.EditorBuffer.CursorY = currY + 1
	}
}

func (a *App) ToggleCommandPrompt(active bool) {
	if active {
		_, h := a.Screen.Size()
		a.UiElements = append(a.UiElements, &CommandPrompt{Input: rendering.TextInput{X: 1, Y: h - 1, CursorPos: 0, Buffer: []rune{}}})
	} else {
		if len(a.UiElements) > 0 {
			index := len(a.UiElements) - 1
			a.UiElements[index] = nil
			a.UiElements = a.UiElements[:index]
		}
	}
}

func DrawAppState(screen tcell.Screen, appState *App) {
	screen.Clear()
	_, h := screen.Size()
	screen.PutStrStyled(0, h-2, appState.Machine.Mode.GetMode(), tcell.StyleDefault)

	appState.EditorBuffer.Draw(screen)

	for _, elem := range appState.UiElements {
		elem.Draw(screen)
	}

	screen.Show()
}

type CommandPrompt struct {
	Input rendering.TextInput
}

var _ rendering.Drawable = (*CommandPrompt)(nil)

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
