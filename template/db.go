package {{Name}}

import (
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

var allModels = []interface{}{
	&User{},
}

func NewDatabase(config *Config) *pg.DB {
	db := pg.Connect(&pg.Options{
		User:     config.DBUser,
		Database: config.DBName,
	})
	orm.SetTableNameInflector(func(s string) string { return s })
	db.AddQueryHook(Logger{})
	return db
}

func ResetDatabase(db *pg.DB) {
	for _, m := range allModels {
		if err := db.DropTable(m, &orm.DropTableOptions{
			IfExists: true,
		}); err != nil {
			log.Print(err)
		}
	}
	for _, m := range allModels {
		if err := db.CreateTable(m, &orm.CreateTableOptions{
			IfNotExists: true,
		}); err != nil {
			log.Print(err)
		}
	}
}
