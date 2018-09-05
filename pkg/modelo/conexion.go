package modelo

import (
  "bufio"
  "fmt"
  "net"
  "log"
  "time"
  "strings"
)

var contadorConexiones = 0

/**
 * Una Conexion contiene la información y las acciones que realiza un
 * cliente al conectarse al servidor.   Tiene un nombre, las salas en las
 * que se encuentra, un canal de mensajes entrantes y salientes, y un
 * lector y un escritor para comunicarse.
 */
type Conexion struct {
	nombre   string
  serial   int
	salas    map[string]*Sala
	entrante chan *Mensaje
	saliente chan string
	conn     net.Conn
	lector   *bufio.Reader
	escritor *bufio.Writer
}

/**
 * Regresa una nueva Conexion con la conexión de red dada, inicia
 * un lector y un escritor que reciben y envían información por el socket.
 */
func NuevaConexion(conn net.Conn) *Conexion {
	escritor := bufio.NewWriter(conn)
	lector := bufio.NewReader(conn)

	conexion := &Conexion{
		nombre:   fmt.Sprintf(CLIENTE_NOMBRE, contadorConexiones),
    serial:   contadorConexiones,
		salas:    make(map[string]*Sala),
		entrante: make(chan *Mensaje),
		saliente: make(chan string),
		conn:     conn,
		lector:   lector,
		escritor: escritor,
	}
  contadorConexiones++
	conexion.Escucha()
	return conexion
}

/**
 * Asigna nombre a la conexión.
 */
func (conexion *Conexion) SetNombre(nombre string) {
  conexion.nombre = nombre
}

/**
 * Regresa el serial de la conexión.
 */
func (conexion *Conexion) Serial() int {
  return conexion.serial
}


/**
 * Lee las cadenas en el socket de la Conexion, les da formato de
 * Mensaje y los envía al canal entrante de la Conexion.
 */
func (conexion *Conexion) Lee() {
  log.SetFlags(log.LstdFlags | log.Lshortfile)
	for {
    receptaculo := make([]byte, 1024)
    n, err := conexion.lector.Read(receptaculo)
		//str, err := conexion.lector.ReadString('\n')
		if err != nil {
			log.Println(err)
			break
		}
    str := string(receptaculo[:n])
		message := NuevoMensaje(time.Now(), conexion, strings.TrimSuffix(str, "\n"))
		conexion.entrante <- message
	}
	close(conexion.entrante)
	log.Printf("Se cerró el hilo de lectura del canal entrante de %v\n", conexion.nombre)
}

/**
 * Lee los mensajes del canal saliente y los escribe en el socket.
 */
func (conexion *Conexion) Escribe() {
	for str := range conexion.saliente {
		_, err := conexion.escritor.WriteString(str)
		if err != nil {
			log.Println(err)
			break
		}
		err = conexion.escritor.Flush()
		if err != nil {
			log.Println(err)
			break
		}
	}
	log.Printf("Se cerró el hilo de escritura de %v\n", conexion.nombre)
}

/**
 * Inicia dos goroutines, la primera lee del canal de salida del cliente y
 * escribe a su socket, y la segunda lee del socket del cliente, y escribe
 * en su canal de entrada.
 */
func (conexion *Conexion) Escucha() {
  go conexion.Escribe()
  go conexion.Lee()
}

/**
 * Cierra la conexión del cliente.
 */
func (conexion *Conexion) Terminar() {
  conexion.conn.Close()
}
