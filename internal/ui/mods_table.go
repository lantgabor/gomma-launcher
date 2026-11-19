package ui

import (
	"cmp"
	"strconv"
	"strings"
	"unsafe"

	coreglib "github.com/diamondburned/gotk4/pkg/core/glib"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/diamondburned/gotk4/pkg/pango"
)

// tableColumn describes one column of the mods table.
type tableColumn struct {
	title   string
	expand  bool
	newCell func() gtk.Widgetter        // builds an empty cell widget
	bind    func(gtk.Widgetter, modRow) // fills a cell for a mod
	compare func(a, b modRow) int       // sort order for the column
}

// newModsTable shows the mods in a ColumnView. Its header stays visible while
// scrolling, and clicking a column header sorts by that column (click again to
// flip ascending/descending). Rows start in load order, which the "#" column
// shows and can sort back to.
func newModsTable(mods []modRow) *gtk.ColumnView {
	// The list model holds only each mod's index in mods, as a string.
	// Cells and sorters look the mod up by that index.
	indices := make([]string, len(mods))
	for i := range mods {
		indices[i] = strconv.Itoa(i)
	}
	sorted := gtk.NewSortListModel(gtk.NewStringList(indices), nil)

	table := gtk.NewColumnView(gtk.NewNoSelection(sorted))
	table.SetShowRowSeparators(true)
	sorted.SetSorter(table.Sorter())

	for _, col := range modColumns {
		factory := gtk.NewSignalListItemFactory()
		factory.ConnectSetup(func(obj *coreglib.Object) {
			listItem(obj).SetChild(col.newCell())
		})
		factory.ConnectBind(func(obj *coreglib.Object) {
			item := listItem(obj)
			i := modIndex(item.Item())
			col.bind(item.Child(), mods[i])
		})

		sorter := gtk.NewCustomSorter(func(a, b unsafe.Pointer) int {
			return col.compare(mods[modIndex(coreglib.Take(a))], mods[modIndex(coreglib.Take(b))])
		})

		column := gtk.NewColumnViewColumn(col.title, &factory.ListItemFactory)
		column.SetExpand(col.expand)
		column.SetSorter(&sorter.Sorter)
		table.AppendColumn(column)
	}
	return table
}

var modColumns = []tableColumn{
	{
		title:   "#",
		newCell: newCellLabel,
		bind: func(w gtk.Widgetter, mod modRow) {
			w.(*gtk.Label).SetText(strconv.Itoa(mod.LoadOrder))
		},
		compare: func(a, b modRow) int { return cmp.Compare(a.LoadOrder, b.LoadOrder) },
	},
	{
		title:   "Name",
		expand:  true,
		newCell: newCellLabel,
		bind: func(w gtk.Widgetter, mod modRow) {
			w.(*gtk.Label).SetText(mod.Name)
			setClass(w, "error", mod.Status == "Error")
		},
		compare: func(a, b modRow) int { return strings.Compare(a.Name, b.Name) },
	},
	{
		title:   "Downloading",
		newCell: newCellProgressBar,
		bind: func(w gtk.Widgetter, mod modRow) {
			bindProgress(w.(*gtk.ProgressBar), mod.Downloading)
		},
		compare: func(a, b modRow) int { return cmp.Compare(a.Downloading, b.Downloading) },
	},
	{
		title:   "Extracting",
		newCell: newCellProgressBar,
		bind: func(w gtk.Widgetter, mod modRow) {
			bindProgress(w.(*gtk.ProgressBar), mod.Extracting)
		},
		compare: func(a, b modRow) int { return cmp.Compare(a.Extracting, b.Extracting) },
	},
	{
		title:   "Status",
		newCell: newCellLabel,
		bind: func(w gtk.Widgetter, mod modRow) {
			w.(*gtk.Label).SetText(mod.Status)
			setClass(w, "warning", mod.Status == "Mismatch")
			setClass(w, "error", mod.Status == "Error")
		},
		compare: func(a, b modRow) int { return strings.Compare(a.Status, b.Status) },
	},
}

// listItem wraps a factory callback's object. With GTK 4.12+ it is really a
// GtkColumnViewCell (a GtkListItem subclass) that gotk4 v0.4.1 does not know.
func listItem(obj *coreglib.Object) *gtk.ListItem {
	return &gtk.ListItem{Object: obj}
}

// modIndex reads the index stored in a list model item (a GtkStringObject).
func modIndex(obj *coreglib.Object) int {
	i, _ := strconv.Atoi((&gtk.StringObject{Object: obj}).String())
	return i
}

func newCellLabel() gtk.Widgetter {
	label := gtk.NewLabel("")
	label.SetXAlign(0)
	label.SetEllipsize(pango.EllipsizeEnd)
	return label
}

func newCellProgressBar() gtk.Widgetter {
	bar := gtk.NewProgressBar()
	bar.SetShowText(true)
	bar.SetVAlign(gtk.AlignCenter)
	bar.SetSizeRequest(160, -1)
	return bar
}

// bindProgress shows the bar only when there is progress to show.
func bindProgress(bar *gtk.ProgressBar, fraction float64) {
	bar.SetVisible(fraction > 0)
	bar.SetFraction(fraction)
}
