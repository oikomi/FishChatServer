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
	"time"
	"encoding/json"
	"github.com/oikomi/FishChatServer/log"
	"github.com/oikomi/FishChatServer/base"
	"github.com/oikomi/FishChatServer/libnet"
	"github.com/oikomi/FishChatServer/protocol"
	"github.com/oikomi/FishChatServer/storage/redis_store"
	"github.com/oikomi/FishChatServer/storage/mongo_store"
)

type Manager struct {
	cfg          *ManagerConfig
	sessionStore *redis_store.SessionStore
	topicStore   *redis_store.TopicStore
	mongoStore   *mongo_store.MongoStore
}   

func NewManager(cfg *ManagerConfig) *Manager {
	return &Manager {
		cfg : cfg,
		sessionStore       : redis_store.NewSessionStore(redis_store.NewRedisStore(&redis_store.RedisStoreOptions {
			Network        : "tcp",
			Address        : cfg.Redis.Addr + cfg.Redis.Port,
			ConnectTimeout : time.Duration(cfg.Redis.ConnectTimeout)*time.Millisecond,
			ReadTimeout    : time.Duration(cfg.Redis.ReadTimeout)*time.Millisecond,
			WriteTimeout   : time.Duration(cfg.Redis.WriteTimeout)*time.Millisecond,
			Database       : 1,
			KeyPrefix      : base.COMM_PREFIX,
		})),
		topicStore         : redis_store.NewTopicStore(redis_store.NewRedisStore(&redis_store.RedisStoreOptions {
			Network        : "tcp",
			Address        : cfg.Redis.Addr + cfg.Redis.Port,
			ConnectTimeout : time.Duration(cfg.Redis.ConnectTimeout)*time.Millisecond,
			ReadTimeout    : time.Duration(cfg.Redis.ReadTimeout)*time.Millisecond,
			WriteTimeout   : time.Duration(cfg.Redis.WriteTimeout)*time.Millisecond,
			Database       : 1,
			KeyPrefix      : base.COMM_PREFIX,
		})),
		mongoStore         : mongo_store.NewMongoStore(cfg.Mongo.Addr, cfg.Mongo.Port, cfg.Mongo.User, cfg.Mongo.Password),
	}
}

func (self *Manager)connectMsgServer(ms string) (*libnet.Session, error) {
	client, err := libnet.Dial("tcp", ms)
	if err != nil {
		log.Error(err.Error())
		panic(err)
	}

	return client, err
}

func (self *Manager)parseProtocol(cmd []byte, session *libnet.Session) error {
	var err error
	var c protocol.CmdInternal
	//var cType interface{}
	
	err = json.Unmarshal(cmd, &c)
	if err != nil {
		log.Error("error:", err)
		return err
	}
	
	pp := NewProtoProc(self)
	
	log.Info(c)
	log.Info(c.CmdName)

	switch c.CmdName {
		case protocol.CACHE_SESSION_CMD:
			var scc SessionCacheCmd
			err = json.Unmarshal(cmd, &scc)
			if err != nil {
				log.Error("error:", err)
				return err
			}
			err = pp.procCacheSession(scc, session)
			if err != nil {
				log.Error("error:", err)
				return err
			}
			
		case protocol.CACHE_TOPIC_CMD:
			var tcc TopicCacheCmd
			err = json.Unmarshal(cmd, &tcc)
			if err != nil {
				log.Error("error:", err)
				return err
			}
			err = pp.procCacheTopic(tcc, session)
			if err != nil {
				log.Error("error:", err)
				return err
			}
		
		
		case protocol.STORE_SESSION_CMD:
			var ssd SessionStoreCmd
			err = json.Unmarshal(cmd, &ssd)
			if err != nil {
				log.Error("error:", err)
				return err
			}

			err := pp.procStoreSession(ssd.GetAnyData(), session)
			if err != nil {
				log.Error("error:", err)
				return err
			}
		

		case protocol.STORE_TOPIC_CMD:
			err = pp.procStoreTopic(c.GetAnyData(), session)
			if err != nil {
				log.Error("error:", err)
				return err
			}
		}

	return err
}

func (self *Manager)handleMsgServerClient(msc *libnet.Session) {
	msc.Process(func(msg *libnet.InBuffer) error {
		log.Info("msg_server", msc.Conn().RemoteAddr().String(),"say:", string(msg.Data))
		
		self.parseProtocol(msg.Data, msc)
		
		return nil
	})
}

func (self *Manager)subscribeChannels() error {
	log.Info("subscribeChannels")
	var msgServerClientList []*libnet.Session
	for _, ms := range self.cfg.MsgServerList {
		msgServerClient, err := self.connectMsgServer(ms)
		if err != nil {
			log.Error(err.Error())
			return err
		}
		cmd := protocol.NewCmdSimple(protocol.SUBSCRIBE_CHANNEL_CMD)
		cmd.AddArg(protocol.SYSCTRL_CLIENT_STATUS)
		cmd.AddArg(self.cfg.UUID)
		
		err = msgServerClient.Send(libnet.Json(cmd))
		if err != nil {
			log.Error(err.Error())
			return err
		}
		
		cmd = protocol.NewCmdSimple(protocol.SUBSCRIBE_CHANNEL_CMD)
		cmd.AddArg(protocol.SYSCTRL_TOPIC_STATUS)
		cmd.AddArg(self.cfg.UUID)
		
		err = msgServerClient.Send(libnet.Json(cmd))
		if err != nil {
			log.Error(err.Error())
			return err
		}

		cmd = protocol.NewCmdSimple(protocol.SUBSCRIBE_CHANNEL_CMD)
		cmd.AddArg(protocol.STORE_CLIENT_INFO)
		cmd.AddArg(self.cfg.UUID)
		
		err = msgServerClient.Send(libnet.Json(cmd))
		if err != nil {
			log.Error(err.Error())
			return err
		}
		
		msgServerClientList = append(msgServerClientList, msgServerClient)
	}

	for _, msc := range msgServerClientList {
		go self.handleMsgServerClient(msc)
	}
	return nil
}