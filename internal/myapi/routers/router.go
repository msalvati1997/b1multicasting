package routers

import (
	"b1multicasting/internal/myapi/controllers"
	beego "github.com/beego/beego/v2/server/web"
)

func init() {
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
