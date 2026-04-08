package app

import (
	"errors"
	"os"
	"path/filepath"

	editorApi "github.com/antoni-ostrowski/gvim/internal/editor-api"
	"github.com/antoni-ostrowski/gvim/internal/machine"
	"github.com/antoni-ostrowski/gvim/internal/rendering"
	"github.com/gdamore/tcell/v3"
)

type App struct {
	Machine      editorApi.VimMachine
	UiElements   []editorApi.UiElement
	QuitChn      chan struct{}
	Screen       tcell.Screen
	EditorBuffer editorApi.EditorBuffer
	ArgPath      string
}

var _ editorApi.EditorApi = (*App)(nil)

func NewApp(screen tcell.Screen, argPath string) *App {
	app := &App{
		Machine:      &machine.VimMachine{Mode: &machine.NormalMode{}},
		QuitChn:      make(chan struct{}, 1),
		Screen:       screen,
		UiElements:   []editorApi.UiElement{},
		EditorBuffer: NewEditorBuffer(""),
		ArgPath:      argPath,
	}
	absPath, err := filepath.Abs(argPath)
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
func (a *App) CurrentMode() editorApi.EditorMode {
	return a.Machine.GetMode()
}

func (a *App) Buffer() editorApi.EditorBuffer {
	return a.EditorBuffer
}

func (a *App) ToggleCommandPrompt(active bool) {
	if active {
		_, h := a.Screen.Size()
		a.UiElements = append(a.UiElements, &CommandPrompt{Input: rendering.TextInput{X: 1, Y: h - 1, CursorPos: 0, Buffer: []rune{}}})
	} else {
		if len(a.UiElements) > 0 {
			index := len(a.UiElements) - 1
			a.UiElements[index] = nil
			a.UiElements = a.UiElements[:index]
		}
	}
}
