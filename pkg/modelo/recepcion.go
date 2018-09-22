package modelo

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
    sala := NuevaSala("recepcion", nil)
    recepcion := Recepcion{sala}
    return &recepcion
}

/**
 * Añade al cliente a la sala de chat.
 */
func (recepcion *Recepcion) Agrega(conexion *Conexion) {
	recepcion.conexiones[conexion] = true
}

/**
 * Elimina al cliente de la recepción.
 */
func (recepcion *Recepcion) Elimina(conexion *Conexion) {
	delete (recepcion.conexiones, conexion)
}
