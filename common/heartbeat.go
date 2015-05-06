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

package common

import (
	"sync"
	"time"
	"github.com/oikomi/FishChatServer/log"
	"github.com/oikomi/FishChatServer/libnet"
	"github.com/oikomi/FishChatServer/protocol"
)

type HeartBeat struct {
	name       string
	session    *libnet.Session
	mu         sync.Mutex
	timeout    time.Duration
	expire     time.Duration
	fails      uint64
	threshold  uint64
}

func NewHeartBeat(name string, session *libnet.Session, timeout time.Duration, expire time.Duration, limit uint64) *HeartBeat {
	return &HeartBeat {
		name      : name,
		session   : session,
		timeout   : timeout,
		expire    : expire,
		threshold : limit,
	}
}

func (self *HeartBeat) ResetFailures() {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.fails = 0
}

func (self *HeartBeat) ChangeThreshold(thres uint64) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.threshold = thres
}

func (self *HeartBeat) Beat() {
	timer := time.NewTicker(self.timeout * time.Second)
	//ttl := time.After(self.expire * time.Second)
	for {
		select {
		case <-timer.C:
			go func() {
				cmd := protocol.NewCmdSimple(protocol.SEND_PING_CMD)
				cmd.AddArg(protocol.PING)
				err := self.session.Send(libnet.Json(cmd))
				if err != nil {
					log.Error(err.Error())
				}
			}()
		//case <-ttl:
			//break
		}
	}
}

func (self *HeartBeat) Receive() {
	timeout := time.After(self.timeout)
	for {
		select {
		case <-timeout:
			self.fails = self.fails + 1
			if self.fails > self.threshold {
				return
			}
		}
	}
}

