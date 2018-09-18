package modelo

import (
  "fmt"
)

/**
 * Una sala de chat contiene su nombre, una lista de los clientes actualmente
 * conectados y la historia de todos los mensajes y notificaciones que han sido
 * enviados en la sala.
 */
type Sala struct {
	nombre      string
  propietario *Conexion
	conexiones  map[*Conexion]bool
  invitados   map[*Conexion]bool
	mensajes    []string
}

/**
 * Crea una sala de chat vacía con el nombre dado.
 */
func NuevaSala(nombre string, propietario *Conexion) *Sala {
	return &Sala{
		nombre:      nombre,
    propietario: propietario,
		conexiones:  make(map[*Conexion]bool),
    invitados:   make(map[*Conexion]bool),
		mensajes:    make([]string, 0),
	}
}

/**
 * Añade al cliente a la sala de chat.
 */
func (sala *Sala) Agrega(conexion *Conexion) {
  if conexion.salas[sala.nombre] != nil {
    return
  }
	conexion.salas[sala.nombre] = sala
	sala.conexiones[conexion] = true
}

/**
 * Añade al cliente a la lista de invitados.
 */
 func (sala *Sala) Invita(conexion *Conexion) {
   conexion.saliente <- fmt.Sprintf("...INVITATION TO JOIN %v ROOM BY %v\n", sala.nombre, sala.propietario.nombre)
   conexion.saliente <- fmt.Sprintf("...TO JOIN: JOIN %v\n", sala.nombre)
   if sala.invitados[conexion] {
     return
   }
   sala.invitados[conexion] = true
 }

/**
 * Envía al cliente el historial de todos los mensajes y noticias
 * que han sido enviados en la sala desde su creación.
 */
func (sala *Sala) Historial(conexion *Conexion) {
  for _, mensaje := range sala.mensajes {
		conexion.saliente <- mensaje
	}
}

/**
 * Elimina al cliente de la sala de chat.
 */
func (sala *Sala) Elimina(conexion *Conexion) {
  delete(sala.conexiones, conexion)
	delete(conexion.salas, sala.nombre)
  delete(sala.invitados, conexion)
}

/**
 * Envía el mensaje a todos los clientes conectados al servidor.
 */
func (sala *Sala) Transmite(mensaje string) {
	sala.mensajes = append(sala.mensajes, mensaje)
	for conexion, _ := range sala.conexiones {
		conexion.saliente <- mensaje
	}
}

/**
 * Envía el mensaje a todos los demás clientes en el servidor.
 */
func (sala *Sala) TransmiteOtros(mensaje string, conexion *Conexion) {
	sala.mensajes = append(sala.mensajes, mensaje)
	for otraConexion, _ := range sala.conexiones {
    if conexion != otraConexion {
		    otraConexion.saliente <- mensaje
    }
	}
}
