package vista

import(
	"github.com/gotk3/gotk3/gtk"
)

type MenuPrincipal struct{
	Menubar 		*gtk.MenuBar
	ConectarMI 		*gtk.MenuItem
	DesconectarMI	*gtk.MenuItem
	InvitarMI		*gtk.MenuItem
	InvitacionesMI	*gtk.MenuItem
}

func SetupMenuPrincipal() *MenuPrincipal{
	menubar := setup_menu_bar()

	connMI := setup_menu_item_mnemonic("_Conexión")
	salaMI := setup_menu_item_mnemonic("_Sala")

	connMenu := setup_menu()
	salaMenu := setup_menu()

	connMI.SetSubmenu(connMenu)
	salaMI.SetSubmenu(salaMenu)

	conectarMI := setup_menu_item_label("Conectar")
	desconectarMI := setup_menu_item_label("Desconectar")

	invitarMI := setup_menu_item_label("Invitar")
	invitacionesMI := setup_menu_item_label("Invitaciones")

	connMenu.Append(conectarMI)
	connMenu.Append(desconectarMI)

	salaMenu.Append(invitarMI)
	salaMenu.Append(invitacionesMI)

	menubar.Append(connMI)
	menubar.Append(salaMI)

	return &MenuPrincipal{
		Menubar: 		menubar,
		ConectarMI: 	conectarMI,
		DesconectarMI: 	desconectarMI,
		InvitarMI:		invitarMI,
		InvitacionesMI:	invitacionesMI,
	}
}