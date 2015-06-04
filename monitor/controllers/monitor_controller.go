//
// Copyright 2014-2015 Hong Miao (miaohong@miaohong.org). All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

