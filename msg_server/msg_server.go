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
	"time"
	"flag"
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer/base"
	"github.com/oikomi/FishChatServer/libnet"
	"github.com/oikomi/FishChatServer/storage"
)

/*
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
const char* build_time(void) {
	static const char* psz_build_time = "["__DATE__ " " __TIME__ "]";
	return psz_build_time;
}
*/
import "C"

var (
	buildTime = C.GoString(C.build_time())
)

func BuildTime() string {
	return buildTime
}

const VERSION string = "0.10"

func init() {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", "false")
}

func version() {
	fmt.Printf("msg_server version %s Copyright (c) 2014 Harold Miao (miaohong@miaohong.org)  \n", VERSION)
}

var InputConfFile = flag.String("conf_file", "msg_server.json", "input conf file name")   

func handleSession(ms *MsgServer, session *libnet.Session) {
	session.Process(func(msg *libnet.InBuffer) error {
		err := ms.parseProtocol(msg.Data, session)
		if err != nil {
			glog.Error(err.Error())
		}
		
		return nil
	})
}

func main() {
	version()
	fmt.Printf("built on %s\n", BuildTime())
	flag.Parse()
	cfg := NewMsgServerConfig(*InputConfFile)
	err := cfg.LoadConfig()
	if err != nil {
		glog.Error(err.Error())
		return
	}
	
	rs := storage.NewRedisStore(&storage.RedisStoreOptions {
			Network        : "tcp",
			Address        : cfg.Redis.Port,
			ConnectTimeout : time.Duration(cfg.Redis.ConnectTimeout)*time.Millisecond,
			ReadTimeout    : time.Duration(cfg.Redis.ReadTimeout)*time.Millisecond,
			WriteTimeout   : time.Duration(cfg.Redis.WriteTimeout)*time.Millisecond,
			Database       : 1,
			KeyPrefix      : base.COMM_PREFIX,
		})

	ms := NewMsgServer(cfg, rs)

	ms.server, err = libnet.Listen(cfg.TransportProtocols, cfg.Listen)
	if err != nil {
		panic(err)
	}
	glog.Info("msg_server running at  ", ms.server.Listener().Addr().String())
	
	ms.createChannels()

	go ms.scanDeadSession()

	go ms.sendMonitorData()

	ms.server.Serve(func(session *libnet.Session) {
		glog.Info("a new client ", session.Conn().RemoteAddr().String(), " | come in")
		go handleSession(ms, session)
	})
}
