package modelo

const (
	CONN_PUERTO = ":5000"
	CONN_TIPO = "tcp"

	CMD_PREFIJO   = "/"
	CMD_CREAR     = CMD_PREFIJO + "crear"
	CMD_LISTA     = CMD_PREFIJO + "lista"
	CMD_ENTRAR    = CMD_PREFIJO + "entrar"
	CMD_SALIR     = CMD_PREFIJO + "salir"
	CMD_AYUDA     = CMD_PREFIJO + "ayuda"
	CMD_HISTORIAL = CMD_PREFIJO + "historial"
	CMD_NOMBRE    = CMD_PREFIJO + "nombre"
	CMD_TERMINAR  = CMD_PREFIJO + "terminar"

	CLIENTE_NOMBRE  = "Cliente %v"
	SERVIDOR_NOMBRE = "Servidor"

	ERROR_PREFIJO   = "Error: "
	ERROR_CREAR     = ERROR_PREFIJO + "Ya existe una sala con ese nombre.\n"
	ERROR_ENTRAR    = ERROR_PREFIJO + "No existe una sala con ese nombre.\n"
	ERROR_HISTORIAL = ERROR_PREFIJO + "El vestíbulo no tiene historial.\n"
	ERROR_SALIR     = ERROR_PREFIJO + "No puedes salir del vestíbulo.\n"

	EVENTO_PREFIJO         = "Noticia: "
	EVENTO_ENTRA_SALA      = EVENTO_PREFIJO + "\"%s\" se unió a la sala.\n"
	EVENTO_DEJA_SALA       = EVENTO_PREFIJO + "\"%s\" salió de la sala %v.\n"
	EVENTO_SALA_NOMBRE     = EVENTO_PREFIJO + "\"%s\" cambió su nombre a \"%s\".\n"
	EVENTO_PERSONAL_CREAR  = EVENTO_PREFIJO + "Creaste la sala \"%s\".\n"
	EVENTO_PERSONAL_NOMBRE = EVENTO_PREFIJO + "Cambiaste tu nombre a \"%s\".\n"

	MSJ_CONECTA = "Bienvenido al servidor! Escribe \"/ayuda\" para obtener la lista de comandos.\n"
	MSJ_DESCONECTA = "Desconectado del servidor.\n"
)
