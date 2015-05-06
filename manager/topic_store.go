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

package main

import (
	"github.com/oikomi/FishChatServer/storage"
)

type TopicStoreCmd struct {
	CmdName string
	Args    []string
	AnyData *storage.TopicStoreData
}

func (self TopicStoreCmd)GetCmdName() string {
	return self.CmdName
}

func (self TopicStoreCmd)ChangeCmdName(newName string) {
	self.CmdName = newName
}

func (self TopicStoreCmd)GetArgs() []string {
	return self.Args
}

func (self TopicStoreCmd)AddArg(arg string) {
	self.Args = append(self.Args, arg)
}

func (self TopicStoreCmd)ParseCmd(msglist []string) {
	self.CmdName = msglist[1]
	self.Args = msglist[2:]
}

func (self TopicStoreCmd)GetAnyData() interface{} {
	return self.AnyData
}

