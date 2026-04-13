package cmdprompt

import (
	"fmt"
	"strings"

	"github.com/antoni-ostrowski/gvim/internal/buffer"
	editorApi "github.com/antoni-ostrowski/gvim/internal/editor_api"
	utils "github.com/antoni-ostrowski/gvim/internal/utils"
	"github.com/gdamore/tcell/v3"
	"github.com/spf13/cobra"
)

type CommandPrompt struct {
	Input  buffer.LineBuffer
	active bool
}

var _ editorApi.EditorTool = (*CommandPrompt)(nil)

func New(screen tcell.Screen, api editorApi.EditorApi) *CommandPrompt {
	_, y := screen.Size()

	createCmds(api)

	return &CommandPrompt{
		Input: buffer.LineBuffer{X: 1, Y: y - 1, Buffer: []rune{}},
	}
}

func createCmds(api editorApi.EditorApi) *cobra.Command {
	rootCmd := api.RootCmd()

	quitCmd := &cobra.Command{
		Use:     "quit",
		Aliases: []string{"q"},
		Short:   "Quit the editor",
		Run: func(cmd *cobra.Command, args []string) {
			api.SendQuitSignal()
		},
	}

	writeCmd := &cobra.Command{
		Use:     "write",
		Aliases: []string{"w"},
		Short:   "Write buffer to file",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := api.WriteFile()
			if err != nil {
				return fmt.Errorf("error writing to file: %w", err)
			}
			api.Log(fmt.Sprintf("wrote to file: %s", api.CurrentBufferPath()))
			return nil
		},
	}

	openCmd := &cobra.Command{
		Use:     "open",
		Aliases: []string{"o"},
		Short:   "Open file",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				filename := args[0]
				err := api.OpenFile(filename)
				if err != nil {
					return err
				}
				api.Log(fmt.Sprintf("opened %s", filename))
			}

			return nil
		},
	}

	rootCmd.AddCommand(quitCmd, writeCmd, openCmd)

	return rootCmd
}

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

	rootCmd := api.RootCmd()

	switch event.Key() {
	case tcell.KeyEsc:
		c.active = false
	case tcell.KeyEnter:
		if len(c.Input.Buffer) == 0 {
			return true
		}

		input := string(c.Input.Buffer)
		args := strings.Fields(input)
		rootCmd.SetArgs(args)
		err := rootCmd.Execute()

		if err != nil {
			api.Log(err.Error())
		}
		c.active = false

		return true
	}
	return c.Input.HandleKey(event, api)
}
