// Package ui contains the GTK frontend of the launcher. It only builds the
// widgets; buttons are not wired to any backend logic yet.
package ui

import "github.com/diamondburned/gotk4/pkg/gtk/v4"

const (
	pageMain     = "main"
	pageSettings = "settings"
)

// NewWindow builds the launcher window: a stack holding the main and settings
// pages, with the status bar below. It uses the platform's default decorations.
func NewWindow(app *gtk.Application) *gtk.ApplicationWindow {
	stack := gtk.NewStack()
	stack.SetTransitionType(gtk.StackTransitionTypeSlideLeftRight)
	stack.SetVExpand(true)
	stack.AddNamed(newMainView(func() { stack.SetVisibleChildName(pageSettings) }), pageMain)
	stack.AddNamed(newSettingsView(func() { stack.SetVisibleChildName(pageMain) }), pageSettings)
	stack.SetVisibleChildName(pageMain)

	content := gtk.NewBox(gtk.OrientationVertical, 0)
	content.Append(stack)
	content.Append(gtk.NewSeparator(gtk.OrientationHorizontal))
	content.Append(newStatusBar())

	window := gtk.NewApplicationWindow(app)
	window.SetTitle("GAMMA Launcher")
	window.SetDefaultSize(1000, 640)
	window.SetChild(content)
	return window
}

// newStatusBar shows the addon status line shared by both pages.
func newStatusBar() *gtk.Box {
	bar := gtk.NewBox(gtk.OrientationHorizontal, 24)
	setMargins(bar, 6)
	bar.SetMarginStart(12)
	bar.SetMarginEnd(12)
	for _, text := range []string{
		"ModDB addons status: Up-to-date",
		"GAMMA Addon status: Up-to-date",
	} {
		label := gtk.NewLabel(text)
		label.AddCSSClass("dim-label")
		bar.Append(label)
	}
	return bar
}
