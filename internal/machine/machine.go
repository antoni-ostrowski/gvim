package machine

import (
	editorApi "github.com/antoni-ostrowski/gvim/internal/editor-api"
	"github.com/gdamore/tcell/v3"
)

type VimMachine struct {
	Mode        editorApi.EditorMode
	pendingKeys []rune
}

var _ editorApi.VimMachine = (*VimMachine)(nil)

func (m *VimMachine) Handler(event *tcell.EventKey, api editorApi.EditorApi) {
	// if handler returned diff mode, next key input will get handled by that "next" mode
	if nextMode := m.Mode.KeyHandler(event, api); nextMode != nil {
		m.Mode = nextMode
	}
}
func (m *VimMachine) GetMode() editorApi.EditorMode {
	return m.Mode
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

type insertMode struct{}

var _ editorApi.EditorMode = (*insertMode)(nil)

func (m *insertMode) KeyHandler(event *tcell.EventKey, editorApi editorApi.EditorApi) editorApi.EditorMode {
	if res := handleShared(event, editorApi); res != nil {
		return res
	}

	return nil
}

type visualMode struct{}

var _ editorApi.EditorMode = (*visualMode)(nil)

func (m *visualMode) KeyHandler(event *tcell.EventKey, editorApi editorApi.EditorApi) editorApi.EditorMode {
	if res := handleQuitSignals(event, editorApi); res != nil {
		return res
	}

	if res := handleModeSwitch(event, editorApi); res != nil {
		return res
	}

	return nil
}

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
	buf := api.Buffer()
	switch event.Key() {
	case tcell.KeyLeft:

		buf.MoveCursor(1, editorApi.DirLeft)
	case tcell.KeyRight:
		buf.MoveCursor(1, editorApi.DirRight)
	case tcell.KeyUp:
		buf.MoveCursor(1, editorApi.DirUp)
	case tcell.KeyDown:
		buf.MoveCursor(1, editorApi.DirDown)
	case tcell.KeyRune:
		str := event.Str()
		switch str {
		case "h":
			buf.MoveCursor(1, editorApi.DirLeft)
		case "l":
			buf.MoveCursor(1, editorApi.DirRight)
		case "k":
			buf.MoveCursor(1, editorApi.DirUp)
		case "j":
			buf.MoveCursor(1, editorApi.DirDown)
		}

	}
	return nil
}

func handleModeSwitch(event *tcell.EventKey, api editorApi.EditorApi) editorApi.EditorMode {
	if event.Key() == tcell.KeyEsc {
		api.ToggleCommandPrompt(false)
		return &NormalMode{}
	}

	if _, ok := api.CurrentMode().(*NormalMode); ok && event.Key() == tcell.KeyRune {
		str := event.Str()
		switch str {
		case "i":
			return &insertMode{}
		case "a":
			return &insertMode{}
		case "v":
			return &visualMode{}
		case "V":
			return &visualMode{}
		}
	}

	return nil
}
