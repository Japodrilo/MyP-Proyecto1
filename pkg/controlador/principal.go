package controlador

import(
	"fmt"
	"log"
	//"regexp"
	"strings"
	"time"

	"github.com/Japodrilo/MyP-Proyecto1/pkg/vista"
	"github.com/Japodrilo/MyP-Proyecto1/pkg/modelo"

	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/gotk3/glib"
)

type Principal struct {
	cliente     *modelo.Cliente
	win			*gtk.Window
	cuaderno	*Cuaderno
	menu		*Menu
	lb			*gtk.ListBox
	renglones   map[string]*gtk.ListBoxRow
}

func NuevaPrincipal() *Principal {
	cliente := modelo.NuevoCliente()
	ventanaPrincipal := vista.SetupVentanaPrincipal()
	cuaderno := NuevoCuaderno(ventanaPrincipal.Nb)
	menu := NuevoMenu(ventanaPrincipal.Menubar)
	renglones := make (map[string]*gtk.ListBoxRow)

	cuaderno.entradas["General"], cuaderno.textos["General"] = cuaderno.AddTab("General")
	cuaderno.entradas["General"].Connect("activate", MainEntryAction(cuaderno, cliente))
	cuaderno.entradas["General"].SetSensitive(false)

	ventanaPrincipal.Menubar.ConectarMI.Connect("activate", PopUpConectar(cliente))
	ventanaPrincipal.Menubar.DesconectarMI.Connect("activate", PopUpDesconectar(cliente))

	principal := &Principal{
		cliente:    cliente,
		win:		ventanaPrincipal.Win,
		cuaderno:	cuaderno,
		menu:		menu,
		lb:			ventanaPrincipal.Lb,
		renglones: renglones,
	}
	ventanaPrincipal.Menubar.CerrarMI.Connect("activate", func () {principal.eliminaPestana()})
	principal.Escucha()
	principal.win.ShowAll()
	return principal
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


func parse(mensaje string) (string, []string) {
	prefijo := ""
	argumentos := strings.Fields(mensaje)
    switch {
    case mensaje == "...SUCCESSFUL IDENTIFICATION\n":
    	prefijo = "ID"
	case strings.HasPrefix(mensaje, "...PUBLIC-"):
    	prefijo = "PUBLIC"
    case strings.HasSuffix(argumentos[0], ":"):
    	prefijo = "DIRECTO"
    case strings.HasPrefix(mensaje, "..."):
    default:
    	prefijo = "USERS"
	}
  	return prefijo, argumentos
}

func (principal *Principal) manejaEntrada(mensaje string) {
	prefijo, argumentos := parse(mensaje)
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

func (principal *Principal) Escucha() {
	go func() {
		for {
			select {
			case mensaje := <- principal.cliente.Entrante:
				principal.manejaEntrada(mensaje)

			case estado := <- principal.cliente.Activo:
				if estado {
					for _, entrada := range principal.cuaderno.entradas {
						entrada.SetSensitive(true)
					}
				} else {
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
				}
			}
		}
	}()
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
		switch user {
		case "General":
			cliente.Saliente <- "PUBLICMESSAGE " + text
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