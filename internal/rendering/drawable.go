package rendering

import (
	editorApi "github.com/antoni-ostrowski/gvim/internal/editor-api"
	"github.com/gdamore/tcell/v3"
)

type Drawable interface {
	Draw(screen tcell.Screen)
	HandleKey(event *tcell.EventKey, editorApi editorApi.EditorApi) bool
}
