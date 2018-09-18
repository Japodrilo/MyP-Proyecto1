package controlador

import(
	"github.com/Japodrilo/MyP-Proyecto1/pkg/vista"

	//"github.com/gotk3/gotk3/gtk"
)

type Menu struct {
	*vista.MenuPrincipal
}

func NuevoMenu(menu *vista.MenuPrincipal) *Menu {
	return &Menu{
		menu,
	}
}