package app

import "github.com/gdamore/tcell/v3"

func DrawAppState(screen tcell.Screen, appState *App) {
	screen.Clear()

	appState.EditorBuffer.Draw(screen)
	DrawStatusLine(screen, appState)

	for _, elem := range appState.UiElements {
		elem.Draw(screen)
	}

	screen.Show()
}
