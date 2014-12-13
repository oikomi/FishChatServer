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

package storage

import (
	"time"
	"sync"
	"errors"
	"github.com/garyburd/redigo/redis"
)

var (
	ErrNoKeyPrefix = errors.New("cannot get session keys without a key prefix")
)

type RedisStoreOptions struct {
	Network              string
	Address              string
	ConnectTimeout       time.Duration
	ReadTimeout          time.Duration
	WriteTimeout         time.Duration
	Database             int           // Redis database to use for session keys
	KeyPrefix            string        // If set, keys will be KeyPrefix:SessionID (semicolon added)
	BrowserSessServerTTL time.Duration // Defaults to 2 days
}

type RedisStore struct {
	opts        *RedisStoreOptions
	conn        redis.Conn
	rwMutex     sync.Mutex
}

// Create a redis session store with the specified options.
func NewRedisStore(opts *RedisStoreOptions) *RedisStore {
	var err error
	rs := &RedisStore{
		opts : opts, 
		conn : nil,
		}
	rs.conn, err = redis.DialTimeout(opts.Network, opts.Address, opts.ConnectTimeout,
		opts.ReadTimeout, opts.WriteTimeout)
	if err != nil {
		panic(err)
	}
	return rs
}
