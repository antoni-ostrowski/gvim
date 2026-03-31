package machine

import (
	editorApi "github.com/antoni-ostrowski/gvim/internal/editor-api"
	"github.com/gdamore/tcell/v3"
)

type VimMachine struct {
	Mode        EditorMode
	pendingKeys []rune
}

func (m *VimMachine) Handler(event *tcell.EventKey, editorApi editorApi.EditorApi) {
	// if handler returned diff mode, next key input will get handled by that "next" mode
	if nextMode := m.Mode.KeyHandler(event, editorApi); nextMode != nil {
		m.Mode = nextMode
	}
}

type EditorMode interface {
	GetMode() string
	KeyHandler(event *tcell.EventKey, editorApi editorApi.EditorApi) EditorMode
}

type NormalMode struct{}

var _ EditorMode = (*NormalMode)(nil)

func (m *NormalMode) KeyHandler(event *tcell.EventKey, editorApi editorApi.EditorApi) EditorMode {
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

func (m *NormalMode) GetMode() string { return "normal" }

type insertMode struct{}

var _ EditorMode = (*insertMode)(nil)

func (m *insertMode) KeyHandler(event *tcell.EventKey, editorApi editorApi.EditorApi) EditorMode {
	if res := handleShared(event, editorApi); res != nil {
		return res
	}

	return nil
}

func (m *insertMode) GetMode() string { return "insert" }

type visualMode struct{}

var _ EditorMode = (*visualMode)(nil)

func (m *visualMode) KeyHandler(event *tcell.EventKey, editorApi editorApi.EditorApi) EditorMode {
	if res := handleShared(event, editorApi); res != nil {
		return res
	}

	return nil
}

func (m *visualMode) GetMode() string { return "visual" }

func handleShared(event *tcell.EventKey, editorApi editorApi.EditorApi) EditorMode {
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

func handleQuitSignals(event *tcell.EventKey, editorApi editorApi.EditorApi) EditorMode {
	switch event.Key() {
	case tcell.KeyCtrlC:
		editorApi.SendQuitSignal()
		return nil
	}
	return nil
}

func handleMovement(event *tcell.EventKey, editorApi editorApi.EditorApi) EditorMode {
	switch event.Key() {
	}
	return nil
}

func handleModeSwitch(event *tcell.EventKey, editorApi editorApi.EditorApi) EditorMode {
	if event.Key() == tcell.KeyEsc {
		editorApi.ToggleCommandPrompt(false)
		return &NormalMode{}
	}

	if event.Key() == tcell.KeyRune {
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
