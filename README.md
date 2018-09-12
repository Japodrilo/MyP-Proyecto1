# MyP-Proyecto1

Repositorio para el Proyecto 1 en Modelado y Programación

Lenguaje utilizado: Go
Versión: go version go1.10.3 darwin/amd64

Este proyecto tiene dos clases principales, servidor.go y cliente.go,
ambas generan archivos binarios ejecutables al compilar que son,
respectivamente, un servidor y un cliente de chat en una conexión TCP.

Para compilarlo, la carpeta base debe de estar en el GOPATH, en la
ruta src/github.com/Japodrilo/MyP-Proyecto1, de donde puede utilizarse
el comando "go install ./..." para generar los archivos binarios.

El paquete modelo contiene las siguientes clases:
-conexion: contiene la información de cada cliente.
-constantes: tiene cadenas que se utilizan en todo el paquete.
-mensaje: contiene la información de un mensaje: hora, remitente y texto.
-recepcion: es una sala especial que representa la sala general de chat.
-sala: representa una sala de chat.
-vestibulo: administra las conexiones, recibe los mensajes, los interpreta
            y distribuye las tareas.

La única clase con pruebas reales es mensaje del paquete modelo, las
demás fallan por omisión.
