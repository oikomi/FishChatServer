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
	"sync"
	"time"
	"encoding/json"
	"github.com/oikomi/FishChatServer/log"
	"github.com/oikomi/FishChatServer/base"
	"github.com/oikomi/FishChatServer/libnet"
	"github.com/oikomi/FishChatServer/protocol"
	"github.com/oikomi/FishChatServer/storage/redis_store"
	"github.com/oikomi/FishChatServer/storage/mongo_store"
)

type Router struct {
	cfg                 *RouterConfig
	msgServerClientMap  map[string]*libnet.Session
	sessionCache        *redis_store.SessionCache
	mongoStore          *mongo_store.MongoStore
	topicServerMap      map[string]string
	readMutex           sync.Mutex
}   

func NewRouter(cfg *RouterConfig) *Router {
	return &Router {
		cfg                : cfg,
		msgServerClientMap : make(map[string]*libnet.Session),
		sessionCache       : redis_store.NewSessionCache(redis_store.NewRedisStore(&redis_store.RedisStoreOptions {
					Network :   "tcp",
					Address :   cfg.Redis.Addr + cfg.Redis.Port,
					ConnectTimeout : time.Duration(cfg.Redis.ConnectTimeout)*time.Millisecond,
					ReadTimeout : time.Duration(cfg.Redis.ReadTimeout)*time.Millisecond,
					WriteTimeout : time.Duration(cfg.Redis.WriteTimeout)*time.Millisecond,
					Database :  1,
					KeyPrefix : base.COMM_PREFIX,
		})),
		mongoStore         : mongo_store.NewMongoStore(cfg.Mongo.Addr, cfg.Mongo.Port, cfg.Mongo.User, cfg.Mongo.Password),
		topicServerMap     : make(map[string]string),
	}
}

func (self *Router)connectMsgServer(ms string) (*libnet.Session, error) {
	client, err := libnet.Dial("tcp", ms)
	if err != nil {
		log.Error(err.Error())
		panic(err)
	}

	return client, err
}

func (self *Router)handleMsgServerClient(msc *libnet.Session) {
	msc.Process(func(msg *libnet.InBuffer) error {
		log.Info("msg_server", msc.Conn().RemoteAddr().String()," say: ", string(msg.Data))
		var c protocol.CmdInternal
		pp := NewProtoProc(self)
		err := json.Unmarshal(msg.Data, &c)
		if err != nil {
			log.Error("error:", err)
			return err
		}
		switch c.GetCmdName() {
			case protocol.SEND_MESSAGE_P2P_CMD:
				err := pp.procSendMsgP2P(c, msc)
				if err != nil {
					log.Warning(err.Error())
				}
			case protocol.CREATE_TOPIC_CMD:
				err := pp.procCreateTopic(c, msc)
				if err != nil {
					log.Warning(err.Error())
				}
			case protocol.JOIN_TOPIC_CMD:
				err := pp.procJoinTopic(c, msc)
				if err != nil {
					log.Warning(err.Error())
				}
			case protocol.SEND_MESSAGE_TOPIC_CMD:
				err := pp.procSendMsgTopic(c, msc)
				if err != nil {
					log.Warning(err.Error())
				}
				
			}
		return nil
	})
}

func (self *Router)subscribeChannels() error {
	log.Info("route start to subscribeChannels")
	for _, ms := range self.cfg.MsgServerList {
		msgServerClient, err := self.connectMsgServer(ms)
		if err != nil {
			log.Error(err.Error())
			return err
		}
		cmd := protocol.NewCmdSimple(protocol.SUBSCRIBE_CHANNEL_CMD)
		cmd.AddArg(protocol.SYSCTRL_SEND)
		cmd.AddArg(self.cfg.UUID)
		
		err = msgServerClient.Send(libnet.Json(cmd))
		if err != nil {
			log.Error(err.Error())
			return err
		}
		
		cmd = protocol.NewCmdSimple(protocol.SUBSCRIBE_CHANNEL_CMD)
		cmd.AddArg(protocol.SYSCTRL_TOPIC_SYNC)
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
