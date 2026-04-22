package buffer

import (
	"slices"

	editorApi "github.com/antoni-ostrowski/gvim/internal/editor_api"
	"github.com/antoni-ostrowski/gvim/internal/utils"
	"github.com/gdamore/tcell/v3"
)

// we use go convention, so gap start is inclusive, gap end is exclusive
type GapTextBuffer struct {
	Data             []rune
	GapStart         int
	GapEnd           int
	CursorX, CursorY int
	ScrollOffset     int
	*editorApi.Position
	Style tcell.Style
}

var _ editorApi.TextBuffer = (*GapTextBuffer)(nil)

func NewGapBuffer(text string, pos *editorApi.Position) *GapTextBuffer {
	initGapSize := 1024
	runes := []rune(text)
	totalSize := initGapSize + len(runes)
	data := make([]rune, totalSize)
	copy(data, runes)

	return &GapTextBuffer{Data: data, GapStart: len(runes),
		GapEnd:       totalSize,
		Position:     pos,
		CursorY:      0,
		CursorX:      0,
		ScrollOffset: 0,
		Style:        tcell.StyleDefault,
	}
}

func (e *GapTextBuffer) GetPosition() *editorApi.Position {
	return e.Position
}

func (e *GapTextBuffer) Bytes() []byte {
	first := ([]byte(string(e.Data[:e.GapStart])))
	second := ([]byte(string(e.Data[e.GapEnd:])))
	return slices.Concat(first, second)
}

func (e *GapTextBuffer) SetBytes(content []byte) {
	runes := []rune(string(content))
	gap := make([]rune, 1024)
	e.Data = append(runes, gap...)
	e.GapStart = len(runes)
	e.GapEnd = len(e.Data)
}
func (e *GapTextBuffer) SetCursorX(newPos int) {
	e.CursorX = newPos
}

func (e *GapTextBuffer) Clean() {
	e.Data = []rune{}
	e.CursorY = 0
	e.CursorX = 0
	e.ScrollOffset = 0
	e.GapStart = 0
	e.GapEnd = 0
}

func (e *GapTextBuffer) SetStyle(s tcell.Style) {
	utils.Debuglog("setting style to %v", s)
	e.Style = s
}

func (e *GapTextBuffer) Draw(screen tcell.Screen) {
	utils.Debuglog("drawing with style %v", e.Style)
	drawX := e.Position.BaseX

	lineNum := 0
	colNum := 0
	for i, rune := range e.Data {
		// Track cursor position when we hit the gap
		if i == e.GapStart {
			e.CursorX = colNum
			e.CursorY = lineNum
		}

		// Skip gap buffer contents
		if i >= e.GapStart && i < e.GapEnd {
			continue
		}

		// Check if this line should be drawn
		isVisibleLine := lineNum >= e.ScrollOffset && lineNum < e.ScrollOffset+e.Position.Height

		if rune == '\n' {
			if isVisibleLine {
				// Only draw newline if it fits within width
				if colNum < e.Position.Width {
					screenY := lineNum - e.ScrollOffset + e.Position.BaseY
					screen.PutStrStyled(drawX+colNum, screenY, " ", e.Style)
				}
			}
			lineNum++
			colNum = 0
			continue
		}

		if !isVisibleLine {
			colNum++
			continue
		}

		// Skip if beyond width
		if colNum >= e.Position.Width {
			colNum++
			continue
		}

		screenY := lineNum - e.ScrollOffset + e.Position.BaseY
		screenX := drawX + colNum
		screen.PutStrStyled(screenX, screenY, string(rune), e.Style)
		colNum++
	}

	// Handle cursor at end of buffer (after loop)
	if e.GapStart >= len(e.Data)-e.GapEnd+e.GapStart {
		e.CursorX = colNum
		e.CursorY = lineNum
	}

	// Show cursor on screen (adjusted for scroll)
	cursorScreenY := e.CursorY - e.ScrollOffset + e.Position.BaseY
	cursorScreenX := drawX + e.CursorX
	if e.CursorY >= e.ScrollOffset && e.CursorY < e.ScrollOffset+e.Position.Height {
		screen.ShowCursor(cursorScreenX, cursorScreenY)
	}
}

