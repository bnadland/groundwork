package {{Name}}_test

import (
	"testing"

	"github.com/bnadland/{{Name}}"
	"github.com/go-pg/pg"
)

func NewDatabaseWithTestUser(t *testing.T) *pg.DB {
	db := {{Name}}.NewDatabase(NewTestConfig())
	{{Name}}.ResetDatabase(db)
	if err := {{Name}}.RegisterUser(db, "test", "test"); err != nil {
		t.Fatal(err)
	}
	return db
}

func TestLoginHandler(t *testing.T) {
	s, u, c := NewTestServer(t)
	defer s.Close()

	db := NewDatabaseWithTestUser(t)
	defer db.Close()

	u.Path = "/login"
	t.Log(u)
	t.Log(s.URL)

	u.Path = "/login"
	if err := c.Post(u.String(), map[string]string{
		"Username": "doesnotexist",
		"Password": "test",
	}); err != nil {
		t.Error(err)
	}

	if err := c.Post(u.String(), map[string]string{
		"username": "test",
		"password": "test",
	}); err != nil {
		t.Error(err)
	}

	u.Path = "/login"
	if err := c.Post(u.String(), map[string]string{
		"Username": "test",
		"Password": "test",
	}); err != nil {
		t.Error(err)
	}

	u.Path = "/login"
	if err := c.Post(u.String(), map[string]string{
		"Username": "test",
		"Password": "test",
	}); err != nil {
		t.Error(err)
	}

	c.Wait()
}

func TestLogoutHandler(t *testing.T) {
	s, u, c := NewTestServer(t)
	defer s.Close()
	db := NewDatabaseWithTestUser(t)
	defer db.Close()

	u.Path = "/login"
	if err := c.Post(u.String(), map[string]string{
		"Username": "test",
		"Password": "test",
	}); err != nil {
		t.Error(err)
	}

	u.Path = "/logout"
	if err := c.Post(u.String(), nil); err != nil {
		t.Error(err)
	}

	c.Wait()
}
