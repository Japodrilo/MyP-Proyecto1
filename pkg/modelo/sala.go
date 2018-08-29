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
	conexiones  []*Conexion
	mensajes    []string
}

/**
 * Crea una sala de chat vacía con el nombre dado.
 */
func NuevaSala(nombre string) *Sala {
	return &Sala{
		nombre:      nombre,
		conexiones:  make([]*Conexion, 0),
		mensajes:    make([]string, 0),
	}
}

/**
 * Añade al cliente a la sala de chat.
 */
func (sala *Sala) Agrega(conexion *Conexion) {
	conexion.salas[sala.nombre] = sala
	sala.conexiones = append(sala.conexiones, conexion)
	sala.Transmite(fmt.Sprintf(EVENTO_ENTRA_SALA, conexion.nombre))
}

/**
 * Envía al cliente el historial de todos los mensajes y noticias
 * que han sido enviados en la sala desde su creación.
 */
func (sala *Sala) Historial(conexion *Conexion) {
  for _, message := range sala.mensajes {
		conexion.saliente <- message
	}
}

/**
 * Elimina al cliente de la sala de chat.
 */
func (sala *Sala) Elimina(conexion *Conexion) {
	for i, otherConexion := range sala.conexiones {
		if conexion == otherConexion {
			sala.conexiones = append(sala.conexiones[:i], sala.conexiones[i+1:]...)
			break
		}
	}
	delete(conexion.salas, sala.nombre)
  sala.Transmite(fmt.Sprintf(EVENTO_DEJA_SALA, conexion.nombre, sala.nombre))
}

/**
 * Envía el mensaje a todos los clientes en la sala.
 */
func (sala *Sala) Transmite(message string) {
	sala.mensajes = append(sala.mensajes, message)
	for _, conexion := range sala.conexiones {
		conexion.saliente <- message
	}
}
