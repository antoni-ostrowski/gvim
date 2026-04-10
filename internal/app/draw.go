package app

import (
	"github.com/antoni-ostrowski/gvim/internal/machine"
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
	case *machine.NormalMode:
		screen.SetCursorStyle(tcell.CursorStyleDefault)
	case *machine.InsertMode:
		screen.SetCursorStyle(tcell.CursorStyleBlinkingBlock)
	case *machine.VisualMode:
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
	case *machine.NormalMode:
		return "NORMAL"
	case *machine.InsertMode:
		return "INSERT"
	case *machine.VisualMode:
		return "VISUAL"
	}
	return "UNKNOWN"
}
