//
// Copyright 2014-2015 Hong Miao (miaohong@miaohong.org). All Rights Reserved.
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

package mongo_store

import (
	"time"
	"gopkg.in/mgo.v2"
	"github.com/oikomi/FishChatServer/log"
)

type MongoStoreOptions struct {

}

type MongoStore struct {
	opts            *MongoStoreOptions
	session         *mgo.Session
	
	
}

func NewMongoStore(ip string, port string, user string, password string) *MongoStore {
	var url string
	if user == "" && password == "" {
		url = ip + port
	} else {
		url = user + ":" + password + "@" + ip + port
	}
	
	log.Info("connect to mongo : " , url)
	maxWait := time.Duration(5 * time.Second)
	session, err := mgo.DialWithTimeout(url, maxWait)
	if err != nil {
		panic(err)
	}
	return &MongoStore {
		session : session,
	}
}

func (self *MongoStore)Init() {
	
}

func (self *MongoStore)Close() {
	self.session.Close()
}
