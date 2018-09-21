package vista

import(
	"github.com/gotk3/gotk3/gtk"
)

type MenuPrincipal struct{
	Menubar 	  *gtk.MenuBar
	ConectarMI    *gtk.MenuItem
	DesconectarMI *gtk.MenuItem
	ActivoMI      *gtk.MenuItem
	AlejadoMI     *gtk.MenuItem
	OcupadoMI     *gtk.MenuItem
	CerrarMI      *gtk.MenuItem
	CrearMI       *gtk.MenuItem
	InvitarMI     *gtk.MenuItem
	SalasMI	      *gtk.MenuItem
}

func SetupMenuPrincipal() *MenuPrincipal{
	menubar := SetupMenuBar()

	connMI := SetupMenuItemMnemonic("_Conexión")
	edosMI := SetupMenuItemMnemonic("_Estado")
	tabsMI := SetupMenuItemMnemonic("_Pestañas")
	salaMI := SetupMenuItemMnemonic("_Sala")
	

	connMenu := SetupMenu()
	edosMenu := SetupMenu()
	tabsMenu := SetupMenu()
	salaMenu := SetupMenu()
	


	connMI.SetSubmenu(connMenu)
	edosMI.SetSubmenu(edosMenu)
	tabsMI.SetSubmenu(tabsMenu)
	salaMI.SetSubmenu(salaMenu)
	

	conectarMI := SetupMenuItemLabel("Conectar")
	desconectarMI := SetupMenuItemLabel("Desconectar")

	activoMI := SetupMenuItemLabel("Activo")
	alejadoMI := SetupMenuItemLabel("Alejado")
	ocupadoMI := SetupMenuItemLabel("Ocupado")

	cerrarMI := SetupMenuItemLabel("Cerrar pestaña actual")

	crearMI := SetupMenuItemLabel("Crear")
	invitarMI := SetupMenuItemLabel("Invitar")
	salasMI := SetupMenuItemLabel("Mis salas")

	connMenu.Append(conectarMI)
	connMenu.Append(desconectarMI)

	edosMenu.Append(activoMI)
	edosMenu.Append(alejadoMI)
	edosMenu.Append(ocupadoMI)

	tabsMenu.Append(cerrarMI)

	salaMenu.Append(crearMI)
	salaMenu.Append(invitarMI)
	salaMenu.Append(salasMI)

	menubar.Append(connMI)
	menubar.Append(edosMI)
	menubar.Append(tabsMI)
	menubar.Append(salaMI)

	return &MenuPrincipal{
		Menubar:       menubar,
		ConectarMI:    conectarMI,
		DesconectarMI: desconectarMI,
		ActivoMI:      activoMI,
		AlejadoMI:     alejadoMI,
		OcupadoMI:     ocupadoMI,
		CerrarMI:      cerrarMI,
		CrearMI:       crearMI,
		InvitarMI:     invitarMI,
		SalasMI:       salasMI,
		
	}
}