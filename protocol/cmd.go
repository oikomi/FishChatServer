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

const (
	SELECT_MSG_SERVER_FOR_CLIENT_CMD      = "SELECT_MSG_SERVER_FOR_CLIENT"
)

const (
	//SEND_PING
	SEND_PING_CMD                         = "SEND_PING"
	//SEND_CLIENT_ID CLIENT_ID
	SEND_CLIENT_ID_CMD                    = "SEND_CLIENT_ID"
	//SEND_CLIENT_ID_FOR_TOPIC ID
	SEND_CLIENT_ID_FOR_TOPIC_CMD          = "SEND_CLIENT_ID_FOR_TOPIC"
	//SUBSCRIBE_CHANNEL channelName
	SUBSCRIBE_CHANNEL_CMD                 = "SUBSCRIBE_CHANNEL"
	//SEND_MESSAGE_P2P send2ID send2msg
	SEND_MESSAGE_P2P_CMD                  = "SEND_MESSAGE_P2P"
	//RESP_MESSAGE_P2P  msg fromID
	RESP_MESSAGE_P2P_CMD                  = "RESP_MESSAGE_P2P"
	ROUTE_MESSAGE_P2P_CMD                 = "ROUTE_MESSAGE_P2P"
	CREATE_TOPIC_CMD                      = "CREATE_TOPIC"
	//JOIN_TOPIC TOPIC_NAME CLIENT_ID
	JOIN_TOPIC_CMD                        = "JOIN_TOPIC"
	LOCATE_TOPIC_MSG_ADDR_CMD             = "LOCATE_TOPIC_MSG_ADDR"
	SEND_MESSAGE_TOPIC_CMD                = "SEND_MESSAGE_TOPIC"
	RESP_MESSAGE_TOPIC_CMD                = "RESP_MESSAGE_TOPIC"
)

const (
	STORE_SESSION_CMD       = "STORE_SESSION"
	STORE_TOPIC_CMD         = "STORE_TOPIC"
)

const (
	PING  = "PING"
)

type Cmd interface {
	GetCmdName() string
	ChangeCmdName(newName string)
	GetArgs() []string
	AddArg(arg string)
	ParseCmd(msglist []string)
	GetAnyData() interface{}
}


type CmdSimple struct {
	CmdName string
	Args    []string
}

func NewCmdSimple() *CmdSimple {
	return &CmdSimple {
		CmdName : "",
		Args    : make([]string, 0),
	}
}

func (self CmdSimple)GetCmdName() string {
	return self.CmdName
}

func (self CmdSimple)ChangeCmdName(newName string) {
	self.CmdName = newName
}

func (self CmdSimple)GetArgs() []string {
	return self.Args
}

func (self CmdSimple)AddArg(arg string) {
	self.Args = append(self.Args, arg)
}

func (self CmdSimple)ParseCmd(msglist []string) {
	self.CmdName = msglist[1]
	self.Args = msglist[2:]
}

func (self CmdSimple)GetAnyData() interface{} {
	return nil
}

type CmdInternal struct {
	CmdName string
	Args    []string
	AnyData interface{}
}

func NewCmdInternal(cmdName string, args []string, anyData interface{}) *CmdInternal {
	return &CmdInternal {
		CmdName : cmdName,
		Args    : args,
		AnyData : anyData,
	}
}

func (self CmdInternal)ParseCmd(msglist []string) {
	self.CmdName = msglist[1]
	self.Args = msglist[2:]
}

func (self CmdInternal)GetCmdName() string {
	return self.CmdName
}

func (self CmdInternal)ChangeCmdName(newName string) {
	self.CmdName = newName
}

func (self CmdInternal)GetArgs() []string {
	return self.Args
}

func (self CmdInternal)AddArg(arg string) {
	self.Args = append(self.Args, arg)
}

func (self CmdInternal)SetAnyData(a interface{}) {
	self.AnyData = a
}

func (self CmdInternal)GetAnyData() interface{} {
	return self.AnyData
}

type ClientIDCmd struct {
	CmdName  string
	ClientID string
}

type SendMessageP2PCmd struct {
	CmdName string
	ID      string
	Msg     string
}