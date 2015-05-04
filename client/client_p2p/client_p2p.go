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
	"encoding/json"
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

	fmt.Println("req GatewayServer...")

	gatewayClient, err := libnet.Dial("tcp", cfg.GatewayServer)
	if err != nil {
		panic(err)
	}
	
	//glog.Info("...")
	cmd := protocol.NewCmdSimple(protocol.REQ_MSG_SERVER_CMD)
	
	err = gatewayClient.Send(libnet.Json(cmd))
	if err != nil {
		glog.Error(err.Error())
	}
	
	fmt.Println("input id :")
	var input string
	if _, err := fmt.Scanf("%s\n", &input); err != nil {
		glog.Error(err.Error())
	}
	var c protocol.CmdSimple
	err = gatewayClient.ProcessOnce(func(msg *libnet.InBuffer) error {
		glog.Info(string(msg.Data))
		err = json.Unmarshal(msg.Data, &c)
		if err != nil {
			glog.Error("error:", err)
		}
		return nil
	})
	if err != nil {
		glog.Error(err.Error())
	}

	gatewayClient.Close()

	fmt.Println("req GatewayServer end...")

	msgServerClient, err := libnet.Dial("tcp", string(c.GetArgs()[0]))
	if err != nil {
		panic(err)
	}
	
	fmt.Println("send your id...")
	cmd = protocol.NewCmdSimple(protocol.SEND_CLIENT_ID_CMD)
	cmd.AddArg(input)
	
	err = msgServerClient.Send(libnet.Json(cmd))
	if err != nil {
		glog.Error(err.Error())
	}
	
	go heartBeat(cfg, msgServerClient)
	
	//glog.Info("the msg you want to send...")


	go msgServerClient.Process(func(msg *libnet.InBuffer) error {
		//glog.Info(string(msg.Data))
		err = json.Unmarshal(msg.Data, &c)
		if err != nil {
			glog.Error("error:", err)
		}
		fmt.Println(c.GetArgs()[1] + "  says : " + c.GetArgs()[0])

		return nil
	})

	for {
		cmd = protocol.NewCmdSimple(protocol.SEND_MESSAGE_P2P_CMD)
		
		fmt.Println("send the id you want to talk :")
		if _, err = fmt.Scanf("%s\n", &input); err != nil {
			glog.Error(err.Error())
		}
		
		cmd.AddArg(input)
		
		fmt.Println("input msg :")
		if _, err = fmt.Scanf("%s\n", &input); err != nil {
			glog.Error(err.Error())
		}
		
		cmd.AddArg(input)
		
		fmt.Println("input my ID :")
		if _, err = fmt.Scanf("%s\n", &input); err != nil {
			glog.Error(err.Error())
		}
		
		cmd.AddArg(input)
		
		err = msgServerClient.Send(libnet.Json(cmd))
		if err != nil {
			glog.Error(err.Error())
		}
	}
	
	defer msgServerClient.Close()
	
	// msgServerClient.Process(func(msg *libnet.InBuffer) error {
	// 	glog.Info(string(msg.Data))
	// 	return nil
	// })
	
	glog.Flush()
}
