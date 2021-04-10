package main

import (
	"fmt"
	"os"
	"strings"

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
	if val, ok := config.String("HelpEditor"); ok == nil {
		models.Config.HelpEditor = val
	}
	if val, ok := config.String("github"); ok == nil {
		models.Config.Github = val
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
	web.SetStaticPath("/help", "site/public/")

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
	// config.hugo.yaml
	var configHugoSrc = fmt.Sprintf("%s%s%s/config.hugo.yaml",
		models.Config.HugoRacine,
		"/themes/",
		models.Config.HugoTheme,
	)
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

	// --> /content/site
	err = shutil.CreateDir(models.Config.HugoRacine + "/content/site")
	if err != nil {
		logs.Error("initConfigHugo", err)
		return
	}

	// config.template.yaml
	configHugoSrc = fmt.Sprintf("%s%s%s/config.template.yaml",
		models.Config.HugoRacine,
		"/themes/",
		models.Config.HugoTheme,
	)
	// --> /config/_default/config.yaml
	configHugoDst = fmt.Sprintf("%s/config/_default/config.yaml", models.Config.HugoRacine)
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
		// configHugoDst = fmt.Sprintf("%s/content/site/config.yaml", models.Config.HugoRacine)
		// _, err = os.Open(configHugoDst)
		// if os.IsNotExist(err) {
		// 	err = shutil.CopyFile(configHugoSrc, configHugoDst, false)
		// 	if err != nil {
		// 		msg := fmt.Sprintf("Copie [%s] vers [%s] : %v", configHugoSrc, configHugoDst, err)
		// 		logs.Error(msg)
		// 		return
		// 	}
		// }
	}

	// remplacer le modèle par défaut par celui du thème
	configHugoSrc = fmt.Sprintf("%s%s%s/archetypes/default.md",
		models.Config.HugoRacine,
		"/themes/",
		models.Config.HugoTheme,
	)
	configHugoDst = fmt.Sprintf("%s/archetypes/default.md", models.Config.HugoRacine)
	_, err = os.Open(configHugoSrc)
	if !os.IsNotExist(err) {
		err = shutil.CopyFile(configHugoSrc, configHugoDst, false)
		if err != nil {
			msg := fmt.Sprintf("Copie [%s] vers [%s] : %v", configHugoSrc, configHugoDst, err)
			logs.Error(msg)
			return
		}
	}

	// Si theme beedream
	if models.Config.HugoTheme == "beedream" {
		// renommer config.toml et le remplacer par config.root.yaml
		configHugoDst = fmt.Sprintf("%s/config.toml", models.Config.HugoRacine)
		_, err = os.Open(configHugoDst)
		if !os.IsNotExist(err) {
			err = os.Rename(configHugoDst, strings.ReplaceAll(configHugoDst, "config.toml", "config.toml.old"))
			if err != nil {
				msg := fmt.Sprintf("Rename [%s] : %v", configHugoDst, err)
				logs.Error(msg)
				return
			}
		}
		// pour remplacer par config.root.yaml
		configHugoSrc = fmt.Sprintf("%s%s%s/config.root.yaml",
			models.Config.HugoRacine,
			"/themes/",
			models.Config.HugoTheme,
		)
		// --> /config.yaml
		configHugoDst = fmt.Sprintf("%s/config.yaml", models.Config.HugoRacine)
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
