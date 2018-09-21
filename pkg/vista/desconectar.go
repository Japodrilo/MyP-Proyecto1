package vista

import (
	"github.com/gotk3/gotk3/gtk"
)

type Desconectar struct {
	Win           *gtk.Window
	DesconectarB  *gtk.Button
}

func NuevaDesconectar() *Desconectar {
	win := SetupPopupWindow("Desconectar", 345, 48)
	box := SetupBox()
	grid := SetupGrid(gtk.ORIENTATION_HORIZONTAL)
	mensaje := SetupLabel("¿Deseas terminar tu conexión con el servidor?")
	espacio1 := SetupLabel("    ")
	espacio1.SetHExpand(true)
	espacio2 := SetupLabel("    ")
	espacio2.SetHExpand(true)
	desconectar := SetupButton("Desconectar")
	box.Add(mensaje)
	grid.Add(espacio1)
	grid.Add(desconectar)
	grid.Add(espacio2)
	box.Add(grid)
	win.Add(box)
	win.ShowAll()

	return &Desconectar{
		Win: win,
		DesconectarB: desconectar,
	}
}
