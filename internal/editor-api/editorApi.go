package editorApi

import (
	"github.com/gdamore/tcell/v3"
)

type EditorApi interface {
	SendQuitSignal()
	Buffer() EditorBuffer
	WriteFile() error
}

type EditorMode interface {
	KeyHandler(event *tcell.EventKey, buf EditorBuffer)
}

type EditorBuffer interface {
	Drawable
	MoveCursor(amount int, direction Direction)
	InsertCharAtCurrPos(char rune)
	DeleteCharBeforeCursor()
	InsertNewLine()
	UpsertNewLine()
	JumpToLineStart()
	JumpToLineEnd()
	Bytes() []byte
}

type VimMachine interface {
	Handler(event *tcell.EventKey, buf EditorBuffer)
	GetMode() EditorMode
}

type EditorTool interface {
	Drawable
	KeyHandler
}

type KeyHandler interface {
	HandleKey(event *tcell.EventKey, api EditorApi) bool
}

type Drawable interface {
	Draw(screen tcell.Screen)
}

type Direction int

const (
	DirLeft Direction = iota
	DirRight
	DirUp
	DirDown
)
