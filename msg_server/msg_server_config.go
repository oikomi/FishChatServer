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
	"os"
	"encoding/json"
	"time"
	"github.com/oikomi/FishChatServer/log"
)

type MsgServerConfig struct {
	configfile               string
	LocalIP                  string
	TransportProtocols       string
	Listen                   string
	LogFile                  string
	ScanDeadSessionTimeout   time.Duration
	Expire                   time.Duration
	MonitorBeatTime          time.Duration
	SessionManagerServerList []string
	Redis struct { 
		Addr string 
		Port string
		ConnectTimeout time.Duration
		ReadTimeout time.Duration
		WriteTimeout time.Duration
	} 
}

func NewMsgServerConfig(configfile string) *MsgServerConfig {
	return &MsgServerConfig{
		configfile : configfile,
	}
}

func (self *MsgServerConfig)LoadConfig() error {
	file, err := os.Open(self.configfile)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	defer file.Close()

	dec := json.NewDecoder(file)
	err = dec.Decode(&self)
	if err != nil {
		return err
	}
	return nil
}

func (self *MsgServerConfig)DumpConfig() {
	//fmt.Printf("Mode: %s\nListen: %s\nServer: %s\nLogfile: %s\n", 
	//cfg.Mode, cfg.Listen, cfg.Server, cfg.Logfile)
}
