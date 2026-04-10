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

func (m *VimMachine) Handler(event *tcell.EventKey, buf editorApi.EditorBuffer) {
	modeSwitchHandler := func() editorApi.EditorMode {
		switch m.Mode.(type) {
		case *NormalMode:
			if event.Key() == tcell.KeyRune {
				str := event.Str()
				switch str {
				case "i":
					return &InsertMode{}
				case "a":
					buf.MoveCursor(1, editorApi.DirRight)
					return &InsertMode{}
				case "A":
					buf.JumpToLineEnd()
					return &InsertMode{}
				case "v":
					return &VisualMode{}
				case "V":
					return &VisualMode{}
				}
			}
		case *InsertMode:
			if event.Key() == tcell.KeyEsc {
				return &NormalMode{}
			}
		case *VisualMode:
			if event.Key() == tcell.KeyEsc {
				return &NormalMode{}
			}
		}
		return nil
	}

	// if input was changing mode, we change mode and return early
	if nextMode := modeSwitchHandler(); nextMode != nil {
		m.Mode = nextMode
		return
	}

	m.Mode.KeyHandler(event, buf)
}
func (m *VimMachine) GetMode() editorApi.EditorMode {
	return m.Mode
}

type NormalMode struct{}

var _ editorApi.EditorMode = (*NormalMode)(nil)

func (m *NormalMode) KeyHandler(event *tcell.EventKey, buf editorApi.EditorBuffer) {
	handleMovement(event, buf)
	handleSharedKeys(event, buf)

	if event.Key() == tcell.KeyEnter {
		buf.MoveCursor(1, editorApi.DirDown)
		return
	}

}

type InsertMode struct{}

var _ editorApi.EditorMode = (*InsertMode)(nil)

func (m *InsertMode) KeyHandler(event *tcell.EventKey, buf editorApi.EditorBuffer) {
	handleSharedKeys(event, buf)

	if event.Key() == tcell.KeyBackspace {
		buf.DeleteCharBeforeCursor()
		return
	}

	if event.Key() == tcell.KeyEnter {
		buf.InsertNewLine()
		return
	}

	if event.Key() == tcell.KeyRune {
		r := []rune(event.Str())[0]
		buf.InsertCharAtCurrPos(r)
		return
	}

}

type VisualMode struct{}

var _ editorApi.EditorMode = (*VisualMode)(nil)

func (m *VisualMode) KeyHandler(event *tcell.EventKey, buf editorApi.EditorBuffer) {
	handleMovement(event, buf)
	handleSharedKeys(event, buf)

}

func handleSharedKeys(event *tcell.EventKey, buf editorApi.EditorBuffer) {
	switch event.Key() {
	case tcell.KeyLeft:
		buf.MoveCursor(1, editorApi.DirLeft)
	case tcell.KeyRight:
		buf.MoveCursor(1, editorApi.DirRight)
	case tcell.KeyUp:
		buf.MoveCursor(1, editorApi.DirUp)
	case tcell.KeyDown:
		buf.MoveCursor(1, editorApi.DirDown)
	}
}

func handleMovement(event *tcell.EventKey, buf editorApi.EditorBuffer) {

	switch event.Key() {
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
		case "o":
			buf.InsertNewLine()
		case "O":
			buf.UpsertNewLine()
		case "$":
			buf.JumpToLineEnd()
		case "0":
			buf.JumpToLineStart()
		}

	}
}
