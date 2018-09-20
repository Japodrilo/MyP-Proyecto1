package vista

import (
	"github.com/gotk3/gotk3/gtk"
)

type Crear struct {
	Win    *gtk.Window
	Nombre *gtk.Entry
	CrearB *gtk.Button
}

func NuevaCrear() *Crear {
	win := SetupPopupWindow("Crear Sala", 300, 90)
	grid := SetupGrid(gtk.ORIENTATION_VERTICAL)

	esquina := SetupLabel("    ")
	nombreL := SetupLabel("Nombre de Sala:")
	nombreE := SetupEntry()

	crear := SetupButton("Crear")

	grid.Add(esquina)
	grid.Attach(nombreL, 1, 1, 1, 1)
	grid.Attach(nombreE, 2, 1, 1, 1)

	grid.Attach(crear, 1, 3, 2, 1)

	win.Add(grid)
	win.ShowAll()

	return &Crear{
		Win:    win,
		Nombre: nombreE,
		CrearB: crear,
	}
}

func NombreSalaOcupado() {
	win := SetupPopupWindow("Nombre ocupado", 400, 48)
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

type Invitar struct {
	Win      *gtk.Window
	SalaE    *gtk.Entry
	NombreE  *gtk.Entry
	InvitarB *gtk.Button
}

func NuevaInvitar() *Invitar {
	win := SetupPopupWindow("Invitar", 245, 115)
	grid := SetupGrid(gtk.ORIENTATION_VERTICAL)

	esquina := SetupLabel("    ")
	salaL := SetupLabel("Sala:")
	salaE := SetupEntry()
	nombreL := SetupLabel("Usuario:")
	nombreE := SetupEntry()

	invitar := SetupButton("Invitar")

	grid.Add(esquina)
	grid.Attach(salaL, 1, 1, 1, 1)
	grid.Attach(nombreL, 1, 2, 1, 1)
	grid.Attach(salaE, 2, 1, 1, 1)
	grid.Attach(nombreE, 2, 2, 1, 1)

	grid.Attach(invitar, 1, 3, 2, 1)

	win.Add(grid)
	win.ShowAll()

	return &Invitar{
		Win: win,
		SalaE: salaE,
		NombreE: nombreE,
		InvitarB: invitar,
	}
}