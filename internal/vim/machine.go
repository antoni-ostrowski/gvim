package vim

import (
	editorApi "github.com/antoni-ostrowski/gvim/internal/editor_api"
	"github.com/gdamore/tcell/v3"
)

type Machine struct {
	Mode        editorApi.VimMode
	pendingKeys []rune
}

var _ editorApi.VimStateMachine = (*Machine)(nil)

func (m *Machine) Handler(event *tcell.EventKey, buf editorApi.TextBuffer) {
	modeSwitchHandler := func() editorApi.VimMode {
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
func (m *Machine) GetMode() editorApi.VimMode {
	return m.Mode
}

type NormalMode struct{}

var _ editorApi.VimMode = (*NormalMode)(nil)

func (m *NormalMode) KeyHandler(event *tcell.EventKey, buf editorApi.TextBuffer) {
	handleMovement(event, buf)
	handleSharedKeys(event, buf)

	if event.Key() == tcell.KeyEnter {
		buf.MoveCursor(1, editorApi.DirDown)
		return
	}

}

type InsertMode struct{}

var _ editorApi.VimMode = (*InsertMode)(nil)

func (m *InsertMode) KeyHandler(event *tcell.EventKey, buf editorApi.TextBuffer) {
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

var _ editorApi.VimMode = (*VisualMode)(nil)

func (m *VisualMode) KeyHandler(event *tcell.EventKey, buf editorApi.TextBuffer) {
	handleMovement(event, buf)
	handleSharedKeys(event, buf)

}

func handleSharedKeys(event *tcell.EventKey, buf editorApi.TextBuffer) {
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

func handleMovement(event *tcell.EventKey, buf editorApi.TextBuffer) {
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
