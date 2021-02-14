package evil

import (
	"fmt"
	"html/template"
	"io"
	"math/rand"
	"strings"

	"github.com/chrisseto/evil/channel"
)

func ID() string {
	alphabet := "abcdefghijklmnopqrstuvwxyz"

	var id strings.Builder

	id.WriteString("phx-")

	for i := 0; i < 5; i++ {
		id.WriteByte(alphabet[rand.Intn(len(alphabet))])
	}

	return id.String()
}

type Renderer struct {
	SessionFactory *SessionFactory
	pages          *template.Template
	views          *template.Template
	components     *template.Template
	_views         map[string]View
}

func NewRenderer(
	pagesGlob string,
	viewsGlob string,
	componentGlob string,
) (*Renderer, error) {
	t := template.New("").Funcs(template.FuncMap{
		"EvilView":      func() string { return "" },
		"EvilComponent": func() string { return "" },
	})

	pages, err := template.Must(t.Clone()).ParseGlob(pagesGlob)
	if err != nil {
		return nil, err
	}

	views, err := template.Must(t.Clone()).ParseGlob(viewsGlob)
	if err != nil {
		return nil, err
	}

	components, err := template.Must(t.Clone()).ParseGlob(componentGlob)
	if err != nil {
		return nil, err
	}

	return &Renderer{
		pages:      pages,
		views:      views,
		components: components,
		_views:     map[string]View{},
	}, nil
}

func (r *Renderer) RegisterView(name string, view View) error {
	r._views[name] = view
	return nil
}

func (r *Renderer) RenderPage(wr io.Writer, page string) error {
	t := template.Must(template.Must(r.pages.Clone()).AddParseTree("", r.components.Tree))

	t.Funcs(template.FuncMap{
		"EvilView": func(name string, data interface{}) (template.HTML, error) {
			var b strings.Builder

			s, err := r.SessionFactory.NewSession(name)
			if err != nil {
				return "", err
			}

			// data-phx-main ??
			if _, err := fmt.Fprintf(
				&b,
				`<%s id="%s" data-phx-view="%s" data-phx-session="%s" data-phx-static="%s">`,
				"div",
				s.ID,
				s.View,
				s.Encode(),
				s.Encode(), // TODO this is incorrect?
			); err != nil {
				return "", err
			}

			if _, err := b.WriteRune('\n'); err != nil {
				return "", err
			}

			if err := r.views.ExecuteTemplate(&b, name, data); err != nil {
				return "", err
			}

			if _, err := b.WriteRune('\n'); err != nil {
				return "", err
			}

			if _, err := fmt.Fprintf(&b, `</%s>`, "div"); err != nil {
				return "", nil
			}

			return template.HTML(b.String()), nil
		},
	})

	return t.ExecuteTemplate(wr, page, nil)
}

func (r *Renderer) Mount(view string, s *Session) error {
	return r._views[view].OnMount(s)
}

func (r *Renderer) Event(view string, s *Session, e *channel.Event) error {
	return r._views[view].HandleEvent(s, e)
}

func (r *Renderer) RenderView(view string, s *Session) (string, error) {
	var b strings.Builder

	args, err := r._views[view].ToArgs(s)
	if err != nil {
		return "", err
	}

	if err := r.views.ExecuteTemplate(&b, view, args); err != nil {
		return "", err
	}

	return b.String(), nil
}
