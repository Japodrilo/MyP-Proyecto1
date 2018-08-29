package modelo

import (
  "log"
  "strings"
  "fmt"
)

/**
 * Un Vestibulo administra a los clientes conectados, las salas existentes,
 * y distribuye los mensajes.
 */
type Vestibulo struct {
	conexiones []*Conexion
	salas      map[string]*Sala
	entrante   chan *Mensaje
	join       chan *Conexion
	leave      chan *Conexion
}

/**
 * Crea un Vestibulo que empieza a escuchar en sus canales.
 */
func NuevoVestibulo() *Vestibulo {
	vestibulo := &Vestibulo{
		conexiones: make([]*Conexion, 0),
		salas:      make(map[string]*Sala),
		entrante:   make(chan *Mensaje),
		join:       make(chan *Conexion),
		leave:      make(chan *Conexion),
	}
  vestibulo.salas["recepcion"] = NuevaSala("recepcion")
	vestibulo.Escucha()
	return vestibulo
}

/**
 * Inicia una gorutine nueva que escucha en los canales del Vestibulo.
 */
func (vestibulo *Vestibulo) Escucha() {
	go func() {
		for {
			select {
			case mensaje := <-vestibulo.entrante:
				vestibulo.Parse(mensaje)
			case conexion := <-vestibulo.join:
				vestibulo.Entrar(conexion)
			case conexion := <-vestibulo.leave:
				vestibulo.Elimina(conexion)
			}
		}
	}()
}

/**
 * Se encarga de los clientes que se conectan al vestíbulo.
 */
func (vestibulo *Vestibulo) Entrar(conexion *Conexion) {
	vestibulo.conexiones = append(vestibulo.conexiones, conexion)
	conexion.saliente <- MSJ_CONECTA
  vestibulo.salas["recepcion"].Agrega(conexion)
	go func() {
		for mensaje := range conexion.entrante {
			vestibulo.entrante <- mensaje
		}
		vestibulo.leave <- conexion
	}()
}

/**
 * Se encarga de los clientes que se desconectan del vestíbulo.
 */
func (vestibulo *Vestibulo) Elimina(conexion *Conexion) {
  for _, sala := range conexion.salas {
    sala.Elimina(conexion)
  }
	for i, otherConexion := range vestibulo.conexiones {
		if conexion == otherConexion {
			vestibulo.conexiones = append(vestibulo.conexiones[:i], vestibulo.conexiones[i+1:]...)
			break
		}
	}
	close(conexion.saliente)
	log.Printf("Se cerró el canal saliente de %v\n", conexion.nombre)
}

/**
 * Se encarga de los mensajes que llegan al vestíbulo.   Si el mensaje
 * contiene un comando, éste es ejecutado por el vestíbulo.   Si no, el
 * mensaje se envía a las salas en las que está el remitente.
 */
func (vestibulo *Vestibulo) Parse(mensaje *Mensaje) {
	switch {
	default:
		vestibulo.EnviaMensaje(mensaje)
	case strings.HasPrefix(mensaje.texto, CMD_CREAR):
		nombre := strings.TrimSuffix(strings.TrimPrefix(mensaje.texto, CMD_CREAR+" "), "\n")
		vestibulo.CreaSala(mensaje.conexion, nombre)
	case strings.HasPrefix(mensaje.texto, CMD_LISTA):
		vestibulo.ListaSalas(mensaje.conexion)
	case strings.HasPrefix(mensaje.texto, CMD_ENTRAR):
		nombre := strings.TrimSuffix(strings.TrimPrefix(mensaje.texto, CMD_ENTRAR+" "), "\n")
		vestibulo.EntrarSala(mensaje.conexion, nombre)
	case strings.HasPrefix(mensaje.texto, CMD_SALIR):
    nombre := strings.TrimSuffix(strings.TrimPrefix(mensaje.texto, CMD_SALIR+" "), "\n")
		vestibulo.DejarSala(mensaje.conexion, nombre)
	case strings.HasPrefix(mensaje.texto, CMD_NOMBRE):
		nombre := strings.TrimSuffix(strings.TrimPrefix(mensaje.texto, CMD_NOMBRE+" "), "\n")
		vestibulo.CambiaNombre(mensaje.conexion, nombre)
	case strings.HasPrefix(mensaje.texto, CMD_AYUDA):
		vestibulo.Ayuda(mensaje.conexion)
  case strings.HasPrefix(mensaje.texto, CMD_HISTORIAL):
    nombre := strings.TrimSuffix(strings.TrimPrefix(mensaje.texto, CMD_HISTORIAL+" "), "\n")
		vestibulo.Historial(mensaje.conexion, nombre)
	case strings.HasPrefix(mensaje.texto, CMD_TERMINAR):
		mensaje.conexion.Quit()
	}
}

