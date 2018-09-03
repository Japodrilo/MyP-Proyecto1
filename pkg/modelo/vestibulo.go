package modelo

import (
  "log"
  "strings"
  "fmt"
  "time"
)

/**
 * Un Vestibulo administra a los clientes conectados, las salas existentes,
 * y distribuye los mensajes.
 */
type Vestibulo struct {
  recepcion  *Recepcion
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
    recepcion:  NuevaRecepcion(),
		conexiones: make([]*Conexion, 0),
		salas:      make(map[string]*Sala),
		entrante:   make(chan *Mensaje),
		join:       make(chan *Conexion),
		leave:      make(chan *Conexion),
	}
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
				vestibulo.Maneja(mensaje)
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
  vestibulo.recepcion.Agrega(conexion)
  log.Printf("Conexión recibida de %v\n" +
             "Serial de Conexion: %v\n", conexion.conn.RemoteAddr(), conexion.serial)
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
  vestibulo.recepcion.Elimina(conexion)
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
 * Función auxiliar para extraer el prefijo de un mensaje.
 */
func (vestibulo *Vestibulo) parse(mensaje *Mensaje) (string, string, *Conexion){
  conexion := mensaje.conexion
  switch {
  case mensaje.texto == CMD_AYUDA:
    return CMD_AYUDA, "", conexion
	case strings.HasPrefix(mensaje.texto, CMD_CREAR):
		nombre := strings.TrimSuffix(strings.TrimPrefix(mensaje.texto, CMD_CREAR+" "), "\n")
    return CMD_CREAR, nombre, conexion
  case strings.HasPrefix(mensaje.texto, CMD_ENTRAR):
		nombre := strings.TrimSuffix(strings.TrimPrefix(mensaje.texto, CMD_ENTRAR+" "), "\n")
    return CMD_ENTRAR, nombre, conexion
  case strings.HasPrefix(mensaje.texto, CMD_HISTORIAL):
    nombre := strings.TrimSuffix(strings.TrimPrefix(mensaje.texto, CMD_HISTORIAL+" "), "\n")
    return CMD_HISTORIAL, nombre, conexion
  case strings.HasPrefix(mensaje.texto, CMD_MENSAJE):
		texto := strings.TrimSuffix(strings.TrimPrefix(mensaje.texto, CMD_MENSAJE+" "), "\n")
    return CMD_MENSAJE, texto, conexion
  case strings.HasPrefix(mensaje.texto, CMD_NOMBRE):
		nombre := strings.TrimSuffix(strings.TrimPrefix(mensaje.texto, CMD_NOMBRE+" "), "\n")
    return CMD_NOMBRE, nombre, conexion
  case mensaje.texto == CMD_SALAS:
		return CMD_SALAS, "", conexion
	case strings.HasPrefix(mensaje.texto, CMD_SALIR):
    nombre := strings.TrimSuffix(strings.TrimPrefix(mensaje.texto, CMD_SALIR+" "), "\n")
    return CMD_SALIR, nombre, conexion
	case strings.HasPrefix(mensaje.texto, CMD_TERMINAR):
    return CMD_TERMINAR, "", conexion
  case mensaje.texto == CMD_USUARIOS:
    return CMD_USUARIOS, "", conexion
  default:
    return "", mensaje.texto, conexion
	}
}

/**
 * Se encarga de los mensajes que llegan al vestíbulo.   Si el mensaje
 * contiene un comando, éste es ejecutado por el vestíbulo.   Si no,
 * envía un mensaje de error al cliente.
 */
func (vestibulo *Vestibulo) Maneja(mensaje *Mensaje) {
  prefijo, nombre, conexion := vestibulo.parse(mensaje)
  switch prefijo {
  case CMD_AYUDA:
    vestibulo.Ayuda(conexion)
  case CMD_CREAR:
    vestibulo.CreaSala(conexion, nombre)
  case CMD_ENTRAR:
    vestibulo.EntraSala(conexion, nombre)
  case CMD_HISTORIAL:
    vestibulo.Historial(conexion, nombre)
  case CMD_MENSAJE:
    vestibulo.EnviaMensaje(conexion, mensaje.String())
  case CMD_NOMBRE:
    vestibulo.CambiaNombre(conexion, nombre)
  case CMD_SALAS:
    vestibulo.ListaSalas(conexion)
  case CMD_SALIR:
    vestibulo.DejarSala(conexion, nombre)
  case CMD_TERMINAR:
    conexion.Terminar()
  case CMD_USUARIOS:
    vestibulo.Usuarios(conexion)
  case "":
    vestibulo.noReconocido(conexion, nombre)
  }
}

