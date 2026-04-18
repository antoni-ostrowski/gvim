package app

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/antoni-ostrowski/gvim/internal/buffer"
	editorApi "github.com/antoni-ostrowski/gvim/internal/editor_api"
	"github.com/antoni-ostrowski/gvim/internal/vim"
	"github.com/gdamore/tcell/v3"
	"github.com/spf13/cobra"
)

type App struct {
	Machine       editorApi.VimStateMachine
	Tools         map[string]editorApi.EditorTool
	QuitChn       chan struct{}
	EventChan     chan tcell.Event
	Screen        tcell.Screen
	EditorBuffer  editorApi.TextBuffer
	Cwd           string
	CurOpenedFile string
	LogMess       string
	rootCmd       *cobra.Command
}

var _ editorApi.EditorApi = (*App)(nil)

func NewApp(screen tcell.Screen, argPath string, eventChan chan tcell.Event) *App {
	contents, cwd, absFilePath, err := parseFilepathArg(argPath)
	if err != nil {
		log.Fatalf("failed to parse arg filepath: %v", err)
	}

	w, h := screen.Size()

	app := &App{
		Machine:      vim.NewMachine(),
		QuitChn:      make(chan struct{}, 1),
		Screen:       screen,
		Tools:        make(map[string]editorApi.EditorTool),
		EditorBuffer: buffer.NewGapBuffer(string(contents), &editorApi.Position{BaseX: 0, BaseY: 0, Width: w, Height: h - 2}),
		EventChan:    eventChan,
		LogMess:      "",
		rootCmd: &cobra.Command{
			Use:           "gvim",
			SilenceErrors: true,
			SilenceUsage:  true,
		},
		Cwd:           cwd,
		CurOpenedFile: absFilePath,
	}

	return app
}

func parseFilepathArg(argPath string) ([]byte, string, string, error) {
	isFile := func(path string) error {
		info, err := os.Stat(path)
		if err == nil && info.IsDir() {
			// Handle opening a directory (maybe show a file explorer?)
			return errors.New("directory found instead of a file")
		} else {
			return nil
		}
	}

	absPath, err := filepath.Abs(argPath)
	if err != nil {
		return []byte{}, "", "", fmt.Errorf("parse arg filepath: %w", err)
	}

	err = isFile(absPath)
	if err != nil {
		return []byte{}, absPath, "", nil
	}

	contents, err := os.ReadFile(absPath)
	if err != nil {
		return []byte{}, "", "", fmt.Errorf("read file: %w", err)
	}

	return contents, "", absPath, nil
}

func (a *App) RootCmd() *cobra.Command {
	return a.rootCmd
}

func (a *App) TriggerEvent(event tcell.Event) {
	a.EventChan <- event
}

func (a *App) Log(mess string) {
	a.LogMess = mess
}

func (a *App) CurrentOpenedFilePath() string {
	return a.CurOpenedFile
}
func (a *App) OpenFile(file string) error {
	absPath, err := filepath.Abs(file)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}

	contents, err := os.ReadFile(absPath)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}

	a.EditorBuffer = buffer.NewGapBuffer(string(contents), a.Buffer().GetPosition())
	return nil
}

func (a *App) WriteFile() error {
	f, err := os.Create(a.CurOpenedFile)
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

func (a *App) Buffer() editorApi.TextBuffer {
	return a.EditorBuffer
}
