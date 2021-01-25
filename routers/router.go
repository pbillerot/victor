package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/pbillerot/victor/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{}, "get:Main")

	beego.Router("/folder/*.*", &controllers.MainController{}, "get:Folder")
	beego.Router("/document/*.*", &controllers.MainController{}, "get:Document;post:Document")
	beego.Router("/image/*.*", &controllers.MainController{}, "get:Image;post:Image")
	beego.Router("/pdf/*.*", &controllers.MainController{}, "get:Pdf")
	// beego.Router("/directory/*.*", &controllers.MainController{}, "get:Directory")
	// beego.Router("/file/*.*", &controllers.MainController{}, "get:File")
	beego.Router("/rn/*.*", &controllers.MainController{}, "post:FileRename")
	beego.Router("/mv/*.*", &controllers.MainController{}, "post:FileMove")
	beego.Router("/new/*.*", &controllers.MainController{}, "post:FileNew")
	beego.Router("/cp/*.*", &controllers.MainController{}, "post:FileCp")
	beego.Router("/rm/*.*", &controllers.MainController{}, "post:FileRm")
	beego.Router("/mkdir/*.*", &controllers.MainController{}, "post:FileMkdir")
	beego.Router("/upload/*.*", &controllers.MainController{}, "post:FileUpload")
	beego.Router("/action/:action", &controllers.MainController{}, "post:Action")

	beego.Router("/api/folders", &controllers.MainController{}, "get:APIFolders")
}
