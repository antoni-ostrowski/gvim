package app

import (
	"os"

	"github.com/gdamore/tcell/v3"
	"github.com/mattn/go-runewidth"
)

type App struct {
	ScreenEventChan chan tcell.Event
	Mode            EditorMode
	CommandPrompt   CommandPrompt
}

type VimMachine struct {
}

type CommandPrompt struct {
	input TextInput
}

type EditorMode interface {
	GetMode() string
	KeyHandler(screen tcell.Screen, app *App, event *tcell.EventKey)
}

type NormalMode struct {
}

func (m *NormalMode) KeyHandler(screen tcell.Screen, app *App, event *tcell.EventKey) {
	switch event.Key() {
	case tcell.KeyCtrlC:
		screen.Fini()
		os.Exit(0)
	case tcell.KeyCtrlV:
		app.Mode = &VisualBlockMode{}
	case tcell.KeyRune:
		str := event.Str()
		DrawDebugStr(screen, str)
		if str == ":" {
			app.Mode = &CommandPromptMode{}
			app.CommandPrompt.input.Buffer = []rune{}
			app.CommandPrompt.input.CursorPos = 0
			return
		}
		if str == "i" {
			app.Mode = &InsertMode{}
			return
		}
		if str == "a" {
			app.Mode = &InsertMode{}
			return
		}

		if str == "v" {
			app.Mode = &VisualMode{}
			return
		}
	}
}

func (m *NormalMode) GetMode() string {
	return "normal"
}

type InsertMode struct {
}

func (m *InsertMode) KeyHandler(screen tcell.Screen, app *App, event *tcell.EventKey) {
	switch event.Key() {
	case tcell.KeyCtrlC:
		screen.Fini()
		os.Exit(0)
	case tcell.KeyEscape:
		app.Mode = &NormalMode{}
	}
}

func (m *InsertMode) GetMode() string {
	return "insert"
}

type VisualMode struct {
}

func (m *VisualMode) KeyHandler(screen tcell.Screen, app *App, event *tcell.EventKey) {
	switch event.Key() {
	case tcell.KeyCtrlC:
		screen.Fini()
		os.Exit(0)
	case tcell.KeyEscape:
		app.Mode = &NormalMode{}
	}
}

func (m *VisualMode) GetMode() string {
	return "visual"
}

type VisualBlockMode struct {
}

func (m *VisualBlockMode) KeyHandler(screen tcell.Screen, app *App, event *tcell.EventKey) {
	switch event.Key() {
	case tcell.KeyCtrlC:
		screen.Fini()
		os.Exit(0)
	case tcell.KeyEscape:
		app.Mode = &NormalMode{}
	}
}

func (m *VisualBlockMode) GetMode() string {
	return "visualBlock"
}

type CommandPromptMode struct {
}

func (m *CommandPromptMode) KeyHandler(screen tcell.Screen, app *App, event *tcell.EventKey) {
	switch event.Key() {
	case tcell.KeyCtrlC:
		screen.Fini()
		os.Exit(0)
	case tcell.KeyRune:
		app.CommandPrompt.input.Y = 10
		app.CommandPrompt.input.X = 10
		app.CommandPrompt.input.HandleKey(event)
	case tcell.KeyEscape:
		app.Mode = &NormalMode{}
	}
}

func (m *CommandPromptMode) GetMode() string {
	return "commandPromptMode"
}

func DrawAppState(screen tcell.Screen, app *App) {
	screen.Clear()
	screen.PutStrStyled(0, 0, "test from renderer func - "+app.Mode.GetMode(), tcell.StyleDefault)

	if app.Mode.GetMode() == "commandPromptMode" {
		app.CommandPrompt.input.Draw(screen)
	}
	screen.Show()
}

func DrawDebugStr(screen tcell.Screen, msg string) {
	w, h := screen.Size()
	screen.PutStrStyled(w-len(msg)-1, h-1, msg, tcell.StyleDefault)
}

type TextInput struct {
	X, Y      int    // Screen position of input field
	CursorPos int    // Position in runes (0 = before first char)
	Buffer    []rune // The text content
}

func (t *TextInput) HandleKey(ev *tcell.EventKey) {
	switch ev.Key() {
	case tcell.KeyRune:
		// Insert character at cursor position
		t.Buffer = append(t.Buffer[:t.CursorPos],
			append([]rune(ev.Str()), t.Buffer[t.CursorPos:]...)...)
		t.CursorPos++

	case tcell.KeyBackspace, tcell.KeyBackspace2:
		if t.CursorPos > 0 {
			t.Buffer = append(t.Buffer[:t.CursorPos-1], t.Buffer[t.CursorPos:]...)
			t.CursorPos--
		}

	case tcell.KeyDelete:
		if t.CursorPos < len(t.Buffer) {
			t.Buffer = append(t.Buffer[:t.CursorPos], t.Buffer[t.CursorPos+1:]...)
		}

	case tcell.KeyLeft:
		if t.CursorPos > 0 {
			t.CursorPos--
		}

	case tcell.KeyRight:
		if t.CursorPos < len(t.Buffer) {
			t.CursorPos++
		}

	case tcell.KeyHome:
		t.CursorPos = 0

	case tcell.KeyEnd:
		t.CursorPos = len(t.Buffer)
	}
}
func (t *TextInput) Draw(s tcell.Screen) {
	// Draw the text
	text := string(t.Buffer)
	s.PutStr(t.X, t.Y, text)

	// Calculate cursor screen position (accounting for Unicode width!)
	cursorScreenX := t.X + t.runeWidth(t.Buffer[:t.CursorPos])
	s.ShowCursor(cursorScreenX, t.Y)
}

// Calculate display width of runes (CJK/emoji = 2 cells)
func (t *TextInput) runeWidth(runes []rune) int {
	w := 0
	for _, r := range runes {
		w += runewidth.RuneWidth(r) // github.com/mattn/go-runewidth
	}
	return w
}
