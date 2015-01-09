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
	"time"
	"encoding/json"
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer/base"
	"github.com/oikomi/FishChatServer/libnet"
	"github.com/oikomi/FishChatServer/storage"
	"github.com/oikomi/FishChatServer/protocol"
)

type Manager struct {
	cfg          *ManagerConfig
	sessionStore *storage.SessionStore
	topicStore   *storage.TopicStore
}   

func NewManager(cfg *ManagerConfig) *Manager {
	return &Manager {
		cfg : cfg,
		sessionStore       : storage.NewSessionStore(storage.NewRedisStore(&storage.RedisStoreOptions {
			Network        : "tcp",
			Address        : cfg.Redis.Port,
			ConnectTimeout : time.Duration(cfg.Redis.ConnectTimeout)*time.Millisecond,
			ReadTimeout    : time.Duration(cfg.Redis.ReadTimeout)*time.Millisecond,
			WriteTimeout   : time.Duration(cfg.Redis.WriteTimeout)*time.Millisecond,
			Database       : 1,
			KeyPrefix      : base.COMM_PREFIX,
		})),
		topicStore         : storage.NewTopicStore(storage.NewRedisStore(&storage.RedisStoreOptions {
			Network        : "tcp",
			Address        : cfg.Redis.Port,
			ConnectTimeout : time.Duration(cfg.Redis.ConnectTimeout)*time.Millisecond,
			ReadTimeout    : time.Duration(cfg.Redis.ReadTimeout)*time.Millisecond,
			WriteTimeout   : time.Duration(cfg.Redis.WriteTimeout)*time.Millisecond,
			Database       : 1,
			KeyPrefix      : base.COMM_PREFIX,
		})),
	}
}

func (self *Manager)connectMsgServer(ms string) (*libnet.Session, error) {
	client, err := libnet.Dial("tcp", ms)
	if err != nil {
		glog.Error(err.Error())
		panic(err)
	}

	return client, err
}

func (self *Manager)parseProtocol(cmd []byte, session *libnet.Session) error {
	var c protocol.CmdInternal
	
	err := json.Unmarshal(cmd, &c)
	if err != nil {
		glog.Error("error:", err)
		return err
	}
	
	pp := NewProtoProc(self)
	
	glog.Info(c)
	glog.Info(c.CmdName)

	switch c.CmdName {
		case protocol.STORE_SESSION_CMD:
			var ssc SessionStoreCmd
			err := json.Unmarshal(cmd, &ssc)
			if err != nil {
				glog.Error("error:", err)
				return err
			}
			pp.procStoreSession(ssc, session)
		case protocol.STORE_TOPIC_CMD:
			var tsc TopicStoreCmd
			err := json.Unmarshal(cmd, &tsc)
			if err != nil {
				glog.Error("error:", err)
				return err
			}
			pp.procStoreTopic(tsc, session)
		}

	return err
}

func (self *Manager)handleMsgServerClient(msc *libnet.Session) {
	msc.Process(func(msg *libnet.InBuffer) error {
		glog.Info("msg_server", msc.Conn().RemoteAddr().String(),"say:", string(msg.Data))
		
		self.parseProtocol(msg.Data, msc)
		
		return nil
	})
}

func (self *Manager)subscribeChannels() error {
	glog.Info("subscribeChannels")
	var msgServerClientList []*libnet.Session
	for _, ms := range self.cfg.MsgServerList {
		msgServerClient, err := self.connectMsgServer(ms)
		if err != nil {
			glog.Error(err.Error())
			return err
		}
		cmd := protocol.NewCmdSimple(protocol.SUBSCRIBE_CHANNEL_CMD)
		cmd.AddArg(protocol.SYSCTRL_CLIENT_STATUS)
		cmd.AddArg(self.cfg.UUID)
		
		err = msgServerClient.Send(libnet.Json(cmd))
		if err != nil {
			glog.Error(err.Error())
			return err
		}
		
		cmd = protocol.NewCmdSimple(protocol.SUBSCRIBE_CHANNEL_CMD)
		cmd.AddArg(protocol.SYSCTRL_TOPIC_STATUS)
		cmd.AddArg(self.cfg.UUID)
		
		err = msgServerClient.Send(libnet.Json(cmd))
		if err != nil {
			glog.Error(err.Error())
			return err
		}
		
		msgServerClientList = append(msgServerClientList, msgServerClient)
	}

	for _, msc := range msgServerClientList {
		go self.handleMsgServerClient(msc)
	}
	return nil
}