package app

import (
	"github.com/antoni-ostrowski/gvim/internal/vim"
	"github.com/gdamore/tcell/v3"
)

func DrawAppState(screen tcell.Screen, appState *App) {
	screen.Clear()

	appState.EditorBuffer.Draw(screen)
	drawStatusLine(screen, appState)

	for _, elem := range appState.Tools {
		elem.Draw(screen)
	}

	switch appState.Machine.GetMode().(type) {
	case *vim.NormalMode:
		screen.SetCursorStyle(tcell.CursorStyleDefault)
	case *vim.InsertMode:
		screen.SetCursorStyle(tcell.CursorStyleBlinkingBlock)
	case *vim.VisualMode:
		screen.SetCursorStyle(tcell.CursorStyleDefault)
	}

	screen.Show()
}

func drawStatusLine(screen tcell.Screen, appState *App) {
	_, h := screen.Size()

	screen.PutStrStyled(0, h-3, appState.LogMess, tcell.StyleDefault)
	screen.PutStrStyled(0, h-2, getCurrentEditorModeName(appState), tcell.StyleDefault)
}

func getCurrentEditorModeName(appState *App) string {
	switch appState.Machine.GetMode().(type) {
	case *vim.NormalMode:
		return "NORMAL"
	case *vim.InsertMode:
		return "INSERT"
	case *vim.VisualMode:
		return "VISUAL"
	}
	return "UNKNOWN"
}
