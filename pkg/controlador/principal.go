package controlador

import(
	"fmt"
	"log"
	"strings"

	"github.com/Japodrilo/MyP-Proyecto1/pkg/vista"
	"github.com/Japodrilo/MyP-Proyecto1/pkg/modelo"

	"github.com/gotk3/gotk3/gtk"
)

type Principal struct {
	cliente     *modelo.Cliente
	win			*gtk.Window
	cuaderno	*Cuaderno
	menu		*Menu
	lb			*gtk.ListBox
}

func NuevaPrincipal() *Principal {
	cliente := modelo.NuevoCliente()
	ventanaPrincipal := vista.SetupVentanaPrincipal()
	cuaderno := NuevoCuaderno(ventanaPrincipal.Nb)
	menu := NuevoMenu(ventanaPrincipal.Menubar)

	cuaderno.entradas["General"], cuaderno.textos["General"] = cuaderno.AddTab("General")
	cuaderno.entradas["General"].Connect("activate", MainEntryAction(cuaderno, cliente))

	ventanaPrincipal.Menubar.ConectarMI.Connect("activate", PopUpConectar(cliente))
	ventanaPrincipal.Menubar.DesconectarMI.Connect("activate", PopUpDesconectar(cliente))

	principal := &Principal{
		cliente:    cliente,
		win:		ventanaPrincipal.Win,
		cuaderno:	cuaderno,
		menu:		menu,
		lb:			ventanaPrincipal.Lb,
	}
	principal.Escucha()
	return principal
}

func (principal *Principal) AddUser(username string) {
	lbr := vista.SetupListBoxRow()
	entry := vista.SetupEntry()
	btn := SetupButtonUser(username, principal.cuaderno, entry, principal.cliente)
	lbr.Add(btn)
	principal.lb.Add(lbr)
	principal.win.ShowAll()
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
			go cliente.Lee(conn)
			go cliente.Escribe1(conn)
			go cliente.Escribe2(conn)
		})
	}
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
	principal.AddUser("Usuario 1")
	principal.AddUser("Usuario 2")
}

func parse(mensaje string) (string, []string) {
	prefijo := ""
	argumentos := strings.Fields(mensaje)
    switch {
	case strings.HasPrefix(mensaje, "...PUBLIC-"):
    	prefijo = "PUBLIC"
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
	}
}

func (principal *Principal) Escucha() {
	go func() {
		for {
			select {
			case mensaje := <- principal.cliente.Entrante:
				principal.manejaEntrada(mensaje)
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

func SetupButtonUser(username string, cuaderno *Cuaderno, entry *gtk.Entry, cliente *modelo.Cliente) *gtk.Button {
	btn, err := gtk.ButtonNewWithLabel(username)
	if err != nil {
		log.Fatal("Unable to create button:", err)
	}
	btn.Connect("clicked", func() {
		box := vista.SetupBox()
		scrwin := vista.SetupScrolledWindow()
		tv := vista.SetupTextView()
		tv.SetVExpand(true)
		box.Add(entry)
		box.Add(scrwin)
		
		nbTab := vista.SetupLabel(username)
		scrwin.Add(tv)
		cuaderno.nb.AppendPage(box, nbTab)
		cuaderno.nb.ShowAll()
		btn.SetSensitive(false)
		cuaderno.nb.SetCurrentPage(-1)
		cuaderno.textos[username] = vista.GetBufferTV(tv)
		cuaderno.entradas[username] = entry
		entry.Connect("activate", MainEntryAction(cuaderno, cliente))
		entry.GrabFocus()
	})
	return btn
}