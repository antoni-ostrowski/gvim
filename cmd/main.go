package main

import (
	"log"
	"os"

	"github.com/antoni-ostrowski/gvim/internal/app"
	"github.com/antoni-ostrowski/gvim/internal/machine"
	"github.com/gdamore/tcell/v3"
)

func main() {
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := screen.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	defer func() {
		screen.Fini()
		os.Exit(0)
	}()
	screen.EnablePaste()
	screen.Clear()
	eventChannel := screen.EventQ()
	appState := &app.App{
		Machine: machine.VimMachine{Mode: &machine.NormalMode{}},
		QuitChn: make(chan struct{}, 1),
		Screen:  screen,
	}

	for {
		select {
		case event := <-eventChannel:
			if event, ok := event.(*tcell.EventKey); ok {
				keyHandled := false
				for _, elem := range appState.UiElements {
					if elem.HandleKey(event, appState) {
						keyHandled = true
						break
					}
				}

				if !keyHandled {
					appState.Machine.Handler(event, appState)
				}

			}
		case <-appState.QuitChn:
			screen.Fini()
			os.Exit(0)
		}

		app.DrawAppState(screen, appState)

	}

}
