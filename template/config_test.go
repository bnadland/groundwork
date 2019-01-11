package {{Name}}_test

import (
	"fmt"

	"github.com/bnadland/{{Name}}"
)

func NewTestConfig() *{{Name}}.Config {
	config := {{Name}}.NewConfigFromEnv()
	config.DBName = fmt.Sprintf("%s_test", config.DBName)
	return config
}
