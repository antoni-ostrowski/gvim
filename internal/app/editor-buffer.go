package app

import (
	"slices"

	editorApi "github.com/antoni-ostrowski/gvim/internal/editor-api"
	"github.com/gdamore/tcell/v3"
)

// we use go convention, so gap start is inclusive, gap end is exclusive
// GapStart: The index of the first empty slot.
// GapEnd: The index of the first valid character after the gap.
type EditorGapBuffer struct {
	Data     []rune
	GapStart int
	GapEnd   int
}

var _ editorApi.EditorBuffer = (*EditorGapBuffer)(nil)

func NewEditorBuffer(text string) *EditorGapBuffer {
	initGapSize := 1024
	runes := []rune(text)
	totalSize := initGapSize + len(runes)
	data := make([]rune, totalSize)
	copy(data, runes)

	return &EditorGapBuffer{Data: data, GapStart: len(runes), GapEnd: totalSize}
}
func (e *EditorGapBuffer) Bytes() []byte {
	first := ([]byte(string(e.Data[:e.GapStart])))
	second := ([]byte(string(e.Data[e.GapEnd:])))
	return slices.Concat(first, second)
}
func (e *EditorGapBuffer) Draw(screen tcell.Screen) {
	drawX, drawY := 0, 0
	drawCursorX, drawCursorY := 0, 0
	for i, rune := range e.Data {
		// if we hit cursor position (so gap start), save the cords
		// thats our cursor position to draw
		if i == e.GapStart {
			drawCursorX = drawX
			drawCursorY = drawY
		}

		// skip drawing if gap buffer
		if i >= e.GapStart && i < e.GapEnd {
			continue
		}

		// handle new line
		if rune == '\n' {
			drawX = 0
			drawY++
			continue
		}

		screen.Put(drawX, drawY, string(rune), tcell.StyleDefault)
		drawX++
	}
	screen.ShowCursor(drawCursorX, drawCursorY)
}

func (e *EditorGapBuffer) MoveCursor(amount int, direction editorApi.Direction) {
	switch direction {
	case editorApi.DirLeft:
		target := max(0, e.GapStart-amount)
		e.MoveGapTo(target)
	case editorApi.DirRight:
		target := min(e.logicalLen(), e.GapStart+amount)
		e.MoveGapTo(target)
	case editorApi.DirUp:
		targetLineStart := e.findLineStart(e.GapStart)
		// how many runes deep is our cursor on the line
		column := e.GapStart - targetLineStart

		for range amount {
			// we are on the first line
			if targetLineStart == 0 {
				break
			}
			// find start of the line above us
			targetLineStart = e.findLineStart(targetLineStart - 1)
		}

		targetLineEnd := e.findLineEnd(targetLineStart)
		targetPos := min(targetLineStart+column, targetLineEnd)
		e.MoveGapTo(targetPos)

	case editorApi.DirDown:
		targetLineStart := e.findLineStart(e.GapStart) // Start with current line
		column := e.GapStart - targetLineStart         // Save the column

		for range amount {
			lineEnd := e.findLineEnd(targetLineStart)
			if lineEnd >= e.logicalLen() {
				break
			}
			targetLineStart = lineEnd + 1 // Step forward to next line
		}

		targetLineEnd := e.findLineEnd(targetLineStart)
		targetPos := min(targetLineStart+column, targetLineEnd)
		e.MoveGapTo(targetPos)

	}
}

func (e *EditorGapBuffer) InsertCharAtCurrPos(char rune) {
	if e.GapStart == e.GapEnd {
		// here we need to rsize the data slice and re create gap buffer
		oldData := e.Data
		oldDataLen := len(oldData)
		newDataSize := oldDataLen * 2
		if newDataSize == 0 {
			newDataSize = 1024
		}
		newData := make([]rune, newDataSize)
		// copy left side of gap buffer
		copy(newData, oldData[:e.GapStart])

		rightSideLen := oldDataLen - e.GapEnd
		newGapEnd := newDataSize - rightSideLen
		copy(newData[newGapEnd:], oldData[e.GapEnd:])
		e.Data = newData
		e.GapEnd = newGapEnd
	}
	e.Data[e.GapStart] = char
	e.GapStart++
}

func (e *EditorGapBuffer) DeleteCharBeforeCursor() {
	if e.GapStart > 0 {
		e.GapStart--
	}
}

func (e *EditorGapBuffer) InsertNewLine() {
	e.MoveCursor(1, editorApi.DirDown)
	e.InsertCharAtCurrPos('\n')
	e.MoveCursor(1, editorApi.DirUp)
}

func (e *EditorGapBuffer) UpsertNewLine() {
	e.InsertCharAtCurrPos('\n')
	e.MoveCursor(1, editorApi.DirUp)
}

func (e *EditorGapBuffer) JumpToLineStart() {
	c := e.findLineStart(e.GapStart)
	e.MoveGapTo(c)
}

func (e *EditorGapBuffer) JumpToLineEnd() {
	c := e.findLineEnd(e.GapStart)
	e.MoveGapTo(c)
}

func (e *EditorGapBuffer) findLineStart(pos int) int {
	for i := pos - 1; i >= 0; i-- {
		if e.charAt(i) == '\n' {
			return i + 1
		}
	}
	return 0

}
func (e *EditorGapBuffer) findLineEnd(pos int) int {
	for i := pos; i < e.logicalLen(); i++ {
		if e.charAt(i) == '\n' {
			return i
		}
	}
	return e.logicalLen()
}

func (e *EditorGapBuffer) charAt(logicalIndex int) rune {
	if logicalIndex < e.GapStart {
		return e.Data[logicalIndex]
	}
	// Add the gap width to skip over the empty space
	gapWidth := e.GapEnd - e.GapStart
	return e.Data[logicalIndex+gapWidth]
}

func (e *EditorGapBuffer) logicalLen() int {
	gapWidth := e.GapEnd - e.GapStart
	return len(e.Data) - gapWidth
}
func (e *EditorGapBuffer) MoveGapTo(target int) {
	// this moves gap according to the target

	// If the target is to the left of the current gap,
	// we move the gap left by shifting characters to the right.
	for e.GapStart > target {
		e.MoveGapLeftByOne()
	}

	// If the target is to the right of the current gap,
	// we move the gap right by shifting characters to the left.
	for e.GapStart < target {
		e.MoveGapRightByOne()
	}
}

func (e *EditorGapBuffer) MoveGapLeftByOne() {
	// copy the rune over to right side of the gap
	if e.GapStart > 0 {
		e.GapStart--
		e.GapEnd--
		e.Data[e.GapEnd] = e.Data[e.GapStart]
	}
}

func (e *EditorGapBuffer) MoveGapRightByOne() {
	if e.GapEnd < len(e.Data) {
		e.Data[e.GapStart] = e.Data[e.GapEnd]
		e.GapStart++
		e.GapEnd++
	}
}
