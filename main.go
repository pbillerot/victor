package main

import (
	"fmt"
	"os"

	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/pbillerot/victor/controllers"
	"github.com/pbillerot/victor/models"
	_ "github.com/pbillerot/victor/routers"
	"github.com/pbillerot/victor/shutil"
	"github.com/spf13/viper"
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
	src := models.Config.Help
	dst := "help"
	_, err = os.Open(src + "/index.html")
	if !os.IsNotExist(err) {
		err = shutil.CopyTree(src, dst, nil)
		if err != nil {
			msg := fmt.Sprintf("Copie [%s] vers [%s] : %v", src, dst, err)
			logs.Error(msg)
		}
	}
	web.SetStaticPath("/help", "help")
	initConfigHugo()
}

func initConfigHugo() {
	// Lecture hugo.yaml -> models.HugoApps.Apps
	models.LoadHugoApps()
	models.Config.HugoApps = models.HugoApps.Apps

	// Init Viper
	viper.SetConfigName("ctx")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("conf/")
	viper.ReadInConfig()

	// Déclaration des url de production à servir
	for _, hugoApp := range models.Config.HugoApps {
		web.SetStaticPath(hugoApp.BaseURL, hugoApp.Folder+"/public")
	}

	// Positionnement sur la dernière hugoapp utilisée
	if viper.GetString("hugoapp") == "" {
		hugoApp := models.GetFirstHugoApp()
		controllers.SetHugoApp(hugoApp)
	} else {
		hugoApp := models.GetHugoApp(viper.GetString("hugoapp"))
		controllers.SetHugoApp(hugoApp)
	}

}
