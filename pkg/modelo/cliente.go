package modelo

import(
	"bufio"
	"fmt"
	//"io"
	//"log"
	"os"
	"net"
	"time"
)

type Cliente struct {
	Nombre       string
	conn 		 net.Conn
	Conectado	 bool
	Identificado bool
	Activo       chan bool
	Entrante	 chan string
	Saliente	 chan string
	lector		 *bufio.Reader
	escritor	 *bufio.Writer
}

func NuevoCliente() *Cliente {
	var conn net.Conn
	var lector *bufio.Reader
	var escritor *bufio.Writer
	return &Cliente{
		Nombre: "YO",
		conn: conn,
		Conectado: false,
		Identificado: false,
		Activo: make(chan bool),
		Entrante: make(chan string),
		Saliente: make(chan string),
		lector: lector,
		escritor: escritor,
	}
}

func (cliente *Cliente) Conecta(direccion, puerto string) net.Conn {
	if cliente.Conectado {
		return cliente.conn
	}
	conn, err := net.DialTimeout("tcp", direccion + ":" + puerto, 5 * time.Second)
	if err != nil {
		fmt.Println("No se pudo establecer la conexión.")
		return nil
	}
	cliente.conn = conn
	cliente.Conectado = true
	return conn
}

/**
 * Lee del socket y escribe en la salida estándar.
 */
func (cliente *Cliente) Lee(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		str, err := reader.ReadString('\n')
		if err != nil {
			cliente.Activo <- false
			cliente.Conectado = false
			cliente.Identificado = false
			return
		}
		fmt.Print(str)
		cliente.Entrante <- str
	}
}

/**
 * Lee de la entrada estándar y escribe en el socket.
 */
func (cliente *Cliente) Escribe1(conn net.Conn) {
	lector := bufio.NewReader(os.Stdin)
	escritor := bufio.NewWriter(conn)

	for {
		str, err := lector.ReadString('\n')
		if err != nil {
			return
		}

		_, err = escritor.WriteString(str)
		if err != nil {
			return
		}

		err = escritor.Flush()
		if err != nil {
			return
		}
	}
}

/**
 * Lee de la entrada estándar y escribe en el socket.
 */
func (cliente *Cliente) Escribe2(conn net.Conn) {
	escritor := bufio.NewWriter(conn)

	for {
		select {
		case str := <- cliente.Saliente:
			_, err := escritor.WriteString(str)
			if err != nil {
				cliente.Activo <- false
				cliente.Conectado = false
				cliente.Identificado = false
				return
			}
			err = escritor.Flush()
			if err != nil {
				cliente.Activo <- false
				cliente.Conectado = false
				cliente.Identificado = false
				return
			}
		default:
			continue
		}
	}
}

func (cliente *Cliente) Desconecta() {
	if cliente.Conectado {
		cliente.Saliente <- "DISCONNECT\n"
		cliente.Conectado = false
		cliente.Identificado = false
		cliente.Activo <- false
	}
}