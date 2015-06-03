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
	"github.com/oikomi/FishChatServer/storage/redis_store"
	"github.com/oikomi/FishChatServer/storage/mongo_store"
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

func (self *ProtoProc)procCacheSession(cmd protocol.Cmd, session *libnet.Session) error {
	log.Info("procCacheSession")
	var err error
	log.Info(cmd.GetAnyData())
	err = self.Manager.sessionCache.Set(cmd.GetAnyData().(*redis_store.SessionCacheData))
	if err != nil {
		return err
		log.Error("error:", err)
	}
	log.Info("set sesion id success")
	
	return nil
}

func (self *ProtoProc)procCacheTopic(cmd protocol.Cmd, session *libnet.Session) error {
	log.Info("procCacheTopic")
	var err error
	log.Info(cmd.GetAnyData())
	err = self.Manager.topicCache.Set(cmd.GetAnyData().(*redis_store.TopicCacheData))
	if err != nil {
		return err
		log.Error("error:", err)
	}
	log.Info("set sesion id success")
	
	return nil
}


func (self *ProtoProc)procStoreSession(data interface{}, session *libnet.Session) error {
	log.Info("procStoreSession")
	var err error
	log.Info(data)
	err = self.Manager.mongoStore.Update(mongo_store.DATA_BASE_NAME, mongo_store.CLIENT_INFO_COLLECTION, data)
	if err != nil {
		return err
		log.Error("error:", err)
	}
	
	return nil
}

func (self *ProtoProc)procStoreTopic(data interface{}, session *libnet.Session) error {
	log.Info("procStoreTopic")
	var err error
	log.Info(data)
	
	err = self.Manager.mongoStore.Update(mongo_store.DATA_BASE_NAME, mongo_store.TOPIC_INFO_COLLECTION, data)
	if err != nil {
		return err
		log.Error("error:", err)
	}
	
	return nil
}
