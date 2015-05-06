//
// Copyright 2014 Hong Miao (miaohong@miaohong.org). All Rights Reserved.
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
	"github.com/oikomi/FishChatServer/log"
	"github.com/oikomi/FishChatServer/libnet"
	"github.com/oikomi/FishChatServer/protocol"
	"github.com/oikomi/FishChatServer/storage"
)

func init() {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", "false")
}

type ProtoProc struct {
	Manager   *Manager
}

func NewProtoProc(m *Manager) *ProtoProc {
	return &ProtoProc {
		Manager : m,
	}
}

func (self *ProtoProc)procStoreSession(cmd protocol.Cmd, session *libnet.Session) error {
	log.Info("procStoreSession")
	var err error
	log.Info(cmd.GetAnyData())
	err = self.Manager.sessionStore.Set(cmd.GetAnyData().(*storage.SessionStoreData))
	if err != nil {
		log.Error("error:", err)
	}
	log.Info("set sesion id success")
	
	return nil
}

func (self *ProtoProc)procStoreTopic(cmd protocol.Cmd, session *libnet.Session) error {
	log.Info("procStoreTopic")
	var err error
	log.Info(cmd.GetAnyData())
	err = self.Manager.topicStore.Set(cmd.GetAnyData().(*storage.TopicStoreData))
	if err != nil {
		log.Error("error:", err)
	}
	log.Info("set sesion id success")
	
	return nil
}
