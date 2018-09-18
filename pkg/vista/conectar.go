package vista

import (
	"github.com/gotk3/gotk3/gtk"
)

type Conectar struct {
	Win        *gtk.Window
	DireccionE *gtk.Entry
	PuertoE    *gtk.Entry
	ConectarB  *gtk.Button
}

func NuevaConectar() *Conectar {
	win := SetupPopupWindow("Conectar", 273, 115)
	grid := SetupGrid(gtk.ORIENTATION_VERTICAL)

	esquina := SetupLabel("    ")
	direccionL := SetupLabel("Dirección IP:")
	direccionE := SetupEntry()
	puertoL := SetupLabel("Puerto:")
	puertoE := SetupEntry()

	conectar := SetupButton("Conectar")

	grid.Add(esquina)
	grid.Attach(direccionL, 1, 1, 1, 1)
	grid.Attach(puertoL, 1, 2, 1, 1)
	grid.Attach(direccionE, 2, 1, 1, 1)
	grid.Attach(puertoE, 2, 2, 1, 1)

	grid.Attach(conectar, 1, 3, 2, 1)

	win.Add(grid)
	win.ShowAll()

	return &Conectar{
		Win: win,
		DireccionE: direccionE,
		PuertoE: puertoE,
		ConectarB: conectar,
	}
}