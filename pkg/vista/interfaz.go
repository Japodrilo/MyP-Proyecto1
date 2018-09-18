package vista

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
)

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





