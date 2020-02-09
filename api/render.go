package api

import (
	"net/http"

	"github.com/go-chi/render"
)

// Render is a local render to make things a bit more easy
// it always looks for a status code and sets it before calling render
func Render(w http.ResponseWriter, r *http.Request, v *Response) {
	w.WriteHeader(v.Status.Code)
	render.Render(w, r, v)
}
