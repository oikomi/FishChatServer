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
	"encoding/json"
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer/libnet"
	"github.com/oikomi/FishChatServer/protocol"
)

func init() {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", "false")
}

type Gateway struct {
	cfg     *GatewayConfig
	server  *libnet.Server
}

func NewGateway(cfg *GatewayConfig) *Gateway {
	return &Gateway {
		cfg    : cfg,
		server : new(libnet.Server),
	}
}

func (self *Gateway)parseProtocol(cmd []byte, session *libnet.Session) error {
	var c protocol.CmdSimple
	err := json.Unmarshal(cmd, &c)
	if err != nil {
		glog.Error("error:", err)
		return err
	}
	
	pp := NewProtoProc(self)

	switch c.GetCmdName() {
		case protocol.REQ_MSG_SERVER_CMD:
			pp.procReqMsgServer(&c, session)
		}

	return err
}

