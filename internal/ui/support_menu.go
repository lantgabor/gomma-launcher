package ui

import "github.com/diamondburned/gotk4/pkg/gtk/v4"

// newSupportMenu groups the "Support ..." links into a menu button.
func newSupportMenu() *gtk.MenuButton {
	list := gtk.NewBox(gtk.OrientationVertical, 0)
	for _, text := range []string{
		"Support GAMMA: Donate",
		"Support ModDB",
		"Support GSC Game World",
		"Support modders",
	} {
		button := gtk.NewButtonWithLabel(text)
		button.AddCSSClass("flat")
		list.Append(button)
	}

	popover := gtk.NewPopover()
	popover.SetChild(list)

	menu := gtk.NewMenuButton()
	menu.SetLabel("Support")
	menu.SetDirection(gtk.ArrowUp)
	menu.SetPopover(popover)
	return menu
}
