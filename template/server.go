package {{Name}}

import (
	"net/http"

	"github.com/gobuffalo/packr/v2"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

type Session func(*http.Request) *sessions.Session

func NewSession(config *Config) Session {
	session := sessions.NewCookieStore([]byte(config.SessionKey))
	return func(req *http.Request) *sessions.Session {
		s, err := session.Get(req, config.SessionName)
		if err != nil {
			log.Print(err)
		}
		return s
	}
}

func NewHandler(config *Config) http.Handler {
	db := NewDatabase(config)
	r := NewRenderer(config)
	session := NewSession(config)

	m := mux.NewRouter()
	m.Use(Logger{}.Middleware)
	m.Handle("/assets/{slug:.*}", http.StripPrefix("/assets/", http.FileServer(packr.New("public", "./public"))))

	m.Handle("/", HomeHandler(r))
	m.Handle("/login", LoginHandler(session, r, db))
	m.Handle("/logout", LogoutHandler(session))

	return m
}
