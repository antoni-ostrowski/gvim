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
			newPos := len(e.Lines[e.CursorY])
			if newPos == 0 {
				e.CursorX = newPos
			}
		}
	case editorApi.DirDown:
		e.CursorY = currY + 1

		for len(e.Lines) <= e.CursorY {
			e.Lines = append(e.Lines, []rune{})
		}

		newPos := len(e.Lines[e.CursorY])
		if newPos == 0 {
			e.CursorX = newPos
		}
	}
}

func (e *EditorTextBuffer) InsertCharAtCurrPos(char rune) {
	for len(e.Lines) <= e.CursorY {
		e.Lines = append(e.Lines, []rune{})
	}

	for len(e.Lines[e.CursorY]) <= e.CursorX {
		e.Lines[e.CursorY] = append(e.Lines[e.CursorY], ' ')
	}

	e.Lines[e.CursorY][e.CursorX] = char
	e.CursorX++
}

func (e *EditorTextBuffer) DeleteCharBeforeCursor() {
	if e.CursorY >= len(e.Lines) {
		return
	}
	if e.CursorX == 0 {
		return
	}

	line := e.Lines[e.CursorY]
	if len(line) == 0 {
		return
	}

	if e.CursorX > len(line) {
		e.CursorX = len(line)
	}

	// Delete and write back
	e.Lines[e.CursorY] = append(line[:e.CursorX-1], line[e.CursorX:]...)
	e.CursorX--
}
