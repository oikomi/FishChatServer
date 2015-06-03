package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/oikomi/FishChatServer/monitor/conf"
)

type MonitorController struct {
	beego.Controller
}

func (this *MonitorController) Get() {
	beego.Info("MonitorController Get")
	action := this.GetString(conf.KEY_ACTION)
	if action == "" {
		beego.Error("[para is null] | action ")
		this.Abort("400")
		return
	}
	
	//ifo := NewInfoOperation()
	switch action {
		
	}
}

func (this *MonitorController) Post() {
	beego.Info("MonitorController Post")
	action := this.GetString(conf.KEY_ACTION)
	if action == "" {
		beego.Error("[para is null] | action ")
		this.Abort("400")
		return
	}
	ifo := NewInfoOperation()
	
	switch action {
	case conf.ACTION_LOGIN:
		var ob LoginPostData
		json.Unmarshal(this.Ctx.Input.RequestBody, &ob)
		beego.Info(ob.Username)
		ts, err := ifo.login(ob.Username, ob.Password)
		if err != nil {
			beego.Error(err)
			this.Abort("400")
			return
		}
		this.Data["json"] = ts
		this.ServeJson()
	}
}

type InfoOperation struct {

}

func NewInfoOperation() *InfoOperation {
	return &InfoOperation {
	}
}

func (this *InfoOperation)login(username, password string)  (*LoginStatus, error) {
	ts := NewLoginStatus()
	ts.Status = "1"
	
	if username == "admin" && password == "admin" {
		ts.Status = "0"
	}
	
	return &ts, nil
}

