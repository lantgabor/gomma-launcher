package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func main() {
	// gotk4 forwards GLib/GTK logs to slog; hide INFO-level noise (e.g. Vulkan loader messages).
	slog.SetLogLoggerLevel(slog.LevelInfo)

	app := gtk.NewApplication("com.example.gofrontend", gio.ApplicationFlagsNone)
	app.ConnectActivate(func() { activate(app) })

	if code := app.Run(os.Args); code > 0 {
		os.Exit(code)
	}
}

func activate(app *gtk.Application) {
	label := gtk.NewLabel("Hello, World!")

	button := gtk.NewButtonWithLabel("Click me")
	clicks := 0
	button.ConnectClicked(func() {
		clicks++
		slog.Info("clicks: ", "clicks", clicks)
		if clicks == 1 {
			label.SetText("Clicked 1 time")
		} else {
			label.SetText(fmt.Sprintf("Clicked %d times", clicks))
		}
	})

	box := gtk.NewBox(gtk.OrientationVertical, 12)
	box.SetMarginTop(24)
	box.SetMarginBottom(24)
	box.SetMarginStart(24)
	box.SetMarginEnd(24)
	box.SetVAlign(gtk.AlignCenter)
	box.Append(label)
	box.Append(button)

	window := gtk.NewApplicationWindow(app)
	window.SetTitle("gotk4 Hello World")
	window.SetChild(box)
	window.SetDefaultSize(400, 300)
	window.SetVisible(true)
}
