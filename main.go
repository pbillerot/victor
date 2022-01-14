package main

import (
	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/pbillerot/victor/models"
	_ "github.com/pbillerot/victor/routers"
)

var err error

func main() {
	web.Run()
}

func init() {
	// Initialisation de models.Config
	if val, ok := config.String("github"); ok == nil {
		models.Config.Github = val
	}
	if val, ok := config.String("help"); ok == nil {
		models.Config.Help = val
	}
	// if val, ok := config.String("helpEditor"); ok == nil {
	// 	models.Config.HelpEditor = val
	// }
	if val, ok := config.String("version"); ok == nil {
		models.Config.Version = val
	}
	if val, ok := config.String("appname"); ok == nil {
		models.Config.Appname = val
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

	// Récupération de l'aide en ligne
	// src := models.Config.Help
	// dst := "help"
	// _, err = os.Open(src + "/index.html")
	// if !os.IsNotExist(err) {
	// 	err = shutil.CopyTree(src, dst, nil)
	// 	if err != nil {
	// 		msg := fmt.Sprintf("Copie [%s] vers [%s] : %v", src, dst, err)
	// 		logs.Error(msg)
	// 	}
	// }
	// web.SetStaticPath("/help", "help")
	initConfigHugo()
}

func initConfigHugo() {
	// Lecture hugo.yaml -> models.HugoApps.Apps
	err = models.LoadHugoApps()
	if err != nil {
		logs.Error("LoadHugoApps", err)
		return
	}
	models.Config.HugoApps = models.HugoApps.Apps

	// Déclaration des url de production à servir
	logs.Info("HugoApp: liste des webapp à servir")
	for _, hugoApp := range models.Config.HugoApps {
		web.SetStaticPath(hugoApp.BaseURL, hugoApp.Folder+"/public")
		logs.Info("HugoApp:", hugoApp.Title, hugoApp.Folder)
	}

}
