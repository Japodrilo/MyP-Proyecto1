package controlador

import(
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Japodrilo/MyP-Proyecto1/pkg/vista"
	"github.com/Japodrilo/MyP-Proyecto1/pkg/modelo"

	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/gotk3/glib"
)

type Principal struct {
	cliente      *modelo.Cliente
	win			 *gtk.Window
	cuaderno	 *Cuaderno
	lb			 *gtk.ListBox
	menu         *vista.MenuPrincipal
	renglones    map[string]*gtk.ListBoxRow
	canalSalas   chan int
	salas        []string
	invitaciones []string
}

func NuevaPrincipal() *Principal {
	cliente := modelo.NuevoCliente()
	ventanaPrincipal := vista.SetupVentanaPrincipal()
	cuaderno := NuevoCuaderno(ventanaPrincipal.Nb)
	renglones := make (map[string]*gtk.ListBoxRow)
	canalSalas := make(chan int)
	salas := make([]string, 0)

	cuaderno.entradas["General"], cuaderno.textos["General"] = cuaderno.AddTab("General")
	cuaderno.entradas["General"].Connect("activate", MainEntryAction(cuaderno, cliente))
	cuaderno.entradas["General"].SetSensitive(false)

	ventanaPrincipal.Menubar.ConectarMI.Connect("activate", PopUpConectar(cliente))
	ventanaPrincipal.Menubar.DesconectarMI.Connect("activate", PopUpDesconectar(cliente))
	ventanaPrincipal.Menubar.DesconectarMI.SetSensitive(false)

	ventanaPrincipal.Menubar.ActivoMI.Connect("activate", func () {
		glib.IdleAdd(ventanaPrincipal.Win.SetTitle, "Chat")
		cliente.Saliente <- "STATUS ACTIVE\n"
		})
	ventanaPrincipal.Menubar.AlejadoMI.Connect("activate", func () {
		glib.IdleAdd(ventanaPrincipal.Win.SetTitle, "Chat (ALEJADO)")
		cliente.Saliente <- "STATUS AWAY\n"
		})
	ventanaPrincipal.Menubar.OcupadoMI.Connect("activate", func () {
		glib.IdleAdd(ventanaPrincipal.Win.SetTitle, "Chat (OCUPADO)")
		cliente.Saliente <- "STATUS BUSY\n"
		})
	ventanaPrincipal.Menubar.ActivoMI.SetSensitive(false)
	ventanaPrincipal.Menubar.AlejadoMI.SetSensitive(false)
	ventanaPrincipal.Menubar.OcupadoMI.SetSensitive(false)

	principal := &Principal{
		cliente:    cliente,
		win:		ventanaPrincipal.Win,
		cuaderno:	cuaderno,
		lb:			ventanaPrincipal.Lb,
		menu:       ventanaPrincipal.Menubar,
		renglones:  renglones,
		canalSalas: canalSalas,
		salas:      salas,
	}
	ventanaPrincipal.Menubar.CrearMI.Connect("activate", func () {principal. PopUpCrearSala()})
	ventanaPrincipal.Menubar.CrearMI.SetSensitive(false)
	ventanaPrincipal.Menubar.InvitarMI.Connect("activate", func () {principal. PopUpInvitar()})
	ventanaPrincipal.Menubar.InvitarMI.SetSensitive(false)
	ventanaPrincipal.Menubar.SalasMI.Connect("activate", func () {principal. PopUpMisSalas()})
	ventanaPrincipal.Menubar.SalasMI.SetSensitive(false)
	ventanaPrincipal.Menubar.CerrarMI.Connect("activate", func () {principal.eliminaPestana()})
	principal.Escucha()
	principal.win.ShowAll()
	return principal
}

func (principal *Principal) Escucha() {
	go func() {
		for {
			select {
			case mensaje := <- principal.cliente.Entrante:
				principal.manejaEntrada(mensaje)

			case estado := <- principal.cliente.Activo:
				if estado {
					principal.EncenderTodo()
				} else {
					principal.ApagarTodo()
				}
			}
		}
	}()
}

