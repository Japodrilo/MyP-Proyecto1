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
	CerrarMI        *gtk.MenuItem
}

func SetupMenuPrincipal() *MenuPrincipal{
	menubar := setup_menu_bar()

	connMI := setup_menu_item_mnemonic("_Conexión")
	salaMI := setup_menu_item_mnemonic("_Sala")
	tabsMI := setup_menu_item_mnemonic("_Pestañas")

	connMenu := setup_menu()
	salaMenu := setup_menu()
	tabsMenu := setup_menu()


	connMI.SetSubmenu(connMenu)
	salaMI.SetSubmenu(salaMenu)
	tabsMI.SetSubmenu(tabsMenu)

	conectarMI := setup_menu_item_label("Conectar")
	desconectarMI := setup_menu_item_label("Desconectar")

	invitarMI := setup_menu_item_label("Invitar")
	invitacionesMI := setup_menu_item_label("Invitaciones")

	cerrarMI := setup_menu_item_label("Cerrar pestaña actual")

	connMenu.Append(conectarMI)
	connMenu.Append(desconectarMI)

	salaMenu.Append(invitarMI)
	salaMenu.Append(invitacionesMI)

	tabsMenu.Append(cerrarMI)

	menubar.Append(connMI)
	menubar.Append(salaMI)
	menubar.Append(tabsMI)

	return &MenuPrincipal{
		Menubar: 		menubar,
		ConectarMI: 	conectarMI,
		DesconectarMI: 	desconectarMI,
		InvitarMI:		invitarMI,
		InvitacionesMI:	invitacionesMI,
		CerrarMI:       cerrarMI,
	}
}