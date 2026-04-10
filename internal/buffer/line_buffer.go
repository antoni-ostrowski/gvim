package buffer

import (
	editorApi "github.com/antoni-ostrowski/gvim/internal/editor_api"
	utils "github.com/antoni-ostrowski/gvim/internal/utils"
	"github.com/gdamore/tcell/v3"
	"github.com/mattn/go-runewidth"
)

type LineBuffer struct {
	X, Y      int    // Screen position of input field
	CursorPos int    // Position in runes (0 = before first char)
	Buffer    []rune // The text content
}

var _ editorApi.Drawable = (*LineBuffer)(nil)

func (t *LineBuffer) HandleKey(ev *tcell.EventKey, editorApi editorApi.EditorApi) bool {
	utils.Debuglog("TextInput.HandleKey: Key=%v, Str=%q, CursorPos=%d, BufferLen=%d, Buffer=%q", ev.Key(), ev.Str(), t.CursorPos, len(t.Buffer), string(t.Buffer))
	switch ev.Key() {
	case tcell.KeyRune:
		// Insert character at cursor position
		t.Buffer = append(t.Buffer[:t.CursorPos],
			append([]rune(ev.Str()), t.Buffer[t.CursorPos:]...)...)
		t.CursorPos++
		return true

	case tcell.KeyBackspace, tcell.KeyBackspace2:
		if t.CursorPos > 0 {
			t.Buffer = append(t.Buffer[:t.CursorPos-1], t.Buffer[t.CursorPos:]...)
			t.CursorPos--
		}
		return true

	case tcell.KeyDelete:
		if t.CursorPos < len(t.Buffer) {
			t.Buffer = append(t.Buffer[:t.CursorPos], t.Buffer[t.CursorPos+1:]...)
		}
		return true

	case tcell.KeyLeft:
		if t.CursorPos > 0 {
			t.CursorPos--
		}
		return true

	case tcell.KeyRight:
		if t.CursorPos < len(t.Buffer) {
			t.CursorPos++
		}
		return true

	case tcell.KeyHome:
		t.CursorPos = 0
		return true

	case tcell.KeyEnd:
		t.CursorPos = len(t.Buffer)
		return true
	default:
		return false
	}
}

func (t *LineBuffer) Draw(s tcell.Screen) {
	text := string(t.Buffer)
	s.PutStrStyled(t.X, t.Y, text, tcell.StyleDefault)

	cursorScreenX := t.X + t.runeWidth(t.Buffer[:t.CursorPos])

	s.ShowCursor(cursorScreenX, t.Y)
}

// Calculate display width of runes (CJK/emoji = 2 cells)
func (t *LineBuffer) runeWidth(runes []rune) int {
	w := 0
	for _, r := range runes {
		w += runewidth.RuneWidth(r) // github.com/mattn/go-runewidth
	}
	return w
}
