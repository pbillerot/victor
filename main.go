package main

import (
	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/pbillerot/victor/routers"
)

func main() {
	beego.Run()
}
