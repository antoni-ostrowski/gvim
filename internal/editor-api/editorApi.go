package editorApi

import "github.com/gdamore/tcell/v3"

type EditorApi interface {
	SendQuitSignal()
	ToggleCommandPrompt(active bool)
	CurrentMode() EditorMode
	Buffer() EditorBuffer
}

type EditorMode interface {
	KeyHandler(event *tcell.EventKey, editorApi EditorApi) EditorMode
}

type EditorBuffer interface {
	Drawable
	MoveCursor(amount int, direction Direction)
	InsertCharAtCurrPos(char rune)
	DeleteCharBeforeCursor()
}

type VimMachine interface {
	Handler(event *tcell.EventKey, api EditorApi)
	GetMode() EditorMode
}

type UiElement interface {
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
