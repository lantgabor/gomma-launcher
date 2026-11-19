package ui

import "github.com/diamondburned/gotk4/pkg/gtk/v4"

func newSettingsView(goBack func()) *gtk.Box {
	// Paths: either one installation root (subfolders are created inside it),
	// or, in advanced mode, every path set individually. Only one is shown.
	installGrid := newFormGrid()
	addPathRow(installGrid, 0, "Installation path", "~/Games/GAMMA")

	advancedGrid := newFormGrid()
	for i, field := range []struct{ name, example string }{
		{"Anomaly path", "~/Games/Anomaly"},
		{"GAMMA path", "~/Games/GAMMA"},
		{"Cache path", "~/.cache/gomma-launcher"},
		{"Launcher DB path", "~/.local/share/gomma-launcher/launcher.db"},
	} {
		addPathRow(advancedGrid, i, field.name, field.example)
	}
	advancedGrid.SetVisible(false)

	advanced := gtk.NewCheckButtonWithLabel("Advanced: set each path individually")
	advanced.ConnectToggled(func() {
		installGrid.SetVisible(!advanced.Active())
		advancedGrid.SetVisible(advanced.Active())
	})

	paths := gtk.NewBox(gtk.OrientationVertical, 12)
	paths.Append(advanced)
	paths.Append(installGrid)
	paths.Append(advancedGrid)

	// Cache maintenance.
	purge := gtk.NewButtonWithLabel("Purge Cache")
	purge.AddCSSClass("destructive-action")
	purge.SetHAlign(gtk.AlignStart)

	more := gtk.NewLabel("More settings later here...")
	more.SetXAlign(0)
	more.AddCSSClass("dim-label")

	form := gtk.NewBox(gtk.OrientationVertical, 18)
	setMargins(form, 18)
	form.Append(newSection("Paths", paths))
	form.Append(newSection("Cache", purge))
	form.Append(gtk.NewSeparator(gtk.OrientationHorizontal))
	form.Append(more)

	scroller := gtk.NewScrolledWindow()
	scroller.SetPolicy(gtk.PolicyNever, gtk.PolicyAutomatic)
	scroller.SetVExpand(true)
	scroller.SetChild(form)

	save := gtk.NewButtonWithLabel("Save Settings")
	save.AddCSSClass("suggested-action")

	back := gtk.NewButtonWithLabel("Back")
	back.ConnectClicked(goBack)

	actions := gtk.NewActionBar()
	actions.PackStart(back)
	actions.PackEnd(save)

	view := gtk.NewBox(gtk.OrientationVertical, 0)
	view.Append(scroller)
	view.Append(actions)
	return view
}

// newSection wraps content in a titled frame.
func newSection(title string, child gtk.Widgetter) *gtk.Frame {
	inner := gtk.NewBox(gtk.OrientationVertical, 0)
	setMargins(inner, 12)
	inner.Append(child)

	frame := gtk.NewFrame(title)
	frame.SetChild(inner)
	return frame
}

func newFormGrid() *gtk.Grid {
	grid := gtk.NewGrid()
	grid.SetColumnSpacing(12)
	grid.SetRowSpacing(8)
	return grid
}

func addPathRow(grid *gtk.Grid, row int, name, example string) {
	label := gtk.NewLabel(name)
	label.SetXAlign(0)

	entry := gtk.NewEntry()
	entry.SetPlaceholderText(example)
	entry.SetHExpand(true)

	grid.Attach(label, 0, row, 1, 1)
	grid.Attach(entry, 1, row, 1, 1)
}
