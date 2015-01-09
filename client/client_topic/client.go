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
	"encoding/json"
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer/protocol"
	"github.com/oikomi/FishChatServer/libnet"
)

type Client struct {
	TopicClientMap   map[string]*libnet.Session
	Cfg                  *Config
}

func NewClient(cfg *Config) *Client {
	return &Client {
		Cfg            : cfg,
		TopicClientMap : make(map[string]*libnet.Session),
	}

}

func (self *Client)AddTopicClient(s string, c *libnet.Session) {
	self.TopicClientMap[s] = c
}

func (self *Client)locateTopicAddr(c protocol.CmdSimple) {
	glog.Info("locateTopicAddr")
	serverAddr := c.GetArgs()[0]
	topicName := c.GetArgs()[1]

	msgServerClient, err := libnet.Dial("tcp", serverAddr)
	if err != nil {
		panic(err)
	}
	
	go msgServerClient.Process(func(msg *libnet.InBuffer) error {
		glog.Info(string(msg.Data))
		return nil
	})
	
	self.AddTopicClient(topicName, msgServerClient)
	
	cmd := protocol.NewCmdSimple(protocol.JOIN_TOPIC_CMD)
	cmd.AddArg(topicName)
	cmd.AddArg(gClientID)
	
	err = msgServerClient.Send(libnet.Json(cmd))
	if err != nil {
		glog.Error(err.Error())
	}

}

func (self *Client)parseProtocol(cmd []byte) error {
	var c protocol.CmdSimple
	
	err := json.Unmarshal(cmd, &c)
	if err != nil {
		glog.Error("error:", err)
		return err
	}

	switch c.GetCmdName() {
		case protocol.LOCATE_TOPIC_MSG_ADDR_CMD:
			self.locateTopicAddr(c)

		}
	
	return err
}