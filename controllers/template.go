package controllers

import (
	"strings"

	beego "github.com/beego/beego/v2/server/web"
)

//
//	Fonctions pour les templates
//

// Déclaration des fonctions utilisées dans les templates
func init() {
	beego.AddFuncMap("HugoIncrement", HugoIncrement)
	beego.AddFuncMap("HugoDecrement", HugoDecrement)
	beego.AddFuncMap("BeeReplace", BeeReplace)
	beego.AddFuncMap("BeeSplit", BeeSplit)
}

// BeeSplit BeeSplit strings séparées par une virgule en slice
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

// HugoIncrement as
func HugoIncrement(in int) (out int) {
	in++
	out = in
	return
}

// HugoDecrement as
func HugoDecrement(in int) (out int) {
	in--
	out = in
	return
}
