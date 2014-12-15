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
	"fmt"
	"flag"
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer/libnet"
	"github.com/oikomi/FishChatServer/protocol"
	"github.com/oikomi/FishChatServer/common"
)

var InputConfFile = flag.String("conf_file", "client.json", "input conf file name")   

func init() {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", "false")
}

func heartBeat(cfg Config, msgServerClient *libnet.Session) {
	hb := common.NewHeartBeat("client", msgServerClient, cfg.HeartBeatTime, cfg.Expire, 10)
	hb.Beat()
}

func main() {
	flag.Parse()
	cfg, err := LoadConfig(*InputConfFile)
	if err != nil {
		glog.Error(err.Error())
		return
	}

	p := libnet.PacketN(2, libnet.BigEndianBO, libnet.LittleEndianBF)

	gatewayClient, err := libnet.Dial("tcp", cfg.GatewayServer, p)
	if err != nil {
		panic(err)
	}
	
	fmt.Println("input id :")
	var input string
	if _, err := fmt.Scanf("%s\n", &input); err != nil {
		glog.Error(err.Error())
	}
	
	inMsg, err := gatewayClient.Read()
	if err != nil {
		glog.Error(err.Error())
	}
	glog.Info(string(inMsg.Get()))

	gatewayClient.Close(nil)

	msgServerClient, err := libnet.Dial("tcp", string(inMsg.Get()), p)
	if err != nil {
		panic(err)
	}
	
	glog.Info("test.. send id...")
	cmd := protocol.NewCmdSimple()
	
	cmd.CmdName = protocol.SEND_CLIENT_ID_CMD
	cmd.Args = append(cmd.Args, input)
	
	err = msgServerClient.Send(libnet.JSON {
		cmd,
	})
	if err != nil {
		glog.Error(err.Error())
	}
	
	go heartBeat(cfg, msgServerClient)
	
	glog.Info("test.. send p2p msg...")
	cmd = protocol.NewCmdSimple()
	
	cmd.CmdName = protocol.SEND_MESSAGE_P2P_CMD
	
	fmt.Println("input 2id :")
	if _, err = fmt.Scanf("%s\n", &input); err != nil {
		glog.Error(err.Error())
	}
	
	cmd.Args = append(cmd.Args, input)
	
	fmt.Println("input msg :")
	if _, err = fmt.Scanf("%s\n", &input); err != nil {
		glog.Error(err.Error())
	}
	
	cmd.Args = append(cmd.Args, input)
	
	err = msgServerClient.Send(libnet.JSON {
		cmd,
	})
	if err != nil {
		glog.Error(err.Error())
	}
	
	defer msgServerClient.Close(nil)
	
	msgServerClient.ReadLoop(func(msg libnet.InBuffer) {
		glog.Info(string(msg.Get()))
	})
	
	glog.Flush()
}
