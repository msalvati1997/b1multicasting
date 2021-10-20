package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/msalvati1997/b1multicasting/controllers"
	"log"
)

func init() {
	log.Println("Init root")
	beego.Router("/", &controllers.RegistryController{})
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/registry",
			beego.NSInclude(
				&controllers.RegistryController{},
			),
		),
		beego.NSNamespace("/multicast",
			beego.NSInclude(
				&controllers.MulticastController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
