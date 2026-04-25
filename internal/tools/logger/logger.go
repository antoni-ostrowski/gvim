package logger

import (
	"io"

	"github.com/antoni-ostrowski/gvim/internal/buffer"
	editorApi "github.com/antoni-ostrowski/gvim/internal/editor_api"
	"github.com/antoni-ostrowski/gvim/internal/utils"
	"github.com/antoni-ostrowski/gvim/internal/vim"
	"github.com/gdamore/tcell/v3"
)

type Logger struct {
	Input      editorApi.TextBuffer
	VimMachine editorApi.VimStateMachine
	active     bool
}

var _ editorApi.EditorTool = (*Logger)(nil)

func New(screen tcell.Screen) *Logger {
	w, h := screen.Size()

	c := buffer.NewGapBuffer("", &editorApi.Position{BaseX: 0, BaseY: h - 10, Width: w, Height: 20})

	return &Logger{
		Input:      c,
		VimMachine: vim.NewMachineInsertMode(),
		active:     false,
	}
}

func (l *Logger) HandleKey(event *tcell.EventKey, api editorApi.EditorApi) bool {
	if !l.active {
		return false
	}

	if event.Key() == tcell.KeyESC {
		switch l.VimMachine.GetMode().(type) {
		case *vim.Normal:
			l.active = false
			return true
		}
	}

	return l.VimMachine.Handler(event, l.Input)
}

func (l *Logger) Draw(screen tcell.Screen) {
	if !l.active {
		return
	}

	w, h := screen.Size()
	count := l.Input.LineCount()
	l.Input.GetPosition().BaseY = (h - count) - 2
	l.Input.GetPosition().BaseX = 0
	l.Input.GetPosition().Width = w

	l.Input.Draw(screen)
}

func (l *Logger) Log(mes string) {
	utils.Debuglog("logged this mes %v", mes)
	utils.Debuglog("setting to these bytes %v", []byte(mes))
	l.Input.SetBytes([]byte(mes))

	utils.Debuglog("input %v", string(l.Input.Bytes()))
	l.active = true
}

func (l *Logger) Append(mes string) {
	curBytes := l.Input.Bytes()
	l.Input.SetBytes(append(curBytes, []byte(mes)...))
	l.active = true
}

func (l *Logger) SetStyle(s tcell.Style) {
}

func (a *Logger) LogWriter() io.Writer {
	return a
}

func (a *Logger) Write(p []byte) (n int, err error) {
	s := string(p)
	a.Log(s)
	return 0, nil
}
