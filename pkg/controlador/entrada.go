package controlador

import(
	"strings"
	
	//"github.com/Japodrilo/MyP-Proyecto1/pkg/vista"

	//"github.com/gotk3/gotk3/gtk"
)

func parse(mensaje string) []string {
	r := make ([]string, 0)
	switch {
		case strings.HasPrefix(mensaje, "...PUBLIC-"):
			r[0] = "...PUBLIC-"
			r = append(r, strings.Fields(strings.TrimPrefix(mensaje, "...PUBLIC-"))...)
	}
	return r
}

func ManejaEntrada(mensaje string) {
	entrada := parse(mensaje)

	switch entrada[0] {
	case "...PUBLIC-":

	default:
		return
	}
}