func (principal *Principal) ApagarTodo() {
	principal.win.SetTitle("Chat")
	i := 0
	boton := principal.lb.GetRowAtIndex(i)
	for boton != nil {
		glib.IdleAdd(principal.lb.Remove, boton)
		i++
		boton = principal.lb.GetRowAtIndex(i)
		principal.cuaderno.botones = make(map [string]*gtk.Button)
	}
	for _, entrada := range principal.cuaderno.entradas {
		entrada.SetSensitive(false)
	}
	principal.menu.ConectarMI.SetSensitive(true)
	principal.menu.CrearMI.SetSensitive(false)
	principal.menu.InvitarMI.SetSensitive(false)
	principal.menu.SalasMI.SetSensitive(false)
	principal.menu.ActivoMI.SetSensitive(false)
	principal.menu.AlejadoMI.SetSensitive(false)
	principal.menu.OcupadoMI.SetSensitive(false)
	principal.menu.DesconectarMI.SetSensitive(false)
}

func (principal *Principal) EncenderTodo() {
	for _, entrada := range principal.cuaderno.entradas {
		entrada.SetSensitive(true)
	}
	principal.menu.ConectarMI.SetSensitive(false)
	principal.menu.CrearMI.SetSensitive(true)
	principal.menu.InvitarMI.SetSensitive(true)
	principal.menu.SalasMI.SetSensitive(true)
	principal.menu.ActivoMI.SetSensitive(true)
	principal.menu.AlejadoMI.SetSensitive(true)
	principal.menu.OcupadoMI.SetSensitive(true)
	principal.menu.DesconectarMI.SetSensitive(true)
}


func (principal *Principal) AddUserButton(username string) {
	if principal.cuaderno.botones[username] != nil || principal.cliente.Nombre == username {
		return
	}
	lbr := vista.SetupListBoxRow()
	btn := principal.SetupButtonUser(username)
	lbr.Add(btn)
	principal.cuaderno.botones[username] = btn
	principal.renglones[username] = lbr
	principal.lb.Add(lbr)
	principal.win.ShowAll()
}

func (principal *Principal) eliminaPestana() {
	if principal.cuaderno.nb.GetCurrentPage() == 0 {
			return
	}
	pn := principal.cuaderno.nb.GetCurrentPage()
	page, _ := principal.cuaderno.nb.GetNthPage(pn)
	user, _ := principal.cuaderno.nb.GetTabLabelText(page)
	principal.cuaderno.nb.RemovePage(pn)
	delete(principal.cuaderno.tabs, user)
	delete(principal.cuaderno.textos, user)
	if principal.cuaderno.botones[user] != nil {
		principal.cuaderno.botones[user].SetSensitive(true)
	}
}

func actualizaUsuarios(cliente *modelo.Cliente) {
	ticker := time.NewTicker(5 * time.Second)
	cliente.Saliente <- "USERS\n"
	go func() {for {
		select {
		case activo := <- cliente.Activo:
			if !activo {
				return
			}
		case <- ticker.C:
			cliente.Saliente <- "USERS\n"
		}
	} }()
}

func PopUpConectar(cliente *modelo.Cliente) func() {
	return func() {
		emergente := vista.NuevaConectar()
		emergente.ConectarB.Connect("clicked", func() {
			direccion := vista.GetTextEntry(emergente.DireccionE)
			puerto := vista.GetTextEntry(emergente.PuertoE)
			conn := cliente.Conecta(direccion, puerto)
			if conn == nil {
				vista.PopUpErrorConexion(direccion, puerto)
				return
			} 
			emergente.Win.Close()
			PopUpIdentificarse(cliente)

			go cliente.Lee(conn)
			go cliente.Escribe1(conn)
			go cliente.Escribe2(conn)
		})
	}
}

