package modelo

const (
	CONN_TIPO = "tcp"

	CMD_PREFIJO   = "|="
	CMD_AYUDA     = CMD_PREFIJO + "AYUDA"
	CMD_CREAR     = CMD_PREFIJO + "CREAR"
	CMD_ENTRAR    = CMD_PREFIJO + "ENTRAR"
	CMD_HISTORIAL = CMD_PREFIJO + "HISTORIAL"
	CMD_MENSAJE   = CMD_PREFIJO + "MENSAJE"
	CMD_NOMBRE    = CMD_PREFIJO + "NOMBRE"
	CMD_SALAS     = CMD_PREFIJO + "SALAS"
	CMD_SALIR     = CMD_PREFIJO + "SALIR"
	CMD_TERMINAR  = CMD_PREFIJO + "TERMINAR"
	CMD_USUARIOS  = CMD_PREFIJO + "USUARIOS"

	CLIENTE_NOMBRE  = "Cliente %v"
	SERVIDOR_NOMBRE = "Servidor"

	ERROR_PREFIJO     = "Error: "
	ERROR_CREAR       = ERROR_PREFIJO + "Ya existe una sala con ese nombre.\n"
	ERROR_ENTRAR      = ERROR_PREFIJO + "No existe una sala con ese nombre.\n"
	ERROR_DESCONOCIDO = ERROR_PREFIJO + "Comando no reconocido.\n"
	ERROR_HISTORIAL   = ERROR_PREFIJO + "Solicitaste el historial de una sala a la que no perteneces o no especificaste nombre.\n"
	ERROR_SALIR       = ERROR_PREFIJO + "No perteneces a la sala %v.\n"

	EVENTO_PREFIJO             = "%v - Noticia: "
	EVENTO_ENTRA_SALA          = EVENTO_PREFIJO + "\"%s\" se unió a la sala.\n"
	EVENTO_DEJA_SALA           = EVENTO_PREFIJO + "\"%s\" salió de la sala %v.\n"
	EVENTO_DESCONEXION_USUARIO = EVENTO_PREFIJO + "\"%s\" terminó sesión.\n"
	EVENTO_INICIO_SESION       = EVENTO_PREFIJO + "\"%s\" inició sesión.\n"
	EVENTO_SALA_NOMBRE         = EVENTO_PREFIJO + "\"%s\" cambió su nombre a \"%s\".\n"
	EVENTO_PERSONAL_CREAR      = EVENTO_PREFIJO + "Creaste la sala \"%s\".\n"
	EVENTO_PERSONAL_NOMBRE     = EVENTO_PREFIJO + "Cambiaste tu nombre a \"%s\".\n"

	MSJ_CONECTA = "\n**Bienvenido al servidor! Escribe \"|=AYUDA\" para obtener la lista de comandos.**\n"
	MSJ_DESCONECTA = "\nDesconectado del servidor.\n"
)
