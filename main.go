package main

import (
	"fmt"

	"github.com/astaxie/beego/logs"
	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/server/web"
	"github.com/pbillerot/victor/models"
	_ "github.com/pbillerot/victor/routers"
	"github.com/pbillerot/victor/shutil"
)

var err error

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
	if val, ok := config.String("HelpEditor"); ok == nil {
		models.Config.HelpEditor = val
	}
	if val, ok := config.String("hugo_theme"); ok == nil {
		models.Config.HugoTheme = val
	}
	if val, ok := config.String("version"); ok == nil {
		models.Config.Version = val
	}
	if val, ok := config.String("appname"); ok == nil {
		models.Config.Appname = val
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
	if val, ok := config.String("background"); ok == nil {
		models.Config.Background = val
	}
	if val, ok := config.String("icon"); ok == nil {
		models.Config.Icon = val
	}
	logs.Info("Config", models.Config)

	// Répertoires statiques
	web.SetStaticPath("/content", models.Config.HugoContentDir)
	web.SetStaticPath("/hugo", models.Config.HugoPrivateDir)
	web.SetStaticPath("/", models.Config.HugoPublicDir)

	initConfigHugo()
}

func initConfigHugo() {
	// Tout d'abord recopie /themes/beedream/config.hugo.yaml dans /config/hugo/config.yaml
	err = shutil.CreateDir(models.Config.HugoRacine + "/config")
	if err != nil {
		logs.Error("initConfigHugo", err)
		return
	}
	err = shutil.CreateDir(models.Config.HugoRacine + "/config/_default")
	if err != nil {
		logs.Error("initConfigHugo", err)
		return
	}
	err = shutil.CreateDir(models.Config.HugoRacine + "/config/hugo")
	if err != nil {
		logs.Error("initConfigHugo", err)
		return
	}
	var configHugoSrc = fmt.Sprintf("%s%s%s/config.hugo.yaml",
		models.Config.HugoRacine,
		"/themes/",
		models.Config.HugoTheme,
	)
	var configHugoDst = fmt.Sprintf("%s/config/hugo/config.yaml", models.Config.HugoRacine)

	err = shutil.CopyFile(configHugoSrc, configHugoDst, false)
	if err != nil {
		msg := fmt.Sprintf("Copie [%s] vers [%s] : %v", configHugoSrc, configHugoDst, err)
		logs.Error(msg)
		return
	}

}
