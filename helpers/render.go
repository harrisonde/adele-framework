package helpers

import (
	"net/http"

	"github.com/CloudyKit/jet/v6"
)

func (h *Helpers) Render(w http.ResponseWriter, r *http.Request, template string, variables, data interface{}) error {
	vars := make(jet.VarMap)
	if variables == nil {
		vars = make(jet.VarMap)
	} else {
		vars = variables.(jet.VarMap)
	}

	vars.Set("view", template)
	vars.Set("path", r.URL.Path)

	return h.Redner.Page(w, r, template, vars, data)
}
