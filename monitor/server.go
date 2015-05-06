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
	"sync"
	"time"
	"encoding/json"
	"github.com/oikomi/FishChatServer/log"
	"github.com/oikomi/FishChatServer/base"
	"github.com/oikomi/FishChatServer/libnet"
	"github.com/oikomi/FishChatServer/protocol"
	"github.com/oikomi/FishChatServer/storage"
)

type Monitor struct {
	cfg                 *MonitorConfig
	msgServerClientMap  map[string]*libnet.Session
	sessionStore        *storage.SessionStore
	topicServerMap      map[string]string
	readMutex           sync.Mutex
}   

func NewMonitor(cfg *MonitorConfig) *Monitor {
	return &Monitor {
		cfg                : cfg,
		msgServerClientMap : make(map[string]*libnet.Session),
		sessionStore       : storage.NewSessionStore(storage.NewRedisStore(&storage.RedisStoreOptions {
					Network :   "tcp",
					Address :   cfg.Redis.Port,
					ConnectTimeout : time.Duration(cfg.Redis.ConnectTimeout)*time.Millisecond,
					ReadTimeout : time.Duration(cfg.Redis.ReadTimeout)*time.Millisecond,
					WriteTimeout : time.Duration(cfg.Redis.WriteTimeout)*time.Millisecond,
					Database :  1,
					KeyPrefix : base.COMM_PREFIX,
		})),
		topicServerMap     : make(map[string]string),
	}
}

func (self *Monitor)connectMsgServer(ms string) (*libnet.Session, error) {
	client, err := libnet.Dial("tcp", ms)
	if err != nil {
		log.Error(err.Error())
		panic(err)
	}

	return client, err
}

func (self *Monitor)handleMsgServerClient(msc *libnet.Session) {
	msc.Process(func(msg *libnet.InBuffer) error {
		log.Info("msg_server", msc.Conn().RemoteAddr().String()," say: ", string(msg.Data))
		var c protocol.CmdMonitor
		
		err := json.Unmarshal(msg.Data, &c)
		if err != nil {
			log.Error("error:", err)
			return err
		}

		return nil
	})
}

func (self *Monitor)subscribeChannels() error {
	log.Info("monitor start to subscribeChannels")
	for _, ms := range self.cfg.MsgServerList {
		msgServerClient, err := self.connectMsgServer(ms)
		if err != nil {
			log.Error(err.Error())
			return err
		}
		cmd := protocol.NewCmdSimple(protocol.SUBSCRIBE_CHANNEL_CMD)
		cmd.AddArg(protocol.SYSCTRL_MONITOR)
		cmd.AddArg(self.cfg.UUID)
		
		err = msgServerClient.Send(libnet.Json(cmd))
		if err != nil {
			log.Error(err.Error())
			return err
		}
		
		
		self.msgServerClientMap[ms] = msgServerClient
	}

	for _, msc := range self.msgServerClientMap {
		go self.handleMsgServerClient(msc)
	}
	return nil
}
