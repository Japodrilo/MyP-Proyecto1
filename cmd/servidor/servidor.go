package main

import (
	"log"
	"net"
	"os"
	"github.com/Japodrilo/chatAlterno/pkg/modelo"
)

/**
 * Crea un vestíbulo que escucha las conexiones de los clientes
 * y las conecta al vestíbulo.
 */
func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	vestibulo := modelo.NuevoVestibulo()

	escucha, err := net.Listen(modelo.CONN_TIPO, modelo.CONN_PUERTO)
	if err != nil {
		log.Println("Error: ", err)
		os.Exit(1)
	}
	defer escucha.Close()
	log.Println("Escuchando en el puerto " + modelo.CONN_PUERTO)

	for {
		conn, err := escucha.Accept()
		if err != nil {
			log.Println("Error: ", err)
			continue
		}
		vestibulo.Entrar(modelo.NuevaConexion(conn))
	}
}
