package main

import (
	"github.com/astaxie/beego/logs"
	"github.com/beego/beego/v2/core/config"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/pbillerot/victor/controllers"
	_ "github.com/pbillerot/victor/routers"
)

func init() {

	// Initialisation de controllers.AppConfig
	if val, ok := config.String("hugo_racine"); ok == nil {
		controllers.AppConfig.HugoRacine = val
	}
	if val, ok := config.String("hugo_deploy"); ok == nil {
		controllers.AppConfig.HugoDeploy = val
	}
	if val, ok := config.String("title"); ok == nil {
		controllers.AppConfig.Title = val
	}
	if val, ok := config.String("description"); ok == nil {
		controllers.AppConfig.Description = val
	}
	if val, ok := config.String("favicon"); ok == nil {
		controllers.AppConfig.Favicon = val
	}
	if val, ok := config.String("icon"); ok == nil {
		controllers.AppConfig.Icon = val
	}
	logs.Info("AppConfig", controllers.AppConfig)
}

func main() {
	beego.Run()
}
