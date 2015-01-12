//
// Copyright 2014 Hong Miao. All Rights Reserved.
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

package main

import (
	"flag"
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer/libnet"
	"github.com/oikomi/FishChatServer/common"
	"github.com/oikomi/FishChatServer/protocol"
)

func init() {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", "false")
}

type ProtoProc struct {
	gateway    *Gateway
}

func NewProtoProc(gateway *Gateway) *ProtoProc {
	return &ProtoProc {
		gateway : gateway,
	}
}

func (self *ProtoProc)procReqMsgServer(cmd protocol.Cmd, session *libnet.Session) error {
	glog.Info("procReqMsgServer")
	var err error
	msgServer := common.SelectServer(self.gateway.cfg.MsgServerList, self.gateway.cfg.MsgServerNum)

	resp := protocol.NewCmdSimple(protocol.SELECT_MSG_SERVER_FOR_CLIENT_CMD)
	resp.AddArg(msgServer)
	
	glog.Info(resp)
	
	if session != nil {
		err = session.Send(libnet.Json(resp))
		if err != nil {
			glog.Error(err.Error())
		}
		session.Close()
		glog.Info("client ", session.Conn().RemoteAddr().String(), " | close")
	}
	return nil
}