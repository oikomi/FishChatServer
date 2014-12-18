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
	c := NewClient(&cfg)

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

	c.msgServerClient, err = libnet.Dial("tcp", string(inMsg.Get()), p)
	if err != nil {
		panic(err)
	}
	for {
		fmt.Println("select what you want: ")
		fmt.Println("1: talk to one people")
		var input string
		if _, err := fmt.Scanf("%s\n", &input); err != nil {
			glog.Error(err.Error())
		}
		switch input {
		case "1":
			c.talkP2P()
		default:
			fmt.Println("it's a weekday")
		}
	}
	
	glog.Flush()
}
