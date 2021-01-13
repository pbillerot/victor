package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/pbillerot/victor/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{}, "get:Main")

	beego.Router("/folder/*.*", &controllers.MainController{}, "get:Folder")
	beego.Router("/document/:key", &controllers.MainController{}, "get:Document;post:Document")
	beego.Router("/image/:key", &controllers.MainController{}, "get:Image;post:Image")
	beego.Router("/pdf/:key", &controllers.MainController{}, "get:Pdf")
	beego.Router("/directory/:key", &controllers.MainController{}, "get:Directory")
	beego.Router("/file/:key", &controllers.MainController{}, "get:File")
	beego.Router("/mv/:key", &controllers.MainController{}, "post:FileMv")
	beego.Router("/new/:key", &controllers.MainController{}, "post:FileNew")
	beego.Router("/cp/:key", &controllers.MainController{}, "post:FileCp")
	beego.Router("/rm/:key", &controllers.MainController{}, "post:FileRm")
	beego.Router("/mkdir/:key", &controllers.MainController{}, "post:FileMkdir")
	beego.Router("/upload/:key", &controllers.MainController{}, "post:FileUpload")
	beego.Router("/action/:action", &controllers.MainController{}, "post:Action")

}
