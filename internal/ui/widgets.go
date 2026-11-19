package ui

import "github.com/diamondburned/gotk4/pkg/gtk/v4"

// setMargins sets the same margin on all four sides of a widget.
func setMargins(w gtk.Widgetter, margin int) {
	widget := gtk.BaseWidget(w)
	widget.SetMarginTop(margin)
	widget.SetMarginBottom(margin)
	widget.SetMarginStart(margin)
	widget.SetMarginEnd(margin)
}
