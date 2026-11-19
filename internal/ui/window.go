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
// mods fills the mods table; version is shown in the bottom-right corner.
func NewWindow(app *gtk.Application, mods []Mod, version string) *gtk.ApplicationWindow {
	window := gtk.NewApplicationWindow(app)
	window.SetTitle("GAMMA Launcher")
	window.SetDefaultSize(1000, 640)
	addOpenURLAction(window)

	stack := gtk.NewStack()
	stack.SetTransitionType(gtk.StackTransitionTypeSlideLeftRight)
	stack.SetVExpand(true)
	stack.AddNamed(newMainView(mods, func() { stack.SetVisibleChildName(pageSettings) }), pageMain)
	stack.AddNamed(newSettingsView(func() { stack.SetVisibleChildName(pageMain) }), pageSettings)
	stack.SetVisibleChildName(pageMain)

	content := gtk.NewBox(gtk.OrientationVertical, 0)
	content.Append(stack)
	content.Append(gtk.NewSeparator(gtk.OrientationHorizontal))
	content.Append(newStatusBar(version))
	window.SetChild(content)
	return window
}

// newStatusBar shows the addon status on the left and the version on the right.
func newStatusBar(version string) *gtk.CenterBox {
	status := gtk.NewBox(gtk.OrientationHorizontal, 24)
	for _, text := range []string{
		"ModDB addons status: Up-to-date",
		"GAMMA Addon status: Up-to-date",
	} {
		status.Append(newDimLabel(text))
	}

	bar := gtk.NewCenterBox()
	setMargins(bar, 6)
	bar.SetMarginStart(12)
	bar.SetMarginEnd(12)
	bar.SetStartWidget(status)
	bar.SetEndWidget(newDimLabel(version))
	return bar
}

func newDimLabel(text string) *gtk.Label {
	label := gtk.NewLabel(text)
	label.AddCSSClass("dim-label")
	return label
}
