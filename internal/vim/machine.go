package vim

import (
	editorApi "github.com/antoni-ostrowski/gvim/internal/editor_api"
	"github.com/gdamore/tcell/v3"
)

type Machine struct {
	Mode editorApi.VimMode
}

func NewMachine() *Machine {
	return &Machine{
		Mode: &Normal{},
	}
}

var _ editorApi.VimStateMachine = (*Machine)(nil)

func (m *Machine) Handler(event *tcell.EventKey, buf editorApi.TextBuffer) {
	modeSwitchHandler := func() editorApi.VimMode {
		switch m.Mode.(type) {
		case *Normal:
			if event.Key() == tcell.KeyRune {
				str := event.Str()
				switch str {
				case "i":
					return &Insert{}
				case "a":
					buf.MoveCursor(1, editorApi.DirRight)
					return &Insert{}
				case "A":
					buf.JumpToLineEnd()
					return &Insert{}
				case "v":
					return &Visual{}
				case "V":
					return &Visual{}
				}
			}
		case *Insert:
			if event.Key() == tcell.KeyEsc {
				return &Normal{}
			}
		case *Visual:
			if event.Key() == tcell.KeyEsc {
				return &Normal{}
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

type Normal struct{}

var _ editorApi.VimMode = (*Normal)(nil)

func (m *Normal) KeyHandler(event *tcell.EventKey, buf editorApi.TextBuffer) {
	handleMovement(event, buf)
	handleSharedKeys(event, buf)

	if event.Key() == tcell.KeyEnter {
		buf.MoveCursor(1, editorApi.DirDown)
		return
	}

}

type Insert struct{}

var _ editorApi.VimMode = (*Insert)(nil)

func (m *Insert) KeyHandler(event *tcell.EventKey, buf editorApi.TextBuffer) {
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

type Visual struct{}

var _ editorApi.VimMode = (*Visual)(nil)

func (m *Visual) KeyHandler(event *tcell.EventKey, buf editorApi.TextBuffer) {
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
