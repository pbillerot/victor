package main

import (
	"github.com/astaxie/beego/logs"
	"github.com/beego/beego/v2/core/config"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/pbillerot/victor/models"
	_ "github.com/pbillerot/victor/routers"
)

func main() {
	beego.Run()
}

func init() {
	// Initialisation de models.Config
	if val, ok := config.String("hugo_dir"); ok == nil {
		models.Config.HugoRacine = val
		models.Config.HugoContent = val + "/content"
		models.Config.HugoPrivate = val + "/private"
		models.Config.HugoPublic = val + "/public"
	}
	if val, ok := config.String("test_url"); ok == nil {
		models.Config.HugoPrivateURL = val
	}
	if val, ok := config.String("pub_url"); ok == nil {
		models.Config.HugoPublicURL = val
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

	// Enregistrement de content an tant que répertoire static
	beego.SetStaticPath("/content", models.Config.HugoContent)

}
