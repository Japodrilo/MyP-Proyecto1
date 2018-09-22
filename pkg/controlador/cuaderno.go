package controlador

import(
	"github.com/Japodrilo/MyP-Proyecto1/pkg/vista"

	"github.com/gotk3/gotk3/gtk"
)

/**
 * Estructura que modela las pestañas con conversaciones.
 * Contiene diccionarios que brindan fácil acceso a las
 * entradas de texto, visores de texto y pestañas para cada
 * conversación.   Además, contiene un diccionario de
 * botones correspondiente a los usuarios conectados.
 */
type Cuaderno struct {
	nb 			*gtk.Notebook
	entradas    map[string]*gtk.Entry
	textos		map[string]*gtk.TextBuffer
	botones		map[string]*gtk.Button
	tabs		map[string]int
}

/**
 * Constructor que recibe como parámetro un apuntador a
 * un NoteBook de gtk.
 */
func NuevoCuaderno (nb *gtk.Notebook) *Cuaderno {
	textos := make(map[string]*gtk.TextBuffer)
	entradas := make(map[string]*gtk.Entry)
	botones := make(map[string]*gtk.Button)
	tabs := make(map[string]int)
	
	return &Cuaderno{
		nb:			nb,
		entradas: 	entradas,
		textos:		textos,
		botones:	botones,
		tabs:		tabs,
	}
}

/**
 * Función que añade una pestaña al NoteBook del cuaderno.
 * Recibe como parámetro la cadena para ponerle nombre a la
 * pestaña.
 */
func (cuaderno *Cuaderno) AddTab(name string) (*gtk.Entry, *gtk.TextBuffer) {
	box := vista.SetupBox()
	entry := vista.SetupEntry()
	scrwin := vista.SetupScrolledWindow()
	tv := vista.SetupTextView()
	nbTab := vista.SetupLabel(name)
	
	tv.SetVExpand(true)

	scrwin.Add(tv)
	box.Add(entry)
	box.Add(scrwin)

	cuaderno.nb.SetHExpand(true)
	cuaderno.nb.SetVExpand(true)

	cuaderno.nb.AppendPage(box, nbTab)
	cuaderno.tabs[name] = cuaderno.nb.GetCurrentPage()
	cuaderno.nb.Connect("change-current-page", func() {
		entry.GrabFocus()
	})

	return entry, vista.GetBufferTV(tv)
}