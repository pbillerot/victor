package main

import (
	"github.com/astaxie/beego/logs"
	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/server/web"
	"github.com/pbillerot/victor/models"
	_ "github.com/pbillerot/victor/routers"
)

func main() {
	web.Run()
}

func init() {
	// Initialisation de models.Config
	if val, ok := config.String("hugo_dir"); ok == nil {
		models.Config.HugoRacine = val
		models.Config.HugoContentDir = val + "/content"
		models.Config.HugoPrivateDir = val + "/private"
		models.Config.HugoPublicDir = val + "/public"
	}
	if val, ok := config.String("title"); ok == nil {
		models.Config.Title = val
	}
	if val, ok := config.String("description"); ok == nil {
		models.Config.Description = val
	}
	if val, ok := config.String("favicon"); ok == nil {
		models.Config.Favicon = val
	}
	if val, ok := config.String("icon"); ok == nil {
		models.Config.Icon = val
	}
	logs.Info("Config", models.Config)

	// Répertoires statiques
	web.SetStaticPath("/content", models.Config.HugoContentDir)
	web.SetStaticPath("/hugo", models.Config.HugoPrivateDir)
	web.SetStaticPath("/", models.Config.HugoPublicDir)

}
