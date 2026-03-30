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
	appState := &app.App{ScreenEventChan: eventChannel, Machine: machine.VimMachine{Mode: &machine.NormalMode{}}}

	for {
		app.DrawAppState(screen, appState)

		event := <-eventChannel
		if ev, ok := event.(*tcell.EventKey); ok {
			appState.Machine.Mode.KeyHandler(ev, appState)
		}
	}

}
