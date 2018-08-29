package modelo

import (
  "fmt"
  "net"
  "testing"
  "time"
)

func TestNewMensaje(t *testing.T) {
  conn, err := net.Dial("tcp", "127.0.0.1:1234")
  if err != nil {
    fmt.Println(err)
  }
  conexion := NewConexion(conn)
  hora := time.Date(1980, time.June, 6, 15, 0, 0, 0, time.UTC)
  mensaje := NewMensaje(hora , conexion, "Mensaje de prueba")
  resultado := mensaje.String()
  esperado := fmt.Sprintf("3:00PM - Anonymous: Mensaje de prueba\n")
  if resultado != esperado {
    t.Errorf("Se esperaba %v, se obtuvo %v", esperado, resultado)
  }
}
