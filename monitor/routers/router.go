package routers

import (
	"github.com/oikomi/FishChatServer/monitor/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
}
