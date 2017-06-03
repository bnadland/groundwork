//go:generate esc -prefix public -o public.go public
//go:generate go-bindata -o templates.go templates/...
package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func runServer(listen string, h http.Handler) {
	log.WithField("listen", listen).Print("listening")
	if err := http.ListenAndServe(listen, h); err != nil {
		log.WithError(err).Print("could not listen")
	}
}

func main() {
	c, err := newCache()
	if err != nil {
		log.WithError(err).Print("could not initialize cache")
		return
	} else {
		go c.Run()
	}

	db, err := gorm.Open("postgres", config.Database)
	if err != nil {
		log.WithError(err).Print("could not connect to database")
		return
	}
	defer db.Close()

	r := newRender()

	m := mux.NewRouter()
	m.Handle("/metrics", promhttp.Handler())
	m.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if err := r.HTML(w, 200, "pages/homepage", nil); err != nil {
			log.WithError(err).Print("could not render homepage")
			r.Text(w, 500, "internal error")
		}
	})

	runServer(config.Listen, newMiddlewares(m))
}
