package main

import (
	"log/slog"
	"os"

	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"

	"github.com/lantgabor/gomma-launcher/internal/ui"
)

// Set at build time by the Makefile via -ldflags "-X main.version=... -X main.commit=...".
var (
	version = "dev"
	commit  = ""
)

func main() {
	// gotk4 forwards GLib/GTK logs to slog; hide INFO-level noise (e.g. Vulkan loader messages).
	slog.SetLogLoggerLevel(slog.LevelWarn)

	versionText := "Version " + version
	if commit != "" {
		versionText += " (" + commit + ")"
	}

	app := gtk.NewApplication("io.github.lantgabor.GommaLauncher", gio.ApplicationDefaultFlags)
	app.ConnectActivate(func() {
		ui.NewWindow(app, ui.PlaceholderMods(), versionText).SetVisible(true)
	})

	if code := app.Run(os.Args); code > 0 {
		os.Exit(code)
	}
}
