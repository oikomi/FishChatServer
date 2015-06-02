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

package redis_store

import (
	"sync"
	"time"
	"encoding/json"
	"github.com/garyburd/redigo/redis"
)

type TopicCache struct {
	RS       *RedisStore
	rwMutex  sync.Mutex
}

func NewTopicCache(RS *RedisStore) *TopicCache {
	return &TopicCache {
		RS    : RS,
	}
}

type TopicCacheData struct {
	TopicName     string
	CreaterID     string
	MemberList    []*Member
	MsgServerAddr string
	MaxAge        time.Duration
}

type Member struct {
	ID   string
}

func NewMember(ID string) *Member {
	return &Member {
		ID : ID,
	}
}

func NewTopicCacheData(TopicName string, CreaterID string, MsgServerAddr string) *TopicCacheData {
	return &TopicCacheData {
		TopicName     : TopicName,
		CreaterID     : CreaterID,
		MemberList    : make([]*Member, 0),
		MsgServerAddr : MsgServerAddr,
	}
}

func (self *TopicCacheData)StoreKey() string {
	return self.TopicName
}

func (self *TopicCacheData)AddMember(m *Member) {
	self.MemberList = append(self.MemberList, m)
}

// Get the session from the store.
func (self *TopicCache) Get(k string) (*TopicCacheData, error) {
	self.rwMutex.Lock()
	defer self.rwMutex.Unlock()
	key := k + TOPIC_UNIQ_PREFIX
	if self.RS.opts.KeyPrefix != "" {
		key = self.RS.opts.KeyPrefix + ":" + k + TOPIC_UNIQ_PREFIX
	}
	b, err := redis.Bytes(self.RS.conn.Do("GET", key))
	if err != nil {
		return nil, err
	}
	var sess TopicCacheData
	err = json.Unmarshal(b, &sess)
	if err != nil {
		return nil, err
	}
	return &sess, nil
}

// Save the session into the store.
func (self *TopicCache) Set(sess *TopicCacheData) error {
	self.rwMutex.Lock()
	defer self.rwMutex.Unlock()
	b, err := json.Marshal(sess)
	if err != nil {
		return err
	}
	key := sess.TopicName + TOPIC_UNIQ_PREFIX
	if self.RS.opts.KeyPrefix != "" {
		key = self.RS.opts.KeyPrefix + ":" + sess.TopicName + TOPIC_UNIQ_PREFIX
	}
	ttl := sess.MaxAge
	if ttl == 0 {
		// Browser session, set to specified TTL
		ttl = self.RS.opts.BrowserSessServerTTL
		if ttl == 0 {
			ttl = 2 * 24 * time.Hour // Default to 2 days
		}
	}
	_, err = self.RS.conn.Do("SETEX", key, int(ttl.Seconds()), b)
	if err != nil {
		return err
	}
	return nil
}

// Delete the session from the store.
func (self *TopicCache) Delete(id string) error {
	self.rwMutex.Lock()
	defer self.rwMutex.Unlock()
	key := id + TOPIC_UNIQ_PREFIX
	if self.RS.opts.KeyPrefix != "" {
		key = self.RS.opts.KeyPrefix + ":" + id + TOPIC_UNIQ_PREFIX
	}
	_, err := self.RS.conn.Do("DEL", key)
	if err != nil {
		return err
	}
	return nil
}

// Clear all sessions from the store. Requires the use of a key
// prefix in the store options, otherwise the method refuses to delete all keys.
func (self *TopicCache) Clear() error {
	self.rwMutex.Lock()
	defer self.rwMutex.Unlock()
	vals, err := self.getSessionKeys()
	if err != nil {
		return err
	}
	if len(vals) > 0 {
		self.RS.conn.Send("MULTI")
		for _, v := range vals {
			self.RS.conn.Send("DEL", v)
		}
		_, err = self.RS.conn.Do("EXEC")
		if err != nil {
			return err
		}
	}
	return nil
}
// Get the number of session keys in the store. Requires the use of a
// key prefix in the store options, otherwise returns -1 (cannot tell
// session keys from other keys).
func (self *TopicCache) Len() int {
	self.rwMutex.Lock()
	defer self.rwMutex.Unlock()
	vals, err := self.getSessionKeys()
	if err != nil {
		return -1
	}
	return len(vals)
}
func (self *TopicCache) getSessionKeys() ([]interface{}, error) {
	self.rwMutex.Lock()
	defer self.rwMutex.Unlock()
	if self.RS.opts.KeyPrefix != "" {
		return redis.Values(self.RS.conn.Do("KEYS", self.RS.opts.KeyPrefix+":*"))
	}
	return nil, ErrNoKeyPrefix
}