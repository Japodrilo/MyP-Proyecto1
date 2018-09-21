package vista

import (
	"github.com/gotk3/gotk3/gtk"
)

type Identificar struct {
	Win           *gtk.Window
	Nombre        *gtk.Entry
	IdentificarB  *gtk.Button
}

func NuevaIdentificar() *Identificar {
	win := SetupPopupWindow("Identifícate", 330, 90)
	grid := SetupGrid(gtk.ORIENTATION_VERTICAL)

	esquina := SetupLabel("    ")
	nombreL := SetupLabel("Nombre de Usuario:")
	nombreE := SetupEntry()

	conectar := SetupButton("Identificarse")

	grid.Add(esquina)
	grid.Attach(nombreL, 1, 1, 1, 1)
	grid.Attach(nombreE, 2, 1, 1, 1)

	grid.Attach(conectar, 1, 3, 2, 1)

	win.Add(grid)
	win.ShowAll()

	return &Identificar{
		Win: win,
		Nombre: nombreE,
		IdentificarB: conectar,
	}
}

func NombreOcupado() {
	win := SetupPopupWindow("Nombre ocupado", 500, 48)
	box := SetupBox()
	grid := SetupGrid(gtk.ORIENTATION_HORIZONTAL)
	mensaje := SetupLabel("Nombre ocupado, intenta nuevamente")
	espacio1 := SetupLabel("    ")
	espacio1.SetHExpand(true)
	espacio2 := SetupLabel("    ")
	espacio2.SetHExpand(true)
	desconectar := SetupButtonClick("Cerrar", func() {win.Close()})
	box.Add(mensaje)
	grid.Add(espacio1)
	grid.Add(desconectar)
	grid.Add(espacio2)
	box.Add(grid)
	win.Add(box)
	win.ShowAll()
}