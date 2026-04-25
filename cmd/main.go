package main

import (
	"log"
	"os"

	"github.com/antoni-ostrowski/gvim/internal/app"
	"github.com/antoni-ostrowski/gvim/internal/tools/cmdprompt"
	"github.com/antoni-ostrowski/gvim/internal/tools/logger"
	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

func main() {
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := screen.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	quit := func() {
		screen.Fini()
	}
	defer quit()

	screen.EnablePaste()
	screen.Clear()

	eventChannel := screen.EventQ()

	path := ""
	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	defStyle := tcell.StyleDefault.Background(color.Reset).Foreground(color.Reset)
	screen.SetStyle(defStyle)

	appState := app.NewApp(screen, path, eventChannel)
	loger := logger.New(screen)
	appState.Tools["logger"] = loger
	appState.LoggerTool = loger
	appState.Tools["cmdPrompt"] = cmdprompt.New(screen, appState)

	for {
		select {
		case event := <-eventChannel:
			if event, ok := event.(*tcell.EventKey); ok {

				keyHandled := false

				if event.Key() == tcell.KeyCtrlC {
					appState.SendQuitSignal()
					keyHandled = true
				}

				for _, elem := range appState.Tools {
					if elem.HandleKey(event, appState) {
						keyHandled = true
						break
					}
				}

				if !keyHandled {
					appState.Machine.Handler(event, appState.EditorBuffer)
				}

			}
		case <-appState.QuitChn:
			quit()
		}

		app.DrawAppState(screen, appState)

	}

}
