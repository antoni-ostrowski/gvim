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
func NewMachineInsertMode() *Machine {
	return &Machine{
		Mode: &Insert{},
	}
}

var _ editorApi.VimStateMachine = (*Machine)(nil)

func (m *Machine) Handler(event *tcell.EventKey, buf editorApi.TextBuffer) bool {
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
		return true
	}

	return m.Mode.KeyHandler(event, buf)
}
func (m *Machine) GetMode() editorApi.VimMode {
	return m.Mode
}

type Normal struct{}

var _ editorApi.VimMode = (*Normal)(nil)

func (m *Normal) KeyHandler(event *tcell.EventKey, buf editorApi.TextBuffer) bool {
	if handleSharedKeys(event, buf) {
		return true
	}
	if handleMovement(event, buf) {
		return true
	}

	if event.Key() == tcell.KeyEnter {
		buf.MoveCursor(1, editorApi.DirDown)
		return true
	}

	return false

}

type Insert struct{}

var _ editorApi.VimMode = (*Insert)(nil)

func (m *Insert) KeyHandler(event *tcell.EventKey, buf editorApi.TextBuffer) bool {
	if handleSharedKeys(event, buf) {
		return true
	}

	if event.Key() == tcell.KeyBackspace {
		buf.DeleteCharBeforeCursor()
		return true
	}

	if event.Key() == tcell.KeyEnter {
		buf.InsertNewLine()
		return true
	}

	if event.Key() == tcell.KeyRune {
		r := []rune(event.Str())[0]
		buf.InsertCharAtCurrPos(r)
		return true
	}
	return false
}

type Visual struct{}

var _ editorApi.VimMode = (*Visual)(nil)

func (m *Visual) KeyHandler(event *tcell.EventKey, buf editorApi.TextBuffer) bool {
	if handleSharedKeys(event, buf) {
		return true
	}
	if handleMovement(event, buf) {
		return true
	}
	return false
}

func handleSharedKeys(event *tcell.EventKey, buf editorApi.TextBuffer) bool {
	switch event.Key() {
	case tcell.KeyLeft:
		buf.MoveCursor(1, editorApi.DirLeft)
		return true
	case tcell.KeyRight:
		buf.MoveCursor(1, editorApi.DirRight)
		return true
	case tcell.KeyUp:
		buf.MoveCursor(1, editorApi.DirUp)
		return true
	case tcell.KeyDown:
		buf.MoveCursor(1, editorApi.DirDown)
		return true
	}
	return false
}

func handleMovement(event *tcell.EventKey, buf editorApi.TextBuffer) bool {
	switch event.Key() {
	case tcell.KeyRune:
		str := event.Str()
		switch str {
		case "h":
			buf.MoveCursor(1, editorApi.DirLeft)
			return true
		case "l":
			buf.MoveCursor(1, editorApi.DirRight)
			return true
		case "k":
			buf.MoveCursor(1, editorApi.DirUp)
			return true
		case "j":
			buf.MoveCursor(1, editorApi.DirDown)
			return true
		case "o":
			buf.InsertNewLine()
			return true
		case "O":
			buf.UpsertNewLine()
			return true
		case "$":
			buf.JumpToLineEnd()
			return true
		case "0":
			buf.JumpToLineStart()
			return true
		}
	}
	return false
}
