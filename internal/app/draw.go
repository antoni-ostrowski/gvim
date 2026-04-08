package app

import (
	"github.com/antoni-ostrowski/gvim/internal/machine"
	"github.com/gdamore/tcell/v3"
)

func DrawAppState(screen tcell.Screen, appState *App) {
	screen.Clear()

	appState.EditorBuffer.Draw(screen)
	DrawStatusLine(screen, appState)

	for _, elem := range appState.UiElements {
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
