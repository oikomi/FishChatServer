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
	"time"
	"math/rand"
	"github.com/oikomi/FishChatServer/log"
	"github.com/oikomi/FishChatServer/base"
	"github.com/oikomi/FishChatServer/storage/redis_store"
)

const KeyPrefix string = base.COMM_PREFIX

var DefaultRedisConnectTimeout uint32 = 2000
var DefaultRedisReadTimeout    uint32 = 1000
var DefaultRedisWriteTimeout   uint32 = 1000

var DefaultRedisOptions redis_store.RedisStoreOptions = redis_store.RedisStoreOptions {
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

func GetSessionFromCID(storeOp  interface{}, ID string) (*redis_store.SessionCacheData, error) {
	
	switch storeOp.(type) {
		case *redis_store.SessionCache:	
			session ,err := storeOp.(*redis_store.SessionCache).Get(ID)
			
			if err != nil {
				log.Warningf("no ID : %s", ID)
				return nil, err
			}
			if session != nil {
				log.Info(session)
			}
			
			return session, nil
			
	}
	
	return nil, nil
	
//	session ,err := sessionCache.Get(ID)
	
//	if err != nil {
//		log.Warningf("no ID : %s", ID)
//		return nil, err
//	}
//	if session != nil {
//		log.Info(session)
//	}
	
//	return session, nil
}

func DelSessionFromCID(storeOp  interface{}, ID string) error {
	switch storeOp.(type) {
		case *redis_store.SessionCache:
			err := storeOp.(*redis_store.SessionCache).Delete(ID)
			
			if err != nil {
				log.Warningf("no ID : %s", ID)
				return err
			}			
	}
	
//	err := sessionCache.Delete(ID)
	
//	if err != nil {
//		log.Warningf("no ID : %s", ID)
//		return err
//	}

	return nil
}

func GetTopicFromTopicName(storeOp  interface{}, topicName string) (*redis_store.TopicCacheData, error) {
	switch storeOp.(type) {
		case *redis_store.TopicCache:
			topic ,err := storeOp.(*redis_store.TopicCache).Get(topicName)
			
			if err != nil {
				log.Warningf("no topicName : %s", topicName)
				return nil, err
			}
			if topic != nil {
				log.Info(topic)
			}
			
			return topic, nil			
	}
	
	return nil, nil
	
//	topic ,err := topicCache.Get(topicName)
	
//	if err != nil {
//		log.Warningf("no topicName : %s", topicName)
//		return nil, err
//	}
//	if topic != nil {
//		log.Info(topic)
//	}
	
//	return topic, nil
}

func GetOfflineMsgFromOwnerName(storeOp  interface{}, ownerName string) (*redis_store.OfflineMsgCacheData, error) {
	switch storeOp.(type) {
		case *redis_store.OfflineMsgCache:
			o ,err := storeOp.(*redis_store.OfflineMsgCache).Get(ownerName)
			
			if err != nil {
				log.Warningf("no ownerName : %s", ownerName)
				return nil, err
			}
			if o != nil {
				log.Info(o)
			}
			
			return o, nil		
	}
	
	
	return nil, nil
	
//	o ,err := offlineMsgCache.Get(ownerName)
	
//	if err != nil {
//		log.Warningf("no ownerName : %s", ownerName)
//		return nil, err
//	}
//	if o != nil {
//		log.Info(o)
//	}
	
//	return o, nil
}


