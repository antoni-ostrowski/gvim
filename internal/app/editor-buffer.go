package app

import (
	editorApi "github.com/antoni-ostrowski/gvim/internal/editor-api"
	"github.com/gdamore/tcell/v3"
)

type EditorTextBuffer struct {
	Lines   [][]rune
	CursorX int
	CursorY int
}

var _ editorApi.EditorBuffer = (*EditorTextBuffer)(nil)

func (e *EditorTextBuffer) Draw(screen tcell.Screen) {
	screen.PutStrStyled(e.CursorX, e.CursorY, "ntsear", tcell.StyleDefault)
	screen.ShowCursor(e.CursorX, e.CursorY)
}

func (e *EditorTextBuffer) MoveCursor(amount int, direction editorApi.Direction) {
	currX := e.CursorX
	currY := e.CursorY
	switch direction {
	case editorApi.DirLeft:
		if currX != 0 {
			e.CursorX = currX - amount
		}
	case editorApi.DirRight:
		e.CursorX = currX + amount
	case editorApi.DirUp:
		if currY != 0 {
			e.CursorY = currY - 1
		}
	case editorApi.DirDown:
		e.CursorY = currY + 1
	}
}
