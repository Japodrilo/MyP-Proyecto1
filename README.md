# MyP-Proyecto1

Repositorio para el Proyecto 1 en Modelado y Programación

Lenguaje utilizado: Go
Versiones: go version go1.10.3 darwin/amd64
           go version go1.10.1 linux/amd64
           
Dependencias: [gotk3](https://github.com/gotk3/gotk3)

Este proyecto tiene dos clases principales, cmd/servidor/servidor.go
y cmd/cliente/cliente.go, ambas generan archivos binarios ejecutables
al compilar que son, respectivamente, un servidor y un cliente de chat
en una conexión TCP.

Antes de compilar, gotk3 debe estar instalado en el GOPATH, en la ruta
src/github.com/gotk3/gotk3.   Para instalarlo, basta utilizar el comando

```bash
$ go get github.com/gotk3/gotk3/gtk
```

Para instalarlo en Ubuntu/Debian, se necesitan las siguientes dependencias:
GTK 3.6-3.16, GLib 2.36-2.40, y Cairo 1.10 or 1.12.

Para instrucciones detalladas, consulte: [installation](https://github.com/gotk3/gotk3/wiki#installation)
Las dependencias pueden obtenerse (Ubuntu/Debian) con el comando.

```bash
$ sudo apt-get install libgtk-3-dev libcairo2-dev libglib2.0-dev
```

Para compilar el proyecto, la carpeta base debe de estar en el GOPATH,
en la ruta src/github.com/Japodrilo/MyP-Proyecto1, de donde puede utilizarse
el comando "go install ./..." para generar los archivos binarios, que se
guardarán en la ruta dada por GOBIN.   También, el cliente puede correrse
con el comando "go run cliente.go" directamente en la carpeta donde se encuentra:
análogamente, el servidor puede correrse con el comando "go run servidor.go puerto"
Las clases en el modelo contienen pruebas unitarias que corren con el comando
"go test ./...", sin embargo, todas las pruebas fallan por omisión.
