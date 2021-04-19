package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
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
		if strings.HasPrefix(val, "/") {
			models.Config.HugoContentDir = val + "/content"
			models.Config.HugoPrivateDir = val + "/private"
			models.Config.HugoPublicDir = val + "/public"
		} else {
			models.Config.HugoContentDir = val + "/content"
			models.Config.HugoPrivateDir = "private"
			models.Config.HugoPublicDir = "public"
		}
	}
	if val, ok := config.String("github"); ok == nil {
		models.Config.Github = val
	}
	if val, ok := config.String("help"); ok == nil {
		models.Config.Help = val
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
	if strings.HasPrefix(models.Config.HugoRacine, "/") {
		web.SetStaticPath("/content", models.Config.HugoContentDir)
		web.SetStaticPath("/hugo", models.Config.HugoPrivateDir)
		web.SetStaticPath("/", models.Config.HugoPublicDir)
	} else {
		// path relatif à la webapp victor
		web.SetStaticPath("/content", models.Config.HugoContentDir)
		web.SetStaticPath("/hugo", models.Config.HugoRacine+"/"+models.Config.HugoPrivateDir)
		web.SetStaticPath("/", models.Config.HugoRacine+"/"+models.Config.HugoPublicDir)
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
	// mkdir /config/hugo
	err = shutil.CreateDir(models.Config.HugoRacine + "/config/hugo")
	if err != nil {
		logs.Error("initConfigHugo", err)
		return
	}
	// copy config.hugo.yaml
	var configHugoSrc = fmt.Sprintf("./conf/config.hugo.yaml")
	var configHugoDst = fmt.Sprintf("%s/config/hugo/config.yaml", models.Config.HugoRacine)
	_, err = os.Open(configHugoSrc)
	if !os.IsNotExist(err) {
		_, err = os.Open(configHugoDst)
		if os.IsNotExist(err) {
			err = shutil.CopyFile(configHugoSrc, configHugoDst, false)
			if err != nil {
				msg := fmt.Sprintf("Copie [%s] vers [%s] : %v", configHugoSrc, configHugoDst, err)
				logs.Error(msg)
				return
			}
		}
	}

}
