package modelo

import (
  "fmt"
  "time"
)



/**
 * La recepción es una sala de chat especial, su nombre por omisión es
 * "recepcion", tiene una lista de los clientes actualmente conectados y
 * la historia de todos los mensajes y notificaciones que han sido enviados.
 * Un cliente no debería de dejar la recepción hasta terminar sesión.
 */
type Recepcion struct {
	*Sala
}

/**
 * Crea una recepción vacía con el nombre por omisión.
 */
func NuevaRecepcion() *Recepcion {
  sala := NuevaSala("recepcion")
  recepcion := Recepcion{sala}
  return &recepcion
}

/**
 * Añade al cliente a la sala de chat.
 */
func (recepcion *Recepcion) Agrega(conexion *Conexion) {
	recepcion.conexiones = append(recepcion.conexiones, conexion)
	recepcion.TransmiteOtros(fmt.Sprintf(EVENTO_INICIO_SESION, time.Now().Format(time.Kitchen), conexion.nombre), conexion)
}

/**
 * Elimina al cliente de la recepción.
 */
func (recepcion *Recepcion) Elimina(conexion *Conexion) {
	for i, otherConexion := range recepcion.conexiones {
		if conexion == otherConexion {
			recepcion.conexiones = append(recepcion.conexiones[:i], recepcion.conexiones[i+1:]...)
			break
		}
	}
  recepcion.Transmite(fmt.Sprintf(EVENTO_DESCONEXION_USUARIO, time.Now().Format(time.Kitchen), conexion.nombre))
}
