package vista

import (
	"github.com/gotk3/gotk3/gtk"
)

/**
 * Estructura que representa a la ventana principal y 
 * sus campos de interés para el controlador.
 */
type VentanaPrincipal struct {
	Win			*gtk.Window
	Nb			*gtk.Notebook
	Menubar 	*MenuPrincipal
	Lb 			*gtk.ListBox
}

/**
 * Constructor, función que dibuja la ventana principal
 * e inicializa los campos que necesita el controlador.
 */
func SetupVentanaPrincipal() *VentanaPrincipal {
	win := SetupWindow("Chat")
	box := SetupBox()
	menubar := SetupMenuPrincipal()
	grid := SetupGrid(gtk.ORIENTATION_HORIZONTAL)
	scrwinusr := SetupScrolledWindow()
	nb := SetupNotebook()

	lb := SetupListBox()
	scrwinusr.Add(lb)

	box.Add(menubar.Menubar)
	box.Add(grid)

	grid.Attach(nb, 0, 0, 2, 2)
	grid.Attach(scrwinusr, 2, 0, 1, 2)

	win.Add(box)
	win.ShowAll()

	return &VentanaPrincipal{
		Win:		win,
		Nb:			nb,
		Menubar:	menubar,
		Lb:			lb,
	}
}
