package controllers

import (
	"strings"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/pbillerot/victor/models"
)

//
//	Fonctions pour les templates
//

// Déclaration des fonctions utilisées dans les templates
func init() {
	beego.AddFuncMap("BeeIncrement", BeeIncrement)
	beego.AddFuncMap("BeeDecrement", BeeDecrement)
	beego.AddFuncMap("BeeReplace", BeeReplace)
	beego.AddFuncMap("BeeSplit", BeeSplit)
	beego.AddFuncMap("BeeSplitBreadcrumb", BeeSplitBreadcrumb)
}

// BeeSplitBreadcrumb /rep1/rep2/rep3/file.ext
func BeeSplitBreadcrumb(path string) (breadcrumb []models.Breadcrumb) {
	reps := strings.Split(path, "/")
	pp := ""
	ll := len(reps)
	isLast := false
	for i, rep := range reps {
		if i == 0 {
			continue
		}
		pp = pp + "/" + rep
		if i == (ll - 1) {
			isLast = true
		}
		breadcrumb = append(breadcrumb, models.Breadcrumb{Base: rep, Path: pp, IsLast: isLast})
	}
	return
}

// BeeSplit strings séparées par un séparateur en slice
func BeeSplit(in string, separateur string) (out []string) {
	if in != "" {
		out = strings.Split(in, separateur)
	} else {
		out = []string{}
	}
	return
}

// BeeReplace as
func BeeReplace(in string, old string, new string) (out string) {
	out = strings.Replace(in, old, new, 1)
	return
}

// BeeIncrement as
func BeeIncrement(in int) (out int) {
	in++
	out = in
	return
}

// BeeDecrement as
func BeeDecrement(in int) (out int) {
	in--
	out = in
	return
}
