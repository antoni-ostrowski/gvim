package main

import (
	"log"
	"os"

	"github.com/antoni-ostrowski/gvim/internal/app"
	"github.com/antoni-ostrowski/gvim/internal/render"
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
	screen.EnablePaste()
	screen.Clear()
	app := &app.App{Message: "stintsr"}

	for {
		render.DrawAppState(screen, app)

		ev := <-screen.EventQ()
		if ev, ok := ev.(*tcell.EventKey); ok {
			// Normal mode - handle main screen input
			switch ev.Key() {
			case tcell.KeyCtrlC:
				screen.Fini()
				os.Exit(0)
			case tcell.KeyRune:
				os.Exit(0)
			}
		}
	}

}
