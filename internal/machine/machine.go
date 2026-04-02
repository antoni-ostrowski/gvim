package machine

import (
	editorApi "github.com/antoni-ostrowski/gvim/internal/editor-api"
	"github.com/gdamore/tcell/v3"
)

type VimMachine struct {
	Mode        editorApi.EditorMode
	pendingKeys []rune
}

func (m *VimMachine) Handler(event *tcell.EventKey, editorApi editorApi.EditorApi) {
	// if handler returned diff mode, next key input will get handled by that "next" mode
	if nextMode := m.Mode.KeyHandler(event, editorApi); nextMode != nil {
		m.Mode = nextMode
	}
}

type NormalMode struct{}

var _ editorApi.EditorMode = (*NormalMode)(nil)

func (m *NormalMode) KeyHandler(event *tcell.EventKey, editorApi editorApi.EditorApi) editorApi.EditorMode {
	if res := handleShared(event, editorApi); res != nil {
		return res
	}

	if event.Key() == tcell.KeyRune {
		switch event.Str() {
		case ":":
			editorApi.ToggleCommandPrompt(true)
		}
	}

	return nil
}

func (m *NormalMode) GetMode() string { return "NORMAL" }

type insertMode struct{}

var _ editorApi.EditorMode = (*insertMode)(nil)

func (m *insertMode) KeyHandler(event *tcell.EventKey, editorApi editorApi.EditorApi) editorApi.EditorMode {
	if res := handleShared(event, editorApi); res != nil {
		return res
	}

	return nil
}

func (m *insertMode) GetMode() string { return "INSERT" }

type visualMode struct{}

var _ editorApi.EditorMode = (*visualMode)(nil)

func (m *visualMode) KeyHandler(event *tcell.EventKey, editorApi editorApi.EditorApi) editorApi.EditorMode {
	if res := handleShared(event, editorApi); res != nil {
		return res
	}

	return nil
}

func (m *visualMode) GetMode() string { return "VISUAL" }

func handleShared(event *tcell.EventKey, editorApi editorApi.EditorApi) editorApi.EditorMode {
	if res := handleQuitSignals(event, editorApi); res != nil {
		return res
	}

	if res := handleMovement(event, editorApi); res != nil {
		return res
	}

	if res := handleModeSwitch(event, editorApi); res != nil {
		return res
	}
	return nil
}

func handleQuitSignals(event *tcell.EventKey, editorApi editorApi.EditorApi) editorApi.EditorMode {
	switch event.Key() {
	case tcell.KeyCtrlC:
		editorApi.SendQuitSignal()
		return nil
	}
	return nil
}

func handleMovement(event *tcell.EventKey, api editorApi.EditorApi) editorApi.EditorMode {
	switch event.Key() {
	case tcell.KeyLeft:
		api.MoveEditorBuffCursor(1, editorApi.Left{})
	case tcell.KeyRight:
		api.MoveEditorBuffCursor(1, editorApi.Right{})
	case tcell.KeyUp:
		api.MoveEditorBuffCursor(1, editorApi.Up{})
	case tcell.KeyDown:
		api.MoveEditorBuffCursor(1, editorApi.Down{})
	case tcell.KeyRune:
		str := event.Str()
		switch str {
		case "h":
			api.MoveEditorBuffCursor(1, editorApi.Left{})
		case "l":
			api.MoveEditorBuffCursor(1, editorApi.Right{})
		case "k":
			api.MoveEditorBuffCursor(1, editorApi.Up{})
		case "j":
			api.MoveEditorBuffCursor(1, editorApi.Down{})
		}

	}
	return nil
}

func handleModeSwitch(event *tcell.EventKey, api editorApi.EditorApi) editorApi.EditorMode {

	if event.Key() == tcell.KeyEsc {
		api.ToggleCommandPrompt(false)
		return &NormalMode{}
	}

	if api.CurrentMode().GetMode() == "NORMAL" && event.Key() == tcell.KeyRune {
		str := event.Str()
		switch str {
		case "i":
			return &insertMode{}
		case "v":
			return &visualMode{}
		case "V":
			return &visualMode{}
		}
	}

	return nil
}
