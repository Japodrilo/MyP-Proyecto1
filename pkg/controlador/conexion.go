package controlador

import(
	"bufio"
	"fmt"
	"net"
	"os"
	"github.com/Japodrilo/MyP-Proyecto1/pkg/vista"

	"github.com/gotk3/gotk3/gtk"
)


func c(direccionE *gtk.Entry, puertoE *gtk.Entry, win *gtk.Window) func() {
	return func() {
		direccion := vista.GetTextEntry(direccionE)
		puerto := vista.GetTextEntry(puertoE)

		conn, err := net.Dial("tcp", direccion + ":" + puerto)
		if err != nil {
			vista.PopUpErrorConexion(direccion, puerto)
			return
		}

		win.Close()

		go func(conn net.Conn) {
			reader := bufio.NewReader(conn)
			for {
				str, err := reader.ReadString('\n')
				if err != nil {
					return
				}
			//fmt.Print(str)
			ManejaEntrada(str)
			}
		}(conn)

		go func(conn net.Conn) {
			lector := bufio.NewReader(os.Stdin)
			escritor := bufio.NewWriter(conn)

			for {
				str, err := lector.ReadString('\n')
				if err != nil {
					fmt.Println(err)
					return
				}

				_, err = escritor.WriteString(str)
				if err != nil {
					fmt.Println(err)
					return
				}

				err = escritor.Flush()
				if err != nil {
					fmt.Println(err)
					return
				}
			}
		}(conn)
	}
}

func PopUpConectar() func() {
	return vista.PopUpConectar(c)
}

func VentanaPrincipal() {
	principal := NuevaPrincipal()
	principal.AddUser("Usuario 1")
	principal.AddUser("Usuario 2")
}
