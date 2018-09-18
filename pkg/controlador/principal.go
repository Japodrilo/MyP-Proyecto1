package controlador

import(
	"github.com/Japodrilo/MyP-Proyecto1/pkg/vista"

	"github.com/gotk3/gotk3/gtk"
)

type Principal struct {
	win			*gtk.Window
	cuaderno	*Cuaderno
	menu		*Menu
	lb			*gtk.ListBox
}

func NuevaPrincipal() *Principal {
	ventanaPrincipal := vista.SetupVentanaPrincipal()
	cuaderno := NuevoCuaderno(ventanaPrincipal.Nb)
	menu := NuevoMenu(ventanaPrincipal.Menubar)

	cuaderno.entradas["General"], cuaderno.textos["General"] = cuaderno.AddTab("General")

	return &Principal{
		win:		ventanaPrincipal.Win,
		cuaderno:	cuaderno,
		menu:		menu,
		lb:			ventanaPrincipal.Lb,
	}
}

func (principal *Principal) AddUser(username string) {
	lbr := vista.SetupListBoxRow()
	entry := vista.SetupEntry()
	btn := vista.SetupButtonUser(username, principal.cuaderno.nb, principal.cuaderno.textos, entry)
	lbr.Add(btn)
	principal.lb.Add(lbr)
	principal.win.ShowAll()
}