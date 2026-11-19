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

// setClass adds or removes a style class. List cells are reused for different
// rows, so a class must also be removed when it no longer applies.
func setClass(w gtk.Widgetter, class string, on bool) {
	if on {
		gtk.BaseWidget(w).AddCSSClass(class)
	} else {
		gtk.BaseWidget(w).RemoveCSSClass(class)
	}
}
