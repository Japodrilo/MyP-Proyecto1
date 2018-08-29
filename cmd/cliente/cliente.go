package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"
	"github.com/Japodrilo/chatAlterno/pkg/modelo"
)

var wg sync.WaitGroup

/**
 * Lee del socket y escribe en la salida estándar.
 */
func Lee(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		str, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf(modelo.MSJ_DESCONECTA)
			wg.Done()
			return
		}
		fmt.Print(str)
	}
}

/**
 * Lee de la entrada estándar y escribe en el socket.
 */
func Escribe(conn net.Conn) {
	lector := bufio.NewReader(os.Stdin)
	escritor := bufio.NewWriter(conn)

	for {
		str, err := lector.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		_, err = escritor.WriteString(str)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = escritor.Flush()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

/**
 * Inicia un hilo de lectura y escritura que se conecta al servidor
 * mediante un socket.
 */
func main() {
	wg.Add(1)

	conn, err := net.Dial(modelo.CONN_TIPO, modelo.CONN_PUERTO)
	if err != nil {
		fmt.Println(err)
	}

	go Lee(conn)
	go Escribe(conn)

	wg.Wait()
}
