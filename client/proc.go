
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
	"github.com/oikomi/FishChatServer/libnet"
)

type Client struct {
	msgServerClient  *libnet.Session
	TopicClientMap   map[string]*libnet.Session
	Cfg              *Config

}

func NewClient(cfg *Config) *Client {
	return &Client{
		Cfg : cfg,
		TopicClientMap : make(map[string]*libnet.Session),
	}
}

func (self *Client)talkP2P() {
	for {
		fmt.Println("select what you want: ")
		fmt.Println("1: send ID you want to talk")
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
}