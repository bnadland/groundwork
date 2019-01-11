package main

import (
	"log"
	"net/http"

	"github.com/bnadland/{{Name}}"
	"github.com/spf13/cobra"
)

func main() {
	config := {{Name}}.NewConfigFromEnv()
	if config == nil {
		return
	}

	cmd := &cobra.Command{}
	cmd.AddCommand(&cobra.Command{
		Use: "server",
		Run: func(cmd *cobra.Command, args []string) {
			server(config)
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use: "reset",
		Run: func(cmd *cobra.Command, args []string) {
			reset(config)
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use:  "register [username] [password]",
		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			register(config, args[0], args[1])
		},
	})

	cmd.Execute()
}

func register(config *{{Name}}.Config, username string, password string) {
	if err := {{Name}}.RegisterUser({{Name}}.NewDatabase(config), username, password); err != nil {
		log.Print(err)
	}
}

func reset(config *{{Name}}.Config) {
	{{Name}}.ResetDatabase({{Name}}.NewDatabase(config))
}

func server(config *{{Name}}.Config) {
	s := &http.Server{
		Addr:    config.Addr,
		Handler: {{Name}}.NewHandler(config),
	}
	log.Printf("listening on %s", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		log.Print(err)
	}
}
