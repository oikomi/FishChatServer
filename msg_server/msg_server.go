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
	"flag"
	"fmt"
	"github.com/golang/glog"
	"github.com/funny/link"
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
	fmt.Printf("msg_server version %s Copyright (c) 2014 Harold Miao (miaohonghit@gmail.com)  \n", VERSION)
}

var InputConfFile = flag.String("conf_file", "msg_server.json", "input conf file name")   

func handleSession(ms *MsgServer, session *link.Session) {
	session.ReadLoop(func(msg link.InBuffer) {
		glog.Info(string(msg.Get()))
		
		err := ms.parseProtocol(msg.Get(), session)
		if err != nil {
			glog.Error(err.Error())
		}
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
	
	ms := NewMsgServer(cfg)
	
	p := link.PacketN(2, link.BigEndianBO, link.LittleEndianBF)
	
	ms.server, err = link.Listen(cfg.TransportProtocols, cfg.Listen, p)
	if err != nil {
		panic(err)
	}
	glog.Info("server start:", ms.server.Listener().Addr().String())
	
	ms.createChannels()
	go ms.scanDeadSession()

	ms.server.AcceptLoop(func(session *link.Session) {
		glog.Info("client ", session.Conn().RemoteAddr().String(), " | in")
		
		go handleSession(ms, session)
	})
}