func (principal *Principal) PopUpCrearSala() {
	emergente := vista.NuevaCrear()
	emergente.CrearB.Connect("clicked", func() {
		nombre := vista.GetTextEntry(emergente.Nombre)
		principal.cliente.Saliente <- "CREATEROOM " + nombre + "\n"
		time.Sleep(400 * time.Millisecond)
		if nombre == "" {
			emergente.Win.Close()
			return	
		}
		select{
			case i := <- principal.canalSalas:
				switch i {
				case 0: 
					glib.IdleAdd(principal.AddTab, "*S*-" + nombre)
					principal.salas = append(principal.salas, nombre)
					emergente.Win.Close()
					return
				case 1: 
					vista.NombreSalaOcupado()
				}
		}
	})
}

func (principal *Principal) PopUpMisSalas() {
	emergente := vista.NuevaMisSalas()
	for _, sala := range principal.salas {
			emergente.SalaCBT.AppendText(sala)
		}
	emergente.AbrirB.Connect("clicked", func() {
		sala := emergente.SalaCBT.GetActiveText()
		pagina := principal.cuaderno.tabs["*S*-" + sala]
		if pagina > 0 {
			principal.cuaderno.nb.SetCurrentPage(pagina)
			principal.cuaderno.entradas["*S*-" + sala].GrabFocus()
			emergente.Win.Close()
			return
		}
		principal.AddTab("*S*-" + sala)
		emergente.Win.Close()
	})
}

func PopUpIdentificarse(cliente *modelo.Cliente) {
	emergente := vista.NuevaIdentificar()
	emergente.IdentificarB.Connect("clicked", func() {
		nombre := vista.GetTextEntry(emergente.Nombre)
		cliente.Saliente <- "IDENTIFY " + nombre + "\n"
		time.Sleep(200 * time.Millisecond)
		if !cliente.Identificado {
			vista.NombreOcupado()
		} else {
			cliente.Nombre = nombre
			emergente.Win.Close()
			cliente.Activo <- true
			actualizaUsuarios(cliente)
		}
	})
}

func (principal *Principal) PopUpInvitar() {
	emergente := vista.NuevaInvitar()
	for usuario, _ := range principal.renglones {
		emergente.NombreCBT.AppendText(usuario)
	}
	for _, sala := range principal.salas {
		emergente.SalaCBT.AppendText(sala)
	}
	emergente.InvitarB.Connect("clicked", func() {
		sala := emergente.SalaCBT.GetActiveText()
		nombre := emergente.NombreCBT.GetActiveText()
		principal.cliente.Saliente <- "INVITE " + sala + " " + nombre + "\n"
		time.Sleep(400 * time.Millisecond)
		select{
			case i := <- principal.canalSalas:
				switch i {
				case 0:
					emergente.Win.Close()
					return
				case 1:
					vista.UsuarioInexistente()
				case 2:
					vista.SalaInexistente()
				case 3:
					vista.NoTePertenece()
				}
		}
	})
}

func PopUpDesconectar(cliente *modelo.Cliente) func() {
	return func() {
		emergente := vista.NuevaDesconectar()
		emergente.DesconectarB.Connect("clicked", func() {
			cliente.Desconecta()
			emergente.Win.Close()
		})
	}
}

func VentanaPrincipal() {
	principal := NuevaPrincipal()
	principal.win.ShowAll()
}

func contiene(rebanada []string, cadena string) bool {
	for _, entrada := range rebanada {
		if entrada == cadena {
			return true
		}
	}
	return false
}


