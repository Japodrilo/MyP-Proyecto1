package vista

import(
	//"fmt"
	"log"

	"github.com/gotk3/gotk3/gtk"
)


func SetupWindow(title string) *gtk.Window {
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}
	win.SetTitle(title)
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})
	win.SetDefaultSize(800, 600)
	win.SetPosition(gtk.WIN_POS_CENTER)
	return win
}

func SetupPopupWindow(title string, width, height int) *gtk.Window {
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}
	win.SetTitle(title)
	win.Connect("destroy", func() {
		win.Close()
	})
	win.SetDefaultSize(width, height)
	win.SetPosition(gtk.WIN_POS_CENTER)
	return win
}

func SetupBox() *gtk.Box {
	box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		log.Fatal("Unable to create box:", err)
	}
	box.SetHomogeneous(false)
	return box
}

func SetupNotebook() *gtk.Notebook {
	nb, err := gtk.NotebookNew()
	if err != nil {
		log.Fatal("Unable to create notebook:", err)
	}
	return nb
}

func SetupScrolledWindow() *gtk.ScrolledWindow {
	scrwin, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		log.Fatal("Unable to create scrolled window:", err)
	}
	scrwin.SetPolicy(1,1)
	scrwin.SetHExpand(true)
	return scrwin
}

func SetupGrid(orient gtk.Orientation) *gtk.Grid {
	grid, err := gtk.GridNew()
	if err != nil {
		log.Fatal("Unable to create grid:", err)
	}
	grid.SetOrientation(orient)
	return grid
}

func SetupTextView() *gtk.TextView {
	tv, err := gtk.TextViewNew()
	if err != nil {
		log.Fatal("Unable to create TextView:", err)
	}
	tv.SetWrapMode(2)
	tv.SetEditable(false)
	tv.SetCursorVisible(false)
	return tv
}

func SetupEntry() *gtk.Entry {
	entry, err := gtk.EntryNew()
	if err != nil {
		log.Fatal("Unable to create Entry:", err)
	}
	return entry
}

func SetupButton(label string) *gtk.Button {
	btn, err := gtk.ButtonNewWithLabel(label)
	if err != nil {
		log.Fatal("Unable to create button:", err)
	}
	return btn
}

func SetupButtonClick(label string, onClick func()) *gtk.Button {
	btn := SetupButton(label)
	btn.Connect("clicked", onClick)
	return btn
}

func setup_list_box() *gtk.ListBox {
	lb, err := gtk.ListBoxNew()
	if err != nil {
		log.Fatal("Unable to create ListBox:", err)
	}
	return lb
}

func SetupLabel(text string) *gtk.Label {
	label, err := gtk.LabelNew(text)
	if err != nil {
		log.Fatal("Unable to create label:", err)
	}
	return label
}

func GetLabelText(label *gtk.Label) string {
	text, err := label.GetText()
	if err != nil {
		log.Fatal("Unable to get text from label:", err)
	}
	return text
}

func GetBufferTV(tv *gtk.TextView) *gtk.TextBuffer {
	buffer, err := tv.GetBuffer()
	if err != nil {
		log.Fatal("Unable to get buffer:", err)
	}
	return buffer
}

func get_buffer_from_entry(entry *gtk.Entry) *gtk.EntryBuffer {
	buffer, err := entry.GetBuffer()
	if err != nil {
		log.Fatal("Unable to get buffer:", err)
	}
	return buffer
}

func GetTextEntry(entry *gtk.Entry) string {
	text, err := entry.GetText()
	if err != nil {
		log.Fatal("Unable to get text from buffer:", err)
	}
	return text
}

func GetTextEntryClean(entry *gtk.Entry) string {
	text, err := entry.GetText()
	if err != nil {
		log.Fatal("Unable to get text from buffer:", err)
	}
	buffer := get_buffer_from_entry(entry)
	buffer.DeleteText(0,-1)
	return text + "\n"
}

func SetupListBoxRow() *gtk.ListBoxRow {
	lbr, err := gtk.ListBoxRowNew()
	if err != nil {
		log.Fatal("Unable to create List Box Row:", err)
	}
	return lbr
}

func setup_menu_bar() *gtk.MenuBar {
	menubar, err := gtk.MenuBarNew()
	if err != nil {
		log.Fatal("Unable to create Menu Bar:", err)
	}
	return menubar
}

func setup_menu() *gtk.Menu {
	menu, err := gtk.MenuNew()
	if err != nil {
		log.Fatal("Unable to create Menu:", err)
	}
	return menu
}

func setup_menu_item_label(text string) *gtk.MenuItem {
	menuitem, err := gtk.MenuItemNewWithLabel(text)
	if err != nil {
		log.Fatal("Unable to create Menu Bar:", err)
	}
	return menuitem
}

func setup_menu_item_mnemonic(text string) *gtk.MenuItem {
	menuitem, err := gtk.MenuItemNewWithMnemonic(text)
	if err != nil {
		log.Fatal("Unable to create Menu Bar:", err)
	}
	return menuitem
}
