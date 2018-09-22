package controlador

import(
	"github.com/Japodrilo/MyP-Proyecto1/pkg/vista"
)

/**
 * Clase que iba a ser intermediaria entre la vista
 * y el controlador de la ventana principal, pero al
 * final no me dió tiempo de descentralizar algunas
 * funciones que quería :/.
 */
type Menu struct {
	*vista.MenuPrincipal
}

/**
 * Su constructor trivial.
 */
func NuevoMenu(menu *vista.MenuPrincipal) *Menu {
	return &Menu{
		menu,
	}
}