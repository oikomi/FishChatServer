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

package protocol

import (
	"github.com/oikomi/FishChatServer/libnet"
	"github.com/oikomi/FishChatServer/storage"
)

type TopicMap   map[string]*Topic

type Topic struct {
	TopicName     string
	MsgAddr       string
	Channel       *libnet.Channel
	TA            *TopicAttribute
	ClientIDList  []string
	TSD           *storage.TopicStoreData
}

func NewTopic(topicName string, msgAddr string, CreaterID string, CreaterSession *libnet.Session) *Topic {
	return &Topic {
		TopicName    : topicName,
		MsgAddr      : msgAddr,
		Channel      : new(libnet.Channel),
		TA           : NewTopicAttribute(CreaterID, CreaterSession),
		ClientIDList : make([]string, 0),
	}
}

func (self *Topic)AddMember(m *storage.Member) {
	self.TSD.MemberList = append(self.TSD.MemberList, m)
}

type TopicAttribute struct {
	CreaterID          string
	CreaterSession     *libnet.Session
}

func NewTopicAttribute(CreaterID string, CreaterSession *libnet.Session) *TopicAttribute {
	return &TopicAttribute {
		CreaterID      : CreaterID,
		CreaterSession : CreaterSession,
	}
}