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
	wom.ConfigurePocketBase(app)
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
