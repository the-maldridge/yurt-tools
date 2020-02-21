package web

import (
	"html/template"
	"io"
	"log"
	"path/filepath"
	"strings"

	"github.com/labstack/echo"
)

// Renderer fulfills the echo Renderer interface and maintains a cache
// of templates.  It also provides a debug mode where changes to
// templates on disk are reflected on the next page reload rather than
// having to restart the server.
type renderer struct {
	tmplDir string

	tmpls map[string]*template.Template
}

func newRenderer(baseDir string) (*renderer, error) {
	r := &renderer{}

	r.tmplDir = baseDir

	if err := r.loadTmpls(); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *renderer) loadTmpls() error {
	base, err := template.ParseGlob(filepath.Join(r.tmplDir, "base", "*.tpl"))
	if err != nil {
		log.Printf("Error loading base template: %v", err)
		return err
	}

	r.tmpls = make(map[string]*template.Template)

	// We can safely throw away this error because the pattern is
	// hard coded.
	pages, _ := filepath.Glob(filepath.Join(r.tmplDir, "pages", "*", "*.tpl"))

	for _, path := range pages {
		page := filepath.Base(filepath.Dir(path))
		pTmpl, err := base.Clone()
		if err != nil {
			return err
		}
		pTmpl, err = pTmpl.ParseGlob(path)
		if err != nil {
			log.Printf("Error parsing page template %s: %v", page, err)
			return err
		}
		r.tmpls[page] = pTmpl
	}

	return nil
}

func (r *renderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	r.loadTmpls()

	d := echo.Map{}
	d["data"] = data
	d["title"] = strings.Title(strings.ReplaceAll(name, "-", " "))
	return r.tmpls[name].ExecuteTemplate(w, "base", d)
}