func (e *GapTextBuffer) MoveCursor(amount int, direction editorApi.Direction) {
	switch direction {
	case editorApi.DirLeft:
		target := max(0, e.GapStart-amount)
		e.MoveGapTo(target)
		e.adjustScrollForCursor()
	case editorApi.DirRight:
		target := min(e.logicalLen(), e.GapStart+amount)
		e.MoveGapTo(target)
		e.adjustScrollForCursor()
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
		e.adjustScrollForCursor()

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
		e.adjustScrollForCursor()

	}

}

func (e *GapTextBuffer) InsertCharAtCurrPos(char rune) {
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

func (e *GapTextBuffer) DeleteCharBeforeCursor() {
	if e.GapStart > 0 {
		e.GapStart--
	}
}

func (e *GapTextBuffer) InsertNewLine() {
	e.MoveCursor(1, editorApi.DirDown)
	e.InsertCharAtCurrPos('\n')
	e.MoveCursor(1, editorApi.DirUp)
}

func (e *GapTextBuffer) UpsertNewLine() {
	e.InsertCharAtCurrPos('\n')
	e.MoveCursor(1, editorApi.DirUp)
}

func (e *GapTextBuffer) JumpToLineStart() {
	c := e.findLineStart(e.GapStart)
	e.MoveGapTo(c)
}

func (e *GapTextBuffer) JumpToLineEnd() {
	c := e.findLineEnd(e.GapStart)
	e.MoveGapTo(c)
}

func (e *GapTextBuffer) findLineStart(pos int) int {
	for i := pos - 1; i >= 0; i-- {
		if e.charAt(i) == '\n' {
			return i + 1
		}
	}
	return 0

}
func (e *GapTextBuffer) findLineEnd(pos int) int {
	for i := pos; i < e.logicalLen(); i++ {
		if e.charAt(i) == '\n' {
			return i
		}
	}
	return e.logicalLen()
}

func (e *GapTextBuffer) charAt(logicalIndex int) rune {
	if logicalIndex < e.GapStart {
		return e.Data[logicalIndex]
	}
	// Add the gap width to skip over the empty space
	gapWidth := e.GapEnd - e.GapStart
	return e.Data[logicalIndex+gapWidth]
}

func (e *GapTextBuffer) logicalLen() int {
	gapWidth := e.GapEnd - e.GapStart
	return len(e.Data) - gapWidth
}
func (e *GapTextBuffer) MoveGapTo(target int) {
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

func (e *GapTextBuffer) MoveGapLeftByOne() {
	// copy the rune over to right side of the gap
	if e.GapStart > 0 {
		e.GapStart--
		e.GapEnd--
		e.Data[e.GapEnd] = e.Data[e.GapStart]
	}
}

func (e *GapTextBuffer) MoveGapRightByOne() {
	if e.GapEnd < len(e.Data) {
		e.Data[e.GapStart] = e.Data[e.GapEnd]
		e.GapStart++
		e.GapEnd++
	}
}
func (e *GapTextBuffer) getCursorLine() int {
	line := 0
	for i := 0; i < e.GapStart && i < len(e.Data); i++ {
		if e.Data[i] == '\n' {
			line++
		}
	}
	return line
}

func (e *GapTextBuffer) adjustScrollForCursor() {
	cursorLine := e.getCursorLine()
	// Scroll up if cursor is above visible area
	if cursorLine < e.ScrollOffset {
		e.ScrollOffset = cursorLine
	}
	// Scroll down if cursor is at or below visible area
	if cursorLine >= e.ScrollOffset+e.Position.Height {
		e.ScrollOffset = cursorLine - e.Position.Height + 1
	}
}
