package {{Name}}

import (
	"net/http"

	"github.com/unrolled/render"
)

func HomeHandler(r *render.Render) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if err := r.HTML(w, 200, "page/index", nil); err != nil {
			log.Print(err)
		}
	})
}
