package editorApi

import "github.com/gdamore/tcell/v3"

type EditorApi interface {
	SendQuitSignal()
	ToggleCommandPrompt(active bool)
	MoveEditorBuffCursor(amount int, direction Direction)
	CurrentMode() EditorMode
}

type EditorMode interface {
	GetMode() string
	KeyHandler(event *tcell.EventKey, editorApi EditorApi) EditorMode
}

type Direction interface {
	isDirection()
}

type Left struct{}
type Right struct{}

func (Left) isDirection()  {}
func (Right) isDirection() {}

type Up struct{}
type Down struct{}

func (Up) isDirection()   {}
func (Down) isDirection() {}
