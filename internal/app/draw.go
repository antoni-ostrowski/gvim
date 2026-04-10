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
	case *vim.Normal:
		screen.SetCursorStyle(tcell.CursorStyleDefault)
	case *vim.Insert:
		screen.SetCursorStyle(tcell.CursorStyleBlinkingBlock)
	case *vim.Visual:
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
	case *vim.Normal:
		return "NORMAL"
	case *vim.Insert:
		return "INSERT"
	case *vim.Visual:
		return "VISUAL"
	}
	return "UNKNOWN"
}
