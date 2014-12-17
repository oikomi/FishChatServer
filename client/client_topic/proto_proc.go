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
	"github.com/oikomi/FishChatServer/libnet"
	"github.com/oikomi/FishChatServer/protocol"
)

func locateTopicAddr(c protocol.CmdSimple) {
	glog.Info("locateTopicAddr")
	serverAddr := c.GetArgs()[0]
	topicName := c.GetArgs()[1]
	
	p := libnet.PacketN(2, libnet.BigEndianBO, libnet.LittleEndianBF)
	
	msgServerClient, err := libnet.Dial("tcp", serverAddr, p)
	if err != nil {
		panic(err)
	}
	
	cmd := protocol.NewCmdSimple()
	
	cmd.CmdName = protocol.JOIN_TOPIC_CMD
	
	cmd.Args = append(cmd.Args, topicName)
	
	cmd.Args = append(cmd.Args, gClientID)
	
	err = msgServerClient.Send(libnet.JSON {
		cmd,
	})
	if err != nil {
		glog.Error(err.Error())
	}

}

func parseProtocol(cmd []byte) error {
	var c protocol.CmdSimple
	
	err := json.Unmarshal(cmd, &c)
	if err != nil {
		glog.Error("error:", err)
		return err
	}
	
	glog.Info(c.CmdName)

	switch c.CmdName {
		case protocol.LOCATE_TOPIC_MSG_ADDR_CMD:
			locateTopicAddr(c)

		}
	
	return err
}