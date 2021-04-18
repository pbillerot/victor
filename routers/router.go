package routers

import (
	"github.com/beego/beego/v2/server/web"
	"github.com/pbillerot/victor/controllers"
)

func init() {

	// web.Router("/hugo/", &controllers.MainController{}, "get:Hugo")

	web.Router("/victor/", &controllers.MainController{}, "get:Main")
	web.Router("/victor/folder/", &controllers.MainController{}, "get:Folder")
	web.Router("/victor/folder/*.*", &controllers.MainController{}, "get:Folder")
	web.Router("/victor/document/*.*", &controllers.MainController{}, "get:Document;post:Document")
	web.Router("/victor/image/*.*", &controllers.MainController{}, "get:Image;post:Image")
	web.Router("/victor/pdf/*.*", &controllers.MainController{}, "get:Pdf")
	// gestion fichiers
	web.Router("/victor/new", &controllers.MainController{}, "post:FileNew")
	web.Router("/victor/mkdir", &controllers.MainController{}, "post:FileMkdir")
	web.Router("/victor/rn/*.*", &controllers.MainController{}, "post:FileRename")
	// sélection multiple
	web.Router("/victor/cp", &controllers.MainController{}, "post:FileCp")
	web.Router("/victor/mv", &controllers.MainController{}, "post:FileMove")
	web.Router("/victor/rm", &controllers.MainController{}, "post:FileRm")
	web.Router("/victor/upload", &controllers.MainController{}, "post:FileUpload")

	web.Router("/victor/action/:action", &controllers.MainController{}, "post:Action")
	web.Router("/victor/api/folders", &controllers.MainController{}, "get:APIFolders")
	web.Router("/victor/api/file", &controllers.MainController{}, "get:APIFile")

}
