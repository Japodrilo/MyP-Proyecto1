package main

import (
	"github.com/Japodrilo/MyP-Proyecto1/pkg/controlador"

	"github.com/gotk3/gotk3/gtk"
)


/**
 * Inicia un hilo de lectura y escritura que se conecta al servidor
 * mediante un socket.
 */
func main() {

	gtk.Init(nil)

	controlador.VentanaPrincipal()

	gtk.Main()

}
