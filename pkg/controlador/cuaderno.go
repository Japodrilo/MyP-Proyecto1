package controlador

import(
	"github.com/Japodrilo/MyP-Proyecto1/pkg/vista"

	"github.com/gotk3/gotk3/gtk"
)

type Cuaderno struct {
	nb 			*gtk.Notebook
	entradas    map[string]*gtk.Entry
	textos		map[string]*gtk.TextBuffer
	botones		map[string]*gtk.Button
	tabs		map[string]int
}

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


func (cuaderno *Cuaderno) EntradaUsuario(usuario string) *gtk.Entry {
	return cuaderno.entradas[usuario]
}

func (cuaderno *Cuaderno) TextoUsuario(usuario string) *gtk.TextBuffer {
	return cuaderno.textos[usuario]
}