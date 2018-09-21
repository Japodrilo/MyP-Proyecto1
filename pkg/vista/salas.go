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
	win := SetupPopupWindow("Nombre ocupado", 250, 48)
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

func SalaInexistente() {
	win := SetupPopupWindow("Sala Inexistente", 250, 48)
	box := SetupBox()
	grid := SetupGrid(gtk.ORIENTATION_HORIZONTAL)
	mensaje := SetupLabel("No existe esa sala")
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

func NoTePertenece() {
	win := SetupPopupWindow("Ella no te pertenece", 270, 48)
	box := SetupBox()
	grid := SetupGrid(gtk.ORIENTATION_HORIZONTAL)
	mensaje := SetupLabel("No eres el dueño de la sala")
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

func UsuarioInexistente() {
	win := SetupPopupWindow("Usuario Inexistente", 250, 48)
	box := SetupBox()
	grid := SetupGrid(gtk.ORIENTATION_HORIZONTAL)
	mensaje := SetupLabel("No existe ese usuario")
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

func EntrarSala() {
	win := SetupPopupWindow("Unirse a sala", 250, 48)
	box := SetupBox()
	grid := SetupGrid(gtk.ORIENTATION_HORIZONTAL)
	mensaje := SetupLabel("No existe ese usuario")
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
	Win       *gtk.Window
	SalaCBT   *gtk.ComboBoxText
	NombreCBT *gtk.ComboBoxText
	InvitarB  *gtk.Button
}

func NuevaInvitar() *Invitar {
	win := SetupPopupWindow("Invitar", 245, 115)
	grid := SetupGrid(gtk.ORIENTATION_VERTICAL)

	esquinaNW := SetupLabel("    ")
	salaL := SetupLabel("Sala:")
	salaCBT := SetupComboBoxText()
	nombreL := SetupLabel("Usuario:")
	nombreCBT := SetupComboBoxText()
	esquinaSE := SetupLabel("    ")

	salaCBT.SetHExpand(true)
	nombreCBT.SetHExpand(true)

	invitar := SetupButton("Invitar")

	grid.Add(esquinaNW)
	grid.Attach(salaL, 1, 1, 1, 1)
	grid.Attach(nombreL, 1, 2, 1, 1)
	grid.Attach(salaCBT, 2, 1, 1, 1)
	grid.Attach(nombreCBT, 2, 2, 1, 1)
	grid.Attach(esquinaSE, 3, 3, 1, 1)

	grid.Attach(invitar, 1, 3, 2, 1)

	win.Add(grid)
	win.ShowAll()

	return &Invitar{
		Win: win,
		SalaCBT: salaCBT,
		NombreCBT: nombreCBT,
		InvitarB: invitar,
	}
}

type MisSalas struct {
	Win     *gtk.Window
	SalaCBT *gtk.ComboBoxText
	AbrirB  *gtk.Button
}

func NuevaMisSalas() *MisSalas {
	win := SetupPopupWindow("Mis salas", 300, 100)
	grid := SetupGrid(gtk.ORIENTATION_VERTICAL)

	esquinaNW := SetupLabel("    ")
	salaL := SetupLabel("Sala:")
	salaCBT := SetupComboBoxText()
	esquinaSE := SetupLabel("    ")

	salaCBT.SetHExpand(true)

	abrir := SetupButton("Abrir conversación")

	grid.Add(esquinaNW)
	grid.Attach(salaL, 1, 1, 1, 1)
	grid.Attach(salaCBT, 2, 1, 1, 1)
	grid.Attach(esquinaSE, 3, 2, 1, 1)

	grid.Attach(abrir, 1, 3, 2, 1)

	win.Add(grid)
	win.ShowAll()

	return &MisSalas{
		Win: win,
		SalaCBT: salaCBT,
		AbrirB: abrir,
	}
}