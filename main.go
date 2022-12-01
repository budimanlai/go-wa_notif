package main

import (
	services "github.com/budimanlai/go-cli-service"
	"github.com/eqto/dbm"
)

func main() {
	srv := services.NewService("config/main.conf")
	srv.AppName = "WhatsApp Send Message"
	srv.Version = "0.0.3"

	cn, e1 := dbm.Connect("mysql", srv.Config.GetString("iam.hostname"), srv.Config.GetInt("iam.port"),
		srv.Config.GetString("iam.username"), srv.Config.GetString("iam.password"), srv.Config.GetString("iam.database"))
	if e1 != nil {
		panic(e1)
	}
	srv.SetDatabase(cn)

	srv.StartHandler(StartService)
	srv.StopHandler(StopService)
	e := srv.Start()
	if e != nil {
		panic(e)
	}
}
