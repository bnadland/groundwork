package {{Name}}

import (
	"net/http"
	"time"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/gorilla/sessions"
	"github.com/hlandau/passlib"
	"github.com/unrolled/render"
)

type User struct {
	ID        int
	CreatedAt time.Time
	UpdatedAt time.Time

	Username string `sql:",unique"`
	Password string
}

func (u *User) BeforeInsert(db orm.DB) error {
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return nil
}

func (u *User) BeforeUpdate(db orm.DB) error {
	u.UpdatedAt = time.Now()
	return nil
}

func RegisterUser(db *pg.DB, username string, password string) error {
	//	passlib.UseDefaults("latest")
	passwordHash, err := passlib.Hash(password)
	if err != nil {
		log.Print(err)
		return err
	}
	user := User{
		Username: username,
		Password: passwordHash,
	}
	if err := db.Insert(&user); err != nil {
		log.Print(err)
		return err
	}
	return nil
}

func GetUserByUsernameAndPassword(db *pg.DB, username string, password string) *User {
	passlib.UseDefaults("latest")
	var user User
	if err := db.Model(&user).Where("username = ?", username).Select(); err != nil {
		log.Print(err)
		return nil
	}
	if newHash, err := passlib.Verify(password, user.Password); err != nil {
		log.Print(err)
		return nil
	} else {
		if newHash != "" {
			user.Password = newHash
			if err := db.Update(&user); err != nil {
				log.Print(err)
			}
		}
	}
	return &user
}

func GetUserFromSession(db *pg.DB, session *sessions.Session) *User {
	if session.Values["userid"] == nil {
		return nil
	}
	var user User
	if err := db.Model(&user).Where("id = ?", session.Values["userid"]).Select(); err != nil {
		log.Print(err)
		session.Values["userid"] = nil
		return nil
	}
	return &user
}

func LoginHandler(session Session, r *render.Render, db *pg.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		session := session(req)
		if user := GetUserFromSession(db, session); user != nil {
			if err := session.Save(req, w); err != nil {
				log.Print(err)
			}
			http.Redirect(w, req, "/", 303)
			return
		}

		if req.Method == "GET" {
			flashes := session.Flashes()
			if err := session.Save(req, w); err != nil {
				log.Print(err)
			}
			r.HTML(w, 200, "page/login", flashes)
			return
		}

		if req.Method == "POST" {
			if err := req.ParseForm(); err != nil {
				log.Print(err)
				http.Error(w, "internal error", 500)
				return
			}

			username := req.Form.Get("Username")
			password := req.Form.Get("Password")
			if username == "" {
				session.AddFlash("Username is required")
			}
			if password == "" {
				session.AddFlash("Password is required")
			}

			if username != "" && password != "" {
				user := GetUserByUsernameAndPassword(db, username, password)
				if user != nil {
					session.Values["userid"] = user.ID
					if err := session.Save(req, w); err != nil {
						log.Print(err)
					}
					http.Redirect(w, req, "/", 303)
					return
				}
				session.AddFlash("We have found no account with that username and password")
			}
			if err := session.Save(req, w); err != nil {
				log.Print(err)
			}
			http.Redirect(w, req, "/login", 303)
			return
		}
	})
}

func LogoutHandler(session Session) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "POST" {
			session := session(req)
			session.Values["userid"] = nil
			if err := session.Save(req, w); err != nil {
				log.Print(err)
			}
			http.Redirect(w, req, "/", 303)
			return
		}
	})
}
