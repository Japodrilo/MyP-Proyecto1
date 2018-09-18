package vista

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
)

type Conectar func(direccion *gtk.Entry, puerto *gtk.Entry, win *gtk.Window) func()

func PopUpConectar(funcionConectar Conectar) func() {
	return func() {
		win := SetupPopupWindow("Conectar", 273, 115)
		grid := SetupGrid(gtk.ORIENTATION_VERTICAL)

		esquina := SetupLabel("    ")
		direccionL := SetupLabel("Dirección IP:")
		direccionE := SetupEntry()
		puertoL := SetupLabel("Puerto:")
		puertoE := SetupEntry()

		conectar := SetupButton("Conectar")
		conectar.Connect("clicked", funcionConectar(direccionE, puertoE, win))

		grid.Add(esquina)
		grid.Attach(direccionL, 1, 1, 1, 1)
		grid.Attach(puertoL, 1, 2, 1, 1)
		grid.Attach(direccionE, 2, 1, 1, 1)
		grid.Attach(puertoE, 2, 2, 1, 1)

		grid.Attach(conectar, 1, 3, 2, 1)

		win.Add(grid)
		win.ShowAll()
	}
}

func PopUpErrorConexion(servidor, puerto string) {
	win := SetupPopupWindow("Error", 500, 48)
	box := SetupBox()
	grid := SetupGrid(gtk.ORIENTATION_HORIZONTAL)
	mensaje := SetupLabel(fmt.Sprintf("No fue posible establecer la conexión con \"%v:%v\"", servidor, puerto))
	espacio1 := SetupLabel("    ")
	espacio1.SetHExpand(true)
	espacio2 := SetupLabel("    ")
	espacio2.SetHExpand(true)
	aceptar := SetupButtonClick("Aceptar", func() {
		win.Close()
	})
	box.Add(mensaje)
	grid.Add(espacio1)
	grid.Add(aceptar)
	grid.Add(espacio2)
	box.Add(grid)
	win.Add(box)
	win.ShowAll()
}

func PopUpDesconectar() {
	win := SetupPopupWindow("Desconectar", 500, 48)
	box := SetupBox()
	grid := SetupGrid(gtk.ORIENTATION_HORIZONTAL)
	mensaje := SetupLabel("¿Deseas terminar tu conexión con el servidor?")
	espacio1 := SetupLabel("    ")
	espacio1.SetHExpand(true)
	espacio2 := SetupLabel("    ")
	espacio2.SetHExpand(true)
	aceptar := SetupButtonClick("Desconectar", func() {
		win.Close()
	})
	box.Add(mensaje)
	grid.Add(espacio1)
	grid.Add(aceptar)
	grid.Add(espacio2)
	box.Add(grid)
	win.Add(box)
	win.ShowAll()
}



