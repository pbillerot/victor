package controllers

/**
	MainController
	Variables globales à l'application
**/
import (
	"html/template"
	"time"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/pbillerot/victor/models"
)

var err error

// MainController as
type MainController struct {
	beego.Controller
}

// Prepare implements Prepare method for loggedRouter.
func (c *MainController) Prepare() {

	// Contexte lié à app.conf
	c.Data["Config"] = models.Config

	// XSRF protection des formulaires
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	// Sera ajouté derrière les urls pour ne pas utiliser le cache des images dynamiques
	c.Data["Composter"] = time.Now().Unix()
	// Contexte de navigation
	c.Data["VictorFolder"] = c.Ctx.Input.Cookie("victor-folder")
	c.Data["VictorFile"] = c.Ctx.Input.Cookie("victor-file")

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
