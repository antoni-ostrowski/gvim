package editorApi

import (
	"github.com/gdamore/tcell/v3"
	"github.com/spf13/cobra"
)

type EditorApi interface {
	SendQuitSignal()
	Buffer() TextBuffer
	WriteFile() error
	CurrentOpenedFilePath() string
	Log(mess string)
	TriggerEvent(event tcell.Event)
	RootCmd() *cobra.Command
	OpenFile(filepath string) error
}

type TextBuffer interface {
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

type VimStateMachine interface {
	Handler(event *tcell.EventKey, buf TextBuffer)
	GetMode() VimMode
}

type VimMode interface {
	KeyHandler(event *tcell.EventKey, buf TextBuffer)
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
