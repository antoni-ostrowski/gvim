package tools

import (
	"fmt"

	utils "github.com/antoni-ostrowski/gvim/internal"
	editorApi "github.com/antoni-ostrowski/gvim/internal/editor-api"
	"github.com/antoni-ostrowski/gvim/internal/rendering"
	"github.com/gdamore/tcell/v3"
)

type CommandPrompt struct {
	Input  rendering.TextInput
	active bool
}

func NewCommandPrompt(screen tcell.Screen) *CommandPrompt {
	_, y := screen.Size()
	return &CommandPrompt{
		Input: rendering.TextInput{X: 1, Y: y - 1, Buffer: []rune{}},
	}
}

var _ editorApi.EditorTool = (*CommandPrompt)(nil)

func (c *CommandPrompt) Draw(screen tcell.Screen) {
	if !c.active {
		return
	}

	c.Input.Draw(screen)
}

func (c *CommandPrompt) HandleKey(event *tcell.EventKey, api editorApi.EditorApi) bool {
	isActivationCombo := event.Key() == tcell.KeyRune && event.Str() == ":"

	if isActivationCombo {
		utils.Debuglog("cmd not active if hit!")
		c.active = true
		if len(c.Input.Buffer) > 0 {
			c.Input.Buffer = []rune{}
			c.Input.CursorPos = 0
		}
		return true
	}

	if c.active == false {
		return false
	}

	switch event.Key() {
	case tcell.KeyEsc:
		c.active = false
	case tcell.KeyEnter:
		if len(c.Input.Buffer) == 0 {
			return true
		}

		if string(c.Input.Buffer[0]) == "q" {
			api.SendQuitSignal()
			return true
		}

		if string(c.Input.Buffer[0]) == "w" {
			err := api.WriteFile()
			if err != nil {
				api.Log(fmt.Errorf("file: err writing to a file: %w", err).Error())
			}
			api.Log(fmt.Sprintf("wrote to file: %s", api.CurrentBufferPath()))
			api.TriggerEvent(tcell.NewEventKey(tcell.KeyEsc, "", tcell.ModNone))

			return true
		}
		return true
	}
	return c.Input.HandleKey(event, api)
}
