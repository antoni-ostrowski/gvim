package render

import (
	"github.com/antoni-ostrowski/gvim/internal/app"
	"github.com/gdamore/tcell/v3"
)

func DrawAppState(screen tcell.Screen, app *app.App) {
	screen.PutStrStyled(2, 2, "test from renderer", tcell.StyleDefault)
	screen.Show()
}
