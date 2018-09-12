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
  recepcion   *Recepcion
  conexionesT map[*Conexion]bool
  conexiones  map[string]*Conexion
	salas       map[string]*Sala
	entrante    chan *Mensaje
	join        chan *Conexion
	leave       chan *Conexion
}

/**
 * Crea un Vestibulo que empieza a escuchar en sus canales.
 */
func NuevoVestibulo() *Vestibulo {
	vestibulo := &Vestibulo{
    recepcion:   NuevaRecepcion(),
		conexionesT: make(map[*Conexion]bool),
    conexiones:  make(map[string]*Conexion),
		salas:       make(map[string]*Sala),
		entrante:    make(chan *Mensaje),
		join:        make(chan *Conexion),
		leave:       make(chan *Conexion),
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
	vestibulo.conexionesT[conexion] = true
	conexion.saliente <- "...Bienvenido al servidor, identifícate por favor:\n"
  conexion.saliente <- "...IDENTIFY username\n"
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
  delete(vestibulo.conexiones, conexion.nombre)
  for _, sala := range conexion.salas {
    sala.Elimina(conexion)
  }
	close(conexion.saliente)
	log.Printf("Se cerró el canal saliente de %v\n", conexion.nombre)
}

/**
 * Función auxiliar para extraer el prefijo de un mensaje.
 */
func (vestibulo *Vestibulo) parse(mensaje *Mensaje) (string, []string, *Conexion){
  prefijo := ""
  argumentos := strings.Fields(mensaje.texto)
  conexion := mensaje.conexion
  switch {
  case mensaje.texto == CMD_AYUDA:
    prefijo = CMD_AYUDA
	case strings.HasPrefix(mensaje.texto, CMD_CREAR):
    prefijo = CMD_CREAR
  case strings.HasPrefix(mensaje.texto, CMD_ENTRAR):
    prefijo = CMD_ENTRAR
  case strings.HasPrefix(mensaje.texto, CMD_HISTORIAL):
    prefijo = CMD_HISTORIAL
  case strings.HasPrefix(mensaje.texto, CMD_INVITAR):
    prefijo = CMD_INVITAR
  case strings.HasPrefix(mensaje.texto, CMD_MENSAJE_DIRECTO):
    prefijo = CMD_MENSAJE_DIRECTO
  case strings.HasPrefix(mensaje.texto, CMD_MENSAJE_PUBLICO):
    prefijo = CMD_MENSAJE_PUBLICO
  case strings.HasPrefix(mensaje.texto, CMD_MENSAJE_SALA):
    prefijo = CMD_MENSAJE_SALA
  case strings.HasPrefix(mensaje.texto, CMD_NOMBRE):
    prefijo = CMD_NOMBRE
  case mensaje.texto == CMD_SALAS:
		prefijo = CMD_SALAS
	case strings.HasPrefix(mensaje.texto, CMD_SALIR):
    prefijo = CMD_SALIR
  case strings.HasPrefix(mensaje.texto, CMD_STATUS):
    prefijo = CMD_STATUS
	case strings.HasPrefix(mensaje.texto, CMD_TERMINAR):
    prefijo = CMD_TERMINAR
  case mensaje.texto == CMD_USUARIOS:
    prefijo = CMD_USUARIOS
	}
  return prefijo, argumentos, conexion
}

/**
 * Se encarga de los mensajes que llegan al vestíbulo.   Si el mensaje
 * contiene un comando, éste es ejecutado por el vestíbulo.   Si no,
 * envía un mensaje de error al cliente.
 */
func (vestibulo *Vestibulo) Maneja(mensaje *Mensaje) {
  prefijo, argumentos, conexion := vestibulo.parse(mensaje)
  switch {
  case prefijo == CMD_NOMBRE && len(argumentos) > 1:
    vestibulo.CambiaNombre(conexion, argumentos[1])
  case conexion.nombre == "":
    vestibulo.Identificate(conexion)
    return
  }
  switch prefijo {
  case CMD_AYUDA:
    vestibulo.Ayuda(conexion)
  case CMD_CREAR:
    vestibulo.CreaSala(conexion, argumentos)
  case CMD_ENTRAR:
    vestibulo.EntraSala(conexion, argumentos)
  case CMD_HISTORIAL:
    vestibulo.Historial(conexion, argumentos[1])
  case CMD_INVITAR:
    vestibulo.Invitar(conexion, argumentos)
  case CMD_MENSAJE_DIRECTO:
    vestibulo.MensajeDirecto(conexion, argumentos)
  case CMD_MENSAJE_PUBLICO:
    vestibulo.MensajePublico(conexion, argumentos)
  case CMD_MENSAJE_SALA:
    vestibulo.MensajeSala(conexion, argumentos)
  case CMD_SALAS:
    vestibulo.ListaSalas(conexion)
  case CMD_SALIR:
    vestibulo.DejarSala(conexion, argumentos[1])
  case CMD_STATUS:
    vestibulo.Estado(conexion, argumentos)
  case CMD_TERMINAR:
    conexion.Terminar()
  case CMD_USUARIOS:
    vestibulo.Usuarios(conexion)
  case "":
    vestibulo.noReconocido(conexion, strings.Join(argumentos, " "))
  }
}

/**
 * Revisa si hay argumentos para el comando ejecutado por el cliente.
 */
func (vestibulo *Vestibulo) hayArgumentos(conexion *Conexion, argumentos []string, numeroArgumentos int) (ok bool) {
  if len(argumentos) < numeroArgumentos {
    vestibulo.mensajesValidos(conexion)
    return false
  }
  return true
}

/**
 * Indica al usuario cuáles son los mensajes válidos.
 */
func (vestibulo *Vestibulo) mensajesValidos(conexion *Conexion) {
  conexion.saliente <- "...LOS MENSAJES VÁLIDOS SON:\n"
  conexion.saliente <- "...IDENTIFY username\n"
  conexion.saliente <- "...STATUS userStatus = {ACTIVE, AWAY, BUSY}\n"
  conexion.saliente <- "...MESSAGE username messageContent\n"
  conexion.saliente <- "...PUBLICMESSAGE messageContent\n"
  conexion.saliente <- "...CREATEROOM roomname\n"
  conexion.saliente <- "...INVITE roomname user1 user2 ...\n"
  conexion.saliente <- "...JOINROOM roomname\n"
  conexion.saliente <- "...ROOMESSAGE roomname messageContent\n"
  conexion.saliente <- "...DISCONNECT\n"
}


/**
 * Revisa si el cliente ya se identificó, de lo contrario,
 * le pide que se identifique.
 */
func (vestibulo *Vestibulo) Identificate(conexion *Conexion) {
  if conexion.nombre == "" {
    conexion.saliente <- "...Debes identificarte primero.\n"
    conexion.saliente <- "...IDENTIFY username\n"
  }
}

/**
 * Cuando el cliente envía un comando no reconocido, se le
 * notifica al mismo con un error, y se imprimen en el log
 * los detalles del mensaje.
 */
func (vestibulo *Vestibulo) noReconocido(conexion *Conexion, texto string) {
  conexion.saliente <- ERROR_DESCONOCIDO
  vestibulo.mensajesValidos(conexion)
  log.Printf("%v envió un comando no reconocido: %v", conexion.nombre, texto)
}

/**
 * Le da formato al mensaje dependiendo de su tipo, directo, público o de sala.
 */
func (vestibulo *Vestibulo) formateaMensaje(conexion *Conexion, argumentos []string) string {
  switch argumentos[0] {
  case CMD_MENSAJE_PUBLICO:
    return fmt.Sprintf("%s - %s: %s\n", time.Now().Format(time.Kitchen), conexion.nombre, strings.Join(argumentos[1:], " "))
  case CMD_MENSAJE_SALA:
    return fmt.Sprintf("%s - %s: %s\n", time.Now().Format(time.Kitchen), conexion.nombre, strings.Join(argumentos[2:], " "))
  case CMD_MENSAJE_DIRECTO:
    return fmt.Sprintf("%s - %s: %s\n", time.Now().Format(time.Kitchen), conexion.nombre, strings.Join(argumentos[2:], " "))
  default:
    return ""
  }
}

/**
 * Intenta invitar a los usuarios a una sala de chat.
 */
func (vestibulo *Vestibulo) Invitar(conexion *Conexion, argumentos []string) {
  if !vestibulo.hayArgumentos(conexion, argumentos, 3) {
    return
  }
  sala := vestibulo.salas[argumentos[1]]
  switch {
  case sala == nil:
    conexion.saliente <- ERROR_ENTRAR_NO_EXISTE
  case sala.propietario.nombre != conexion.nombre:
    conexion.saliente <- ERROR_INVITACION
  default:
    for _, nombre := range argumentos[2:] {
      cliente, ok := vestibulo.conexiones[nombre]
      switch {
        case ok && nombre != sala.propietario.nombre:
          sala.Invita(cliente)
        case ok:
          continue
        default:
          conexion.saliente <- fmt.Sprintf("...USUARIO %v NO ENCONTRADO\n", nombre)
      }
    }
  }
}

/**
 * Intenta mandar el mensaje al usuario objetivo.
 */
func (vestibulo *Vestibulo) MensajeDirecto(conexion *Conexion, argumentos []string) {
  if !vestibulo.hayArgumentos(conexion, argumentos, 3) {
    return
  }
  destinatario, ok := vestibulo.conexiones[argumentos[1]]
  if ok {
    conexion.saliente <- MENSAJE_ENVIADO
    destinatario.saliente <- vestibulo.formateaMensaje(conexion, argumentos)
    log.Printf("%v envió un mensaje público\n", conexion.nombre)
  } else {
    conexion.saliente <- "...USUARIO NO ENCONTRADO\n"
  }
}

/**
 * Intenta mandar el mensaje a todos los usuarios conectados.
 */
func (vestibulo *Vestibulo) MensajePublico(conexion *Conexion, argumentos []string) {
  if !vestibulo.hayArgumentos(conexion, argumentos, 2) {
    return
  }
  vestibulo.recepcion.TransmiteOtros(vestibulo.formateaMensaje(conexion, argumentos), conexion)
  conexion.saliente <- MENSAJE_ENVIADO
	log.Printf("%v envió un mensaje público\n", conexion.nombre)
}

/**
 * Intenta mandar el mensaje a la sala especificada.
 */
func (vestibulo *Vestibulo) MensajeSala(conexion *Conexion, argumentos []string) {
  if !vestibulo.hayArgumentos(conexion, argumentos, 3) {
    return
  }
  sala := conexion.salas[argumentos[1]]
  if sala != nil {
    sala.TransmiteOtros(vestibulo.formateaMensaje(conexion, argumentos), conexion)
    conexion.saliente <- MENSAJE_ENVIADO
    log.Printf("%v envió un mensaje en la sala %v\n", conexion.nombre, argumentos[1])
    return
  }
  conexion.saliente <- fmt.Sprintf(ERROR_SALA, argumentos[1])
}

/**
 * Intenta crear una sala con el nombre dado.
 */
func (vestibulo *Vestibulo) CreaSala(conexion *Conexion, argumentos []string) {
  if !vestibulo.hayArgumentos(conexion, argumentos, 2) {
    return
  }
  if vestibulo.salas[argumentos[1]] != nil {
		conexion.saliente <- ERROR_NOMBRE
		log.Printf("%v trató de crear una sala con un nombre ocupado\n", conexion.nombre)
		return
	}
	sala := NuevaSala(argumentos[1], conexion)
	vestibulo.salas[argumentos[1]] = sala
	conexion.saliente <- fmt.Sprintf(EVENTO_PERSONAL_CREAR, time.Now().Format(time.Kitchen), sala.nombre)
  vestibulo.salas[argumentos[1]].Agrega(conexion)
	log.Printf("%v creó la sala %v", conexion.nombre, argumentos[1])
}

/**
 * Intenta añadir al cliente a la sala con el nombre dado,
 * siempre y cuando éste exista.
 */
func (vestibulo *Vestibulo) EntraSala(conexion *Conexion, argumentos []string) {
  if !vestibulo.hayArgumentos(conexion, argumentos, 2) {
    return
  }
  sala := vestibulo.salas[argumentos[1]]
	if sala == nil {
		conexion.saliente <- ERROR_ENTRAR_NO_EXISTE
		log.Printf("%v intentó entrar a una sala que no existe\n", conexion.nombre)
		return
	}
  if sala.invitados[conexion] {
	  sala.Agrega(conexion)
	  log.Printf("%v entró a la sala %v", conexion.nombre, argumentos[1])
  } else {
    conexion.saliente <- ERROR_ENTRAR_INVITACION
    log.Printf("%v intentó entrar a una sala sin invitación\n", conexion.nombre)
  }
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
  for usuario, _ := range(vestibulo.conexiones) {
      conexion.saliente <- fmt.Sprintf("%v\n", usuario)
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
  conexion.saliente <- fmt.Sprintf(ERROR_SALA, nombre)
  log.Printf("%v intentó dejar una sala a la que no pertenece", conexion.nombre)
}

/**
 * Cambia el nombre del cliente al nombre dado.
 */
func (vestibulo *Vestibulo) CambiaNombre(conexion *Conexion, nombre string) {
  for otraConexion, _ := range vestibulo.conexiones {
		if nombre == otraConexion {
      conexion.saliente <- ERROR_NOMBRE
			return
		}
	}
	conexion.saliente <- fmt.Sprintf(EVENTO_PERSONAL_NOMBRE, time.Now().Format(time.Kitchen), nombre)
	log.Printf("Cliente %v cambió su nombre a %v", conexion.serial, nombre)
  delete(vestibulo.conexiones, conexion.nombre)
  vestibulo.conexiones[nombre] = conexion
  vestibulo.conexionesT[conexion] = false
  conexion.SetNombre(nombre)
}

/**
 * Cambia el nombre del cliente al nombre dado.
 */
func (vestibulo *Vestibulo) Estado(conexion *Conexion, argumentos []string) {
  if !vestibulo.hayArgumentos(conexion, argumentos, 2) {
    return
  }
  conexion.SetStatus(argumentos[1])
	log.Printf("%v cambió su estado a %v", conexion.nombre, argumentos[1])
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
	conexion.saliente <- "Mensajes válidos:\n"
  conexion.saliente <- "CREATEROOM roomname - crea una sala de chat llamada roomname\n"
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
