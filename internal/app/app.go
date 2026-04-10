package app

import (
	"errors"
	"os"
	"path/filepath"

	editorApi "github.com/antoni-ostrowski/gvim/internal/editor-api"
	"github.com/antoni-ostrowski/gvim/internal/machine"
	"github.com/gdamore/tcell/v3"
)

type App struct {
	Machine      editorApi.VimMachine
	Tools        map[string]editorApi.EditorTool
	QuitChn      chan struct{}
	EventChan    chan tcell.Event
	Screen       tcell.Screen
	EditorBuffer editorApi.EditorBuffer
	ArgPath      string
	LogMess      string
}

var _ editorApi.EditorApi = (*App)(nil)

func NewApp(screen tcell.Screen, argPath string, eventChan chan tcell.Event) *App {
	app := &App{
		Machine:      &machine.VimMachine{Mode: &machine.NormalMode{}},
		QuitChn:      make(chan struct{}, 1),
		Screen:       screen,
		Tools:        make(map[string]editorApi.EditorTool),
		EditorBuffer: NewEditorBuffer(""),
		ArgPath:      argPath,
		EventChan:    eventChan,
		LogMess:      "",
	}
	absPath, err := filepath.Abs(argPath)
	app.LogMess = absPath
	if err != nil {
		app.EditorBuffer = NewEditorBuffer("")
	} else {
		err = isFile(absPath)
		if err == nil {
			contents, err := os.ReadFile(absPath)
			if err == nil {
				app.EditorBuffer = NewEditorBuffer(string(contents))
			}
		}
	}

	return app
}

func isFile(path string) error {
	info, err := os.Stat(path)
	if err == nil && info.IsDir() {
		// Handle opening a directory (maybe show a file explorer?)
		return errors.New("directory found instead of a file")
	} else {
		return nil
	}
}
func (a *App) TriggerEvent(event tcell.Event) {
	a.EventChan <- event

}

func (a *App) Log(mess string) {
	a.LogMess = mess
}

func (a *App) CurrentBufferPath() string {
	return a.ArgPath
}

func (a *App) WriteFile() error {
	f, err := os.Create(a.ArgPath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(a.Buffer().Bytes())
	if err != nil {
		return err
	}

	return nil
}

func (a *App) SendQuitSignal() {
	a.QuitChn <- struct{}{}
}

func (a *App) Buffer() editorApi.EditorBuffer {
	return a.EditorBuffer
}
