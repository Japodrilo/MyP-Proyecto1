package modelo

import (
  "time"
  "fmt"
  "strings"
)

/**
 * Un mensaje contiene información sobre el remitente, la hora
 * en el que el mensaje fue enviado y el texto del mensaje.
 */
type Mensaje struct {
	hora   time.Time
	conexion *Conexion
	texto   string
}

/**
 * Crea un nuevo mensaje con la hora, la conexión y el texto dados.
 */
func NuevoMensaje(hora time.Time, conexion *Conexion, texto string) *Mensaje {
	return &Mensaje{
		hora:     hora,
		conexion: conexion,
		texto:    texto,
	}
}

/**
 * Regresa una representación del mensaje como cadena.
 */
func (mensaje *Mensaje) String() string {
	return fmt.Sprintf("%s - %s: %s\n", mensaje.hora.Format(time.Kitchen), mensaje.conexion.nombre, strings.TrimPrefix(mensaje.texto, CMD_MENSAJE_PUBLICO+" "))
}
