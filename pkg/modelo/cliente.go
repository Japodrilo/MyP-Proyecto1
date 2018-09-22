package modelo

import(
	"bufio"
	"fmt"
	"net"
	"time"
)

/**
 * Estructura que modela a un cliente.  Contiene una cadena,
 * con el nombre del cliente, un socket, un canal para comunicar
 * su estado al controlador, dos canales para comunicar mensajes
 * con el controlador, un lector y un escritor.   Además cuenta
 * con dos booleanos que juegan el papel de banderas para señalar
 * estados.
 */
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

/**
 * Constructor básico que inicializa un cliente.
 */
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

/**
 * Método que entabla una conexión remota, recibe como parámetros
 * dos cadenas, indicando la dirección IP y el puerto para entablar
 * la comunicación.   Regresa un socket.
 */
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
 * Método que lee del socket y escribe en el canal Saliente (comunicación
 * con el controlador).   Está pensado para utilizarse en una goroutine.
 */
func (cliente *Cliente) Lee(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		str, err := reader.ReadString('\n')
		if err != nil {
			cliente.Conectado = false
			cliente.Identificado = false
			cliente.Activo <- false
			return
		}
		cliente.Entrante <- str
	}
}


/**
 * Método qu lee del canal Saliente (comunicación con el controlador) y
 * escribe en el socket.   Está pensado para utilizarse en una goroutine.
 */
func (cliente *Cliente) Escribe(conn net.Conn) {
	escritor := bufio.NewWriter(conn)
	for {
		select {
		case str := <- cliente.Saliente:
			_, err := escritor.WriteString(str)
			if err != nil {
				return
			}
			err = escritor.Flush()
			if err != nil {
				return
			}
		default:
			continue
		}
	}
}

/**
 * Método que desconecta al cliente y avisa al controlador
 * que ya no está activo mediante el canal Activo.
 */
func (cliente *Cliente) Desconecta() {
	if cliente.Conectado {
		cliente.Saliente <- "DISCONNECT\n"
		cliente.Conectado = false
		cliente.Identificado = false
		cliente.Activo <- false
	}
}