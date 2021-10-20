package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/msalvati1997/b1multicasting/internal/myapi/controllers"
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
