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
	"sync"
	"time"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/oikomi/FishChatServer/log"
)

type MongoStoreOptions struct {

}

type MongoStore struct {
	opts            *MongoStoreOptions
	session         *mgo.Session
	
	rwMutex         sync.Mutex
	
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
	//self.session.DB("im").C("client_info")
	
}

func (self *MongoStore)Update(db string, c string, data interface{}) error {
	log.Info("MongoStore Update")
	var err error
	self.rwMutex.Lock()
	defer self.rwMutex.Unlock()
	
	op := self.session.DB(db).C(c)
	
	switch data.(type) {
		case *SessionStoreData:
			cid := data.(*SessionStoreData).ClientID
			log.Info("cid : " , cid)
			_, err = op.Upsert(bson.M{"clientid": cid}, data.(*SessionStoreData))
			if err != nil {
				log.Error(err.Error())
				return err
			}	
	}
	
	return err
}



func (self *MongoStore)Close() {
	self.session.Close()
}
