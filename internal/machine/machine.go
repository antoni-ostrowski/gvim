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
	KeyHandler(event *tcell.EventKey, editorActions editorApi.EditorApi) EditorMode
}

type NormalMode struct{}

func (m *NormalMode) KeyHandler(event *tcell.EventKey, editorActions editorApi.EditorApi) EditorMode {
	return nil
}

func (m *NormalMode) GetMode() string { return "normal" }

type insertMode struct{}

func (m *insertMode) KeyHandler(event *tcell.EventKey, editorActions editorApi.EditorApi) EditorMode {
	return nil
}

func (m *insertMode) GetMode() string { return "insert" }

type visualMode struct{}

func (m *visualMode) KeyHandler(event *tcell.EventKey, editorActions editorApi.EditorApi) EditorMode {
	return nil
}

func (m *visualMode) GetMode() string { return "visual" }

type commandPromptMode struct{}

func (m *commandPromptMode) KeyHandler(event *tcell.EventKey, editorActions editorApi.EditorApi) EditorMode {
	return nil
}

func (m *commandPromptMode) GetMode() string { return "commandPromptMode" }
