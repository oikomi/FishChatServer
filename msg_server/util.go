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
	"net"
	"sync"
	"time"
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer/libnet"
	"github.com/oikomi/FishChatServer/protocol"
)

type MonitorBeat struct {
	name       string
	session    *libnet.Session
	mu         sync.Mutex
	timeout    time.Duration
	expire     time.Duration
	fails      uint64
	threshold  uint64
}

func NewMonitorBeat(name string, timeout time.Duration, expire time.Duration, limit uint64) *MonitorBeat {
	return &MonitorBeat {
		name      : name,
		timeout   : timeout,
		expire    : expire,
		threshold : limit,
	}
}

func (self *MonitorBeat) ResetFailures() {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.fails = 0
}

func (self *MonitorBeat) ChangeThreshold(thres uint64) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.threshold = thres
}

func (self *MonitorBeat) Beat(c *libnet.Channel, data *protocol.CmdMonitor) {
	timer := time.NewTicker(self.timeout * time.Second)
	//ttl := time.After(self.expire * time.Second)
	//for {
		select {
		case <-timer.C:
			go func() {
				_, err := c.Broadcast(libnet.Json(data))
				if err != nil {
					glog.Error(err.Error())
					//return err
				}
			}()
		//case <-ttl:
			//break
		}
	//}

	//return nil
}

func (self *MonitorBeat) Receive() {
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

func getHostIP() {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}
	for _, addr := range addrs {
		fmt.Println(addr.String())
	}
}

