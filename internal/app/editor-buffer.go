package app

import (
	utils "github.com/antoni-ostrowski/gvim/internal"
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
	screen.ShowCursor(e.CursorX, e.CursorY)
	for lineIndex, line := range e.Lines {
		for charIndex, char := range line {
			screen.PutStrStyled(charIndex, lineIndex, string(char), tcell.StyleDefault)
		}

	}
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

func (e *EditorTextBuffer) InsertCharAtCurrPos(char rune) {
	utils.Debuglog("cursor x = %v, cursor y = %v ", e.CursorX, e.CursorY)
	utils.Debuglog("lines %v", e.Lines)
	utils.Debuglog("send char %v", char)

	for len(e.Lines) <= e.CursorY {
		e.Lines = append(e.Lines, []rune{})
	}

	for len(e.Lines[e.CursorY]) <= e.CursorX {
		e.Lines[e.CursorY] = append(e.Lines[e.CursorY], ' ')
	}

	e.Lines[e.CursorY][e.CursorX] = char
	e.CursorX++

}
