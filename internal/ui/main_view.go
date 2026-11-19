package ui

import "github.com/diamondburned/gotk4/pkg/gtk/v4"

// newMainView is the main page: action sidebar on the left, the mods table on
// the right, and the Play / Delete Shader Cache bar below the table.
func newMainView(openSettings func()) *gtk.Box {
	view := gtk.NewBox(gtk.OrientationHorizontal, 0)
	view.Append(newSidebar(openSettings))
	view.Append(gtk.NewSeparator(gtk.OrientationVertical))
	view.Append(newModsPanel())
	return view
}

// newSidebar keeps the actions at the top and the support menu at the bottom.
func newSidebar(openSettings func()) *gtk.CenterBox {
	settings := gtk.NewButtonWithLabel("Launcher Settings")
	settings.ConnectClicked(openSettings)

	actions := gtk.NewBox(gtk.OrientationVertical, 6)
	actions.Append(gtk.NewButtonWithLabel("Install/Update G.A.M.M.A."))
	actions.Append(gtk.NewButtonWithLabel("Sync Version Metadata"))
	actions.Append(settings)

	support := gtk.NewBox(gtk.OrientationVertical, 6)
	support.Append(gtk.NewSeparator(gtk.OrientationHorizontal))
	support.Append(newSupportMenu())

	sidebar := gtk.NewCenterBox()
	sidebar.SetOrientation(gtk.OrientationVertical)
	sidebar.SetSizeRequest(220, -1)
	setMargins(sidebar, 12)
	sidebar.SetStartWidget(actions)
	sidebar.SetEndWidget(support)
	return sidebar
}

func newModsPanel() *gtk.Box {
	scroller := gtk.NewScrolledWindow()
	scroller.SetPolicy(gtk.PolicyNever, gtk.PolicyAutomatic)
	scroller.SetVExpand(true)
	scroller.SetChild(newModsTable(loadPlaceholderMods()))

	frame := gtk.NewFrame("Mods")
	setMargins(frame, 12)
	frame.SetChild(scroller)

	play := gtk.NewButtonWithLabel("Play")
	play.AddCSSClass("suggested-action")
	deleteCache := gtk.NewButtonWithLabel("Delete Shader Cache")
	deleteCache.AddCSSClass("destructive-action")

	actions := gtk.NewActionBar()
	actions.PackEnd(play)
	actions.PackEnd(deleteCache)

	panel := gtk.NewBox(gtk.OrientationVertical, 0)
	panel.SetHExpand(true)
	panel.Append(frame)
	panel.Append(actions)
	return panel
}
