package main

import (
	"log/slog"
	"os"

	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"

	"gofrontend/internal/ui"
)

func main() {
	// gotk4 forwards GLib/GTK logs to slog; hide INFO-level noise (e.g. Vulkan loader messages).
	slog.SetLogLoggerLevel(slog.LevelWarn)

	app := gtk.NewApplication("com.example.gofrontend", gio.ApplicationDefaultFlags)
	app.ConnectActivate(func() { ui.NewWindow(app).SetVisible(true) })

	if code := app.Run(os.Args); code > 0 {
		os.Exit(code)
	}
}