func (principal *Principal) parse(mensaje string) (string, []string) {
	prefijo := ""
	argumentos := strings.Fields(mensaje)
    switch {
    case strings.HasPrefix(mensaje, "...INVITATION TO JOIN"):
    	prefijo = "INVITATION_JOIN"
    case (strings.HasPrefix(mensaje, "...USER") && strings.HasSuffix(mensaje, " NOT FOUND\n")):
    	prefijo = "INVITATION_NOT_OK"
    case mensaje == "...ROOM NOT EXIST\n":
    	prefijo = "NO_ROOM"
    case mensaje == "...YOU ARE NOT THE OWNER OF THE ROOM\n":
    	prefijo = "NOT_OWNER"
    case strings.HasPrefix(mensaje, "...INVITATION SENT TO"):
    	prefijo = "INVITATION_OK"
    case len(argumentos) == 2 && (argumentos[1] == "ACTIVE" || argumentos[1] == "AWAY" || argumentos[1] == "BUSY"):
    	prefijo = "STATUS"
    case mensaje == "...SUCCESSFUL IDENTIFICATION\n":
    	prefijo = "ID"
    case mensaje == "...ROOM CREATED\n":
    	prefijo = "ROOM_OK"
    case mensaje == "...ROOM NAME ALREADY IN USE\n":
    	prefijo = "ROOM_NOT_OK"
	case strings.HasPrefix(mensaje, "...PUBLIC-"):
    	prefijo = "PUBLIC"
    case strings.HasPrefix(mensaje, "..."):
    	for _, sala := range principal.salas {
    		if strings.HasPrefix(mensaje, "..." + sala + "-") {
    			prefijo = "SALA"
    		}
    	}
    case strings.HasSuffix(argumentos[0], ":"):
    	prefijo = "DIRECTO"
    default:
    	prefijo = "USERS"
	}
  	return prefijo, argumentos
}

func (principal *Principal) manejaEntrada(mensaje string) {
	prefijo, argumentos := principal.parse(mensaje)
	switch prefijo {
	case "":
		return
	case "PUBLIC":
		nombre := strings.TrimPrefix(argumentos[0], "...PUBLIC-")
		mensaje := strings.Join(argumentos[1:], " ")
		mensaje = nombre + " " + mensaje + "\n"
		principal.cuaderno.textos["General"].InsertAtCursor(mensaje)
	case "DIRECTO":
		nombre := strings.TrimSuffix(argumentos[0], ":")
		mensaje := strings.Join(argumentos[1:], " ")
		mensaje = nombre + ": " + mensaje + "\n"
		if principal.cuaderno.textos[nombre] == nil {
			glib.IdleAdd(principal.AddTab, nombre)
			time.Sleep(100 * time.Millisecond)
		}
		principal.cuaderno.textos[nombre].InsertAtCursor(mensaje)
	case "SALA":
		sala := ""
		for _, s := range principal.salas {
    		if strings.HasPrefix(mensaje, "..." + s + "-") {
    			sala = s
    		}
    	}
    	nombre := strings.TrimPrefix(argumentos[0], "..." + sala + "-")
    	mensaje := strings.Join(argumentos[1:], " ")
    	mensaje = fmt.Sprintf("%v %v\n", nombre, mensaje)
    	if principal.cuaderno.textos["*S*-" + sala] == nil {
    		glib.IdleAdd(principal.AddTab, "*S*-" + sala)
    		time.Sleep(100 * time.Millisecond)
    	}
    	principal.cuaderno.textos["*S*-" + sala].InsertAtCursor(mensaje)
	case "USERS":
		for _, usuario := range argumentos {
			glib.IdleAdd(principal.AddUserButton, usuario)
		}
		for usuario, _ := range principal.cuaderno.botones {
			if !contiene(argumentos, usuario) {
				glib.IdleAdd(principal.lb.Remove, principal.renglones[usuario])
				delete(principal.renglones, usuario)
				delete(principal.cuaderno.botones, usuario)
			}
		}
	case "ID":
		principal.cliente.Identificado = true
		time.Sleep(400 * time.Millisecond)
	case "ROOM_OK":
		principal.canalSalas <- 0
	case "ROOM_NOT_OK":
		principal.canalSalas <- 1
	case "STATUS":
		notificacion := "\nEl usuario %s está %s\n\n"
		switch argumentos[1] {
		case "ACTIVE":
			notificacion := fmt.Sprintf(notificacion, argumentos[0], "ACTIVO")
			principal.cuaderno.textos["General"].InsertAtCursor(notificacion)
		case "AWAY":
			notificacion := fmt.Sprintf(notificacion, argumentos[0], "ALEJADO")
			principal.cuaderno.textos["General"].InsertAtCursor(notificacion)
		case "BUSY":
			notificacion := fmt.Sprintf(notificacion, argumentos[0], "OCUPADO")
			principal.cuaderno.textos["General"].InsertAtCursor(notificacion)
		}
	case "INVITATION_OK":
		principal.canalSalas <- 0
	case "INVITATION_NOT_OK":
		principal.canalSalas <- 1
	case "NO_ROOM":
		principal.canalSalas <- 2
	case "NOT_OWNER":
		principal.canalSalas <- 3
	case "INVITATION_JOIN":
		principal.cliente.Saliente <- "JOINROOM " + argumentos[3] + "\n"
		principal.salas = append(principal.salas, argumentos[3])
		glib.IdleAdd(principal.AddTab, "*S*-" + argumentos[3])

	}
}

