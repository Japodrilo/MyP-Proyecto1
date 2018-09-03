package modelo

import (
  "fmt"
  "time"
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
  if conexion.salas[sala.nombre] != nil {
    return
  }
	conexion.salas[sala.nombre] = sala
	sala.conexiones = append(sala.conexiones, conexion)
	sala.TransmiteOtros(fmt.Sprintf(EVENTO_ENTRA_SALA, time.Now().Format(time.Kitchen), conexion.nombre), conexion)
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
	for i, otherConexion := range sala.conexiones {
		if conexion == otherConexion {
			sala.conexiones = append(sala.conexiones[:i], sala.conexiones[i+1:]...)
			break
		}
	}
	delete(conexion.salas, sala.nombre)
  sala.Transmite(fmt.Sprintf(EVENTO_DEJA_SALA, time.Now().Format(time.Kitchen), conexion.nombre, sala.nombre))
}

/**
 * Envía el mensaje a todos los clientes conectados al servidor.
 */
func (sala *Sala) Transmite(mensaje string) {
	sala.mensajes = append(sala.mensajes, mensaje)
	for _, conexion := range sala.conexiones {
		conexion.saliente <- mensaje
	}
}

/**
 * Envía el mensaje a todos los demás clientes en el servidor.
 */
func (sala *Sala) TransmiteOtros(mensaje string, conexion *Conexion) {
	sala.mensajes = append(sala.mensajes, mensaje)
	for _, otraConexion := range sala.conexiones {
    if conexion != otraConexion {
		    otraConexion.saliente <- mensaje
    }
	}
}
