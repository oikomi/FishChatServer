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

package common

import (
	"time"
	"math/rand"
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer/base"
	"github.com/oikomi/FishChatServer/storage"
)

const KeyPrefix string = base.COMM_PREFIX

var DefaultRedisConnectTimeout uint32 = 2000
var DefaultRedisReadTimeout    uint32 = 1000
var DefaultRedisWriteTimeout   uint32 = 1000

var DefaultRedisOptions storage.RedisStoreOptions = storage.RedisStoreOptions {
	Network        : "tcp",
	Address        : ":6379",
	ConnectTimeout : time.Duration(DefaultRedisConnectTimeout)*time.Millisecond,
	ReadTimeout    : time.Duration(DefaultRedisReadTimeout)*time.Millisecond,
	WriteTimeout   : time.Duration(DefaultRedisWriteTimeout)*time.Millisecond,
	Database       : 1,
	KeyPrefix      : base.COMM_PREFIX,
}

//Just use random to select msg_server
func SelectServer(serverList []string, serverNum int) string {
	return serverList[rand.Intn(serverNum)]
}

func GetSessionFromCID(sessionStore  *storage.SessionStore, ID string) (*storage.SessionStoreData, error) {
	session ,err := sessionStore.Get(ID)
	
	if err != nil {
		glog.Warningf("no ID : %s", ID)
		return nil, err
	}
	if session != nil {
		glog.Info(session)
	}
	
	return session, nil
}

func DelSessionFromCID(sessionStore *storage.SessionStore, ID string) error {
	err := sessionStore.Delete(ID)
	
	if err != nil {
		glog.Warningf("no ID : %s", ID)
		return err
	}

	return nil
}

func GetTopicFromTopicName(topicStore *storage.TopicStore, topicName string) (*storage.TopicStoreData, error) {
	topic ,err := topicStore.Get(topicName)
	
	if err != nil {
		glog.Warningf("no topicName : %s", topicName)
		return nil, err
	}
	if topic != nil {
		glog.Info(topic)
	}
	
	return topic, nil
}

func GetOfflineMsgFromOwnerName(offlineMsgStore *storage.OfflineMsgStore, ownerName string) (*storage.OfflineMsgStoreData, error) {
	o ,err := offlineMsgStore.Get(ownerName)
	
	if err != nil {
		glog.Warningf("no ownerName : %s", ownerName)
		return nil, err
	}
	if o != nil {
		glog.Info(o)
	}
	
	return o, nil
}


