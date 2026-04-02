package app

import (
	editorApi "github.com/antoni-ostrowski/gvim/internal/editor-api"
	"github.com/antoni-ostrowski/gvim/internal/rendering"
	"github.com/gdamore/tcell/v3"
)

type EditorBuffer struct {
	Lines   [][]rune
	CursorX int
	CursorY int
}

var _ rendering.Drawable = (*EditorBuffer)(nil)

func (e *EditorBuffer) Draw(screen tcell.Screen) {
	screen.PutStrStyled(e.CursorX, e.CursorY, "ntsear", tcell.StyleDefault)

	screen.ShowCursor(e.CursorX, e.CursorY)
}

func (e *EditorBuffer) HandleKey(event *tcell.EventKey, editorApi editorApi.EditorApi) bool {
	return true
}