/**
 * Cuando el cliente envía un comando no reconocido, se le
 * notifica al mismo con un error, y se imprimen en el log
 * los detalles del mensaje.
 */
func (vestibulo *Vestibulo) noReconocido(conexion *Conexion, texto string) {
  conexion.saliente <- ERROR_DESCONOCIDO
  log.Printf("%v envió un comando no reconocido: %v", conexion.nombre, texto)
}

/**
 * Intenta mandar el mensaje a las salas de chat en las que se
 * encuentra el cliente.
 */
func (vestibulo *Vestibulo) EnviaMensaje(conexion *Conexion, texto string) {
  vestibulo.recepcion.Transmite(texto)
  for _, sala := range conexion.salas {
    sala.Transmite(texto)
  }
	log.Printf("%v envió un mensaje\n", conexion.nombre)
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
	conexion.saliente <- fmt.Sprintf(EVENTO_PERSONAL_CREAR, time.Now().Format(time.Kitchen), sala.nombre)
  vestibulo.salas[nombre].Agrega(conexion)
	log.Printf("%v creó la sala %v", conexion.nombre, nombre)
}

/**
 * Intenta añadir al cliente a la sala con el nombre dado,
 * siempre y cuando éste exista.
 */
func (vestibulo *Vestibulo) EntraSala(conexion *Conexion, nombre string) {
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
  if nombre == "recepcion" {
		vestibulo.recepcion.Historial(conexion)
		log.Printf("%v solicitó el historial de la recepción\n", conexion.nombre)
		return
	}
  sala, ok := conexion.salas[nombre]
  if ok {
    sala.Historial(conexion)
    log.Printf("%v solicitó el historial de la sala %v\n", conexion.nombre ,nombre)
    return
  }
  conexion.saliente <- ERROR_HISTORIAL
  log.Printf("%v solicitó el historial de una sala de chat a la que no pertenece o no especificó nombre de sala\n", conexion.nombre)
}

/**
 * Envía la lista de usuarios al cliente.
 */
func (vestibulo *Vestibulo) Usuarios(conexion *Conexion) {
  for i, usuario := range(vestibulo.conexiones) {
    conexion.saliente <- fmt.Sprintf("%v. %v\n", i, usuario.nombre)
  }
  log.Printf("%v solicitó la lista de usuarios\n", conexion.nombre)
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
  conexion.saliente <- fmt.Sprintf(ERROR_SALIR, nombre)
  log.Printf("%v intentó dejar una sala a la que no pertenece", conexion.nombre)
}

/**
 * Cambia el nombre del cliente al nombre dado.
 */
func (vestibulo *Vestibulo) CambiaNombre(conexion *Conexion, nombre string) {
	if len(conexion.salas) == 0 {
		conexion.saliente <- fmt.Sprintf(EVENTO_PERSONAL_NOMBRE, time.Now().Format(time.Kitchen), nombre)
	} else {
    for _, sala := range(conexion.salas) {
		  sala.Transmite(fmt.Sprintf(EVENTO_SALA_NOMBRE, time.Now().Format(time.Kitchen), conexion.nombre, nombre))
    }
	}
	log.Printf("%v cambió su nombre a %v", conexion.nombre, nombre)
  conexion.SetNombre(nombre)
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
	conexion.saliente <- "|=AYUDA - muestra todos los comandos\n"
  conexion.saliente <- "|=CREAR foo - crea una sala de chat llamada foo\n"
  conexion.saliente <- "|=ENTRAR foo - entra a la sala de chat foo\n"
  conexion.saliente <- "|=HISTORIAL foo - muestra todos los mensajes y noticias en la sala foo\n"
  conexion.saliente <- "|=MENSAJE foo - envía el mensaje foo\n"
  conexion.saliente <- "|=NOMBRE foo - cambia tu nombre a foo\n"
	conexion.saliente <- "|=SALAS - muestra la lista de salas existentes\n"
	conexion.saliente <- "|=SALIR foo - sale de la sala foo\n"
	conexion.saliente <- "|=TERMINAR - termina el programa\n"
  conexion.saliente <- "|=USUARIOS - muestra la lista de todos los usuarios conectados\n"
	conexion.saliente <- "\n"
	log.Printf("%v solicitó ayuda", conexion.nombre)
}