func(principal *Principal) AddTab(name string) {
	box := vista.SetupBox()
	entry := vista.SetupEntry()
	scrwin := vista.SetupScrolledWindow()
	tv := vista.SetupTextView()
	nbTab := vista.SetupLabel(name)
	principal.cuaderno.textos[name] = vista.GetBufferTV(tv)
	principal.cuaderno.entradas[name] = entry

	entry.Connect("activate", MainEntryAction(principal.cuaderno, principal.cliente))
	
	tv.SetVExpand(true)

	scrwin.Add(tv)
	box.Add(entry)
	box.Add(scrwin)

	principal.cuaderno.nb.SetHExpand(true)
	principal.cuaderno.nb.SetVExpand(true)

	principal.cuaderno.nb.AppendPage(box, nbTab)
	principal.win.ShowAll()
	principal.cuaderno.nb.SetCurrentPage(-1)
	principal.cuaderno.tabs[name] = principal.cuaderno.nb.GetCurrentPage()
}

func MainEntryAction(cuaderno *Cuaderno, cliente *modelo.Cliente) func() {
	return func() {
		pn := cuaderno.nb.GetCurrentPage()
		page, _ := cuaderno.nb.GetNthPage(pn)
		user, _ := cuaderno.nb.GetTabLabelText(page)
		entrada := cuaderno.entradas[user]
		text := vista.GetTextEntryClean(entrada)
		if text == "\n" {
			return
		}
		q := cuaderno.textos[user]
		fmt.Print(text)
		q.InsertAtCursor(cliente.Nombre + ": " + text)
		switch {
		case user == "General":
			cliente.Saliente <- "PUBLICMESSAGE " + text
		case strings.HasPrefix(user, "*S*-"):
			user = strings.TrimPrefix(user, "*S*-")
			cliente.Saliente <- fmt.Sprintf("ROOMESSAGE %v %v", user, text)
		default:
			cliente.Saliente <- "MESSAGE " + user + " " + text
		}
	}
}

func (principal *Principal) SetupButtonUser(username string) *gtk.Button {
	btn, err := gtk.ButtonNewWithLabel(username)
	if err != nil {
		log.Fatal("Unable to create button:", err)
	}

	btn.Connect("clicked", func() {
		if principal.cuaderno.tabs[username] > 0 {
			principal.cuaderno.nb.SetCurrentPage(principal.cuaderno.tabs[username])
			btn.SetSensitive(false)
			return
		}
		box := vista.SetupBox()
		scrwin := vista.SetupScrolledWindow()
		tv := vista.SetupTextView()
		entry := vista.SetupEntry()
		tv.SetVExpand(true)
		box.Add(entry)
		box.Add(scrwin)
		
		nbTab := vista.SetupLabel(username)
		scrwin.Add(tv)
		principal.cuaderno.nb.AppendPage(box, nbTab)
		principal.cuaderno.nb.ShowAll()
		btn.SetSensitive(false)
		principal.cuaderno.nb.SetCurrentPage(-1)
		principal.cuaderno.tabs[username] = principal.cuaderno.nb.GetCurrentPage()
		principal.cuaderno.textos[username] = vista.GetBufferTV(tv)
		principal.cuaderno.entradas[username] = entry
		entry.Connect("activate", MainEntryAction(principal.cuaderno, principal.cliente))
		entry.GrabFocus()
	})
	return btn
}