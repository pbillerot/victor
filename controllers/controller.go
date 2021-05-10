package controllers

/**
	MainController
	Variables globales à l'application
**/
import (
	"html/template"
	"path/filepath"
	"time"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/pbillerot/victor/models"
)

var err error

// MainController as
type MainController struct {
	beego.Controller
}

func setSession(c *MainController, name string, value string) {
	c.SetSession(name, value)
	c.Data[name] = value
	c.Data["Ext"] = filepath.Ext(c.Data["File"].(string))
}

// Prepare implements Prepare method for loggedRouter.
func (c *MainController) Prepare() {
	// Session
	if c.GetSession("Hugo") != nil {
		c.Data["Hugo"] = c.GetSession("Hugo").(models.Hugo)
	} else {
		c.Data["Hugo"] = models.Hugo{}
		c.SetSession("Hugo", models.Hugo{})
	}
	if c.GetSession("Folder") != nil {
		c.Data["Folder"] = c.GetSession("Folder").(string)
	} else {
		c.Data["Folder"] = "/victor"
		c.SetSession("Folder", "/")
	}
	if c.GetSession("File") != nil {
		c.Data["File"] = c.GetSession("File").(string)
	} else {
		c.Data["File"] = ""
		c.SetSession("File", "")
	}
	c.Data["Ext"] = filepath.Ext(c.Data["File"].(string))

	// Contexte lié à app.conf
	c.Data["Config"] = models.Config

	// XSRF protection des formulaires
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	// Sera ajouté derrière les urls pour ne pas utiliser le cache des images dynamiques
	c.Data["Composter"] = time.Now().Unix()
	c.Data["Refresh"] = false
}

// int contains in slice
func containsInt(sl []int, in int) bool {
	for _, v := range sl {
		if v == in {
			return true
		}
	}
	return false
}

// string contains in slice
func containsString(sl []string, in string) bool {
	for _, v := range sl {
		if v == in {
			return true
		}
	}
	return false
}
