package main

import (
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	"github.com/puzzad/wom"
	_ "github.com/puzzad/wom/migrations"
	"log"
)

func main() {
	app := pocketbase.NewWithConfig(&pocketbase.Config{
		DefaultDataDir: "./data",
	})
	migratecmd.MustRegister(app, app.RootCmd, &migratecmd.Options{
		Automigrate: true,
	})
	if err := wom.ConfigurePocketBase(app); err != nil {
		log.Fatal(err)
	}
	if err := app.Execute(); err != nil {
		log.Fatal(err)
	}
}
