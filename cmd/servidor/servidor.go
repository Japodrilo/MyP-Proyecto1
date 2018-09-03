package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"github.com/Japodrilo/MyP-Proyecto1/pkg/modelo"
)


/**
 * Indica al usuario cómo debe usarse el programa.
 */
func uso() {
	fmt.Println("Uso: ./servidor puerto")
	os.Exit(0)
}

/**
 * Crea un vestíbulo que escucha las conexiones de los clientes
 * y las conecta al vestíbulo.
 */
func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	args := os.Args[1:]
	if len(args) != 1 {
		uso()
	}

	vestibulo := modelo.NuevoVestibulo()

	escucha, err := net.Listen(modelo.CONN_TIPO, ":" + args[0])
	if err != nil {
		fmt.Printf("Error al intentar escuchar en el puerto %v, servidor terminado.\n", args[0])
		os.Exit(1)
	}
	defer escucha.Close()
	log.Println("Escuchando en el puerto :" + args[0])

	for {
		conn, err := escucha.Accept()
		if err != nil {
			log.Println("Error: ", err)
			continue
		}
		vestibulo.Entrar(modelo.NuevaConexion(conn))
	}
}
