package ui

import (
	"context"
	"fmt"

	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

// supportLinks are the entries of the Support menu.
// TODO: confirm the GAMMA donation link; it points to the GAMMA project page for now.
var supportLinks = []struct{ label, url string }{
	{"Support GAMMA: Donate", "https://github.com/Grokitach/Stalker_GAMMA"},
	{"Support ModDB", "https://www.moddb.com/"},
	{"Support GSC Game World", "https://www.gsc-game.com/"},
	{"Support modders", "https://www.moddb.com/mods/stalker-anomaly/addons"},
}

// newSupportMenu is a standard GTK menu whose items open supportLinks in the
// browser, through the window's "open-url" action.
func newSupportMenu() *gtk.MenuButton {
	menu := gio.NewMenu()
	for _, link := range supportLinks {
		menu.Append(link.label, fmt.Sprintf("win.open-url('%s')", link.url))
	}

	button := gtk.NewMenuButton()
	button.SetLabel("Support")
	button.SetDirection(gtk.ArrowUp)
	button.SetMenuModel(menu)
	return button
}

// addOpenURLAction adds "win.open-url", which opens its string parameter (a
// URL) with the system's default handler.
func addOpenURLAction(window *gtk.ApplicationWindow) {
	action := gio.NewSimpleAction("open-url", glib.NewVariantType("s"))
	action.ConnectActivate(func(url *glib.Variant) {
		gtk.NewURILauncher(url.String()).Launch(context.Background(), &window.Window, nil)
	})
	window.AddAction(action)
}