/**
 * Intenta mandar el mensaje a las salas de chat en las que se
 * encuentra el cliente.   Si no está en sala alguna, se envía
 * un mensaje de error al cliente.
 */
func (vestibulo *Vestibulo) EnviaMensaje(mensaje *Mensaje) {
  for _, sala := range mensaje.conexion.salas {
    sala.Transmite(mensaje.String())
  }
	log.Printf("%v envió un mensaje\n", mensaje.conexion.nombre)
}

/**
 * Intenta crear una sala con el nombre dado.
 */
func (vestibulo *Vestibulo) CreaSala(conexion *Conexion, nombre string) {
	if vestibulo.salas[nombre] != nil {
		conexion.saliente <- ERROR_CREAR
		log.Printf("%v trató de crear una sala con un nombre ocupado\n", conexion.nombre)
		return
	}
	sala := NuevaSala(nombre)
	vestibulo.salas[nombre] = sala
	conexion.saliente <- fmt.Sprintf(EVENTO_PERSONAL_CREAR, sala.nombre)
	log.Println("client created chat room")
}

/**
 * Intenta añadir al cliente a la sala con el nombre dado,
 * siempre y cuando éste exista.
 */
func (vestibulo *Vestibulo) EntrarSala(conexion *Conexion, nombre string) {
	if vestibulo.salas[nombre] == nil {
		conexion.saliente <- ERROR_ENTRAR
		log.Printf("%v intentó entrar a una sala que no existe\n", conexion.nombre)
		return
	}
	vestibulo.salas[nombre].Agrega(conexion)
	log.Printf("%v entró a la sala %v", conexion.nombre, nombre)
}

/**
 * Envía el historial de la sala al cliente.
 */
func (vestibulo *Vestibulo) Historial(conexion *Conexion, nombre string) {
  if vestibulo.salas[nombre] == nil {
		conexion.saliente <- ERROR_HISTORIAL
		log.Printf("%v solicitó el historial de una sala de chat que no existe\n", conexion.nombre)
		return
	}
  sala, ok := conexion.salas[nombre]
  if ok {
    sala.Historial(conexion)
    log.Printf("%v solicitó el historial de la sala %v\n", conexion.nombre ,nombre)
    return
  }
  conexion.saliente <- ERROR_HISTORIAL
  log.Printf("%v solicitó el historial de una sala de chat a la que no pertenece\n", conexion.nombre)
}


/**
 * Elimina al cliente de la sala con el nombre dado.
 */
func (vestibulo *Vestibulo) DejarSala(conexion *Conexion, nombre string) {
  sala, ok := conexion.salas[nombre]
  if ok {
    sala.Elimina(conexion)
  	log.Printf("%v dejó la sala %v\n", conexion.nombre, nombre)
    return
	}
  conexion.saliente <- ERROR_SALIR
  log.Printf("%v intentó dejar una sala que no existe", conexion.nombre)
}

/**
 * Cambia el nombre del cliente al nombre dado.
 */
func (vestibulo *Vestibulo) CambiaNombre(conexion *Conexion, nombre string) {
	if len(conexion.salas) == 0 {
		conexion.saliente <- fmt.Sprintf(EVENTO_PERSONAL_NOMBRE, nombre)
	} else {
    for _, sala := range(conexion.salas) {
		  sala.Transmite(fmt.Sprintf(EVENTO_SALA_NOMBRE, conexion.nombre, nombre))
    }
	}
	conexion.nombre = nombre
	log.Println("client changed their name")
}

/**
 * Envía al cliente la lista de las salas en existencia.
 */
func (vestibulo *Vestibulo) ListaSalas(conexion *Conexion) {
	conexion.saliente <- "\nSalas de Chat:\n"
	for nombre := range vestibulo.salas {
		conexion.saliente <- fmt.Sprintf("%s\n", nombre)
	}
	conexion.saliente <- "\n"
	log.Printf("%v pidió la lista de las salas", conexion.nombre)
}

/**
 * Envía al cliente la lista de los posibles comandos.
 */
func (vestibulo *Vestibulo) Ayuda(conexion *Conexion) {
	conexion.saliente <- "\n"
	conexion.saliente <- "Comandos:\n"
	conexion.saliente <- "/ayuda - muestra todos los comandos\n"
  conexion.saliente <- "/historial - muestra todos los mensajes y noticias en la sala\n"
	conexion.saliente <- "/lista - lists all chat rooms\n"
	conexion.saliente <- "/crear foo - crea una sala de chat llamada foo\n"
	conexion.saliente <- "/entrar foo - entra a la sala de chat foo\n"
	conexion.saliente <- "/salir foo - sale de la sala foo\n"
	conexion.saliente <- "/nombre foo - cambia tu nombre a foo\n"
	conexion.saliente <- "/terminar - termina el programa\n"
	conexion.saliente <- "\n"
	log.Println("client requested help")
}
