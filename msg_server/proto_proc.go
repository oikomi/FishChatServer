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
	"flag"
	"strconv"
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer/libnet"
	"github.com/oikomi/FishChatServer/base"
	"github.com/oikomi/FishChatServer/protocol"
	"github.com/oikomi/FishChatServer/common"
	"github.com/oikomi/FishChatServer/storage"
)

func init() {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", "false")
}

type ProtoProc struct {
	msgServer    *MsgServer
}

func NewProtoProc(msgServer *MsgServer) *ProtoProc {
	return &ProtoProc {
		msgServer : msgServer,
	}
}

func (self *ProtoProc)procPing(cmd protocol.Cmd, session *libnet.Session) error {
	glog.Info("procPing")
	cid := session.State.(*base.SessionState).ClientID
	self.msgServer.scanSessionMutex.Lock()
	defer self.msgServer.scanSessionMutex.Unlock()
	self.msgServer.sessions[cid].State.(*base.SessionState).Alive = true
	
	return nil
}

func (self *ProtoProc)procClientID(cmd protocol.Cmd, session *libnet.Session) error {
	glog.Info("procClientID")
	var err error
	sessionStoreData := storage.NewSessionStoreData(cmd.GetArgs()[0], session.Conn().RemoteAddr().String(), 
		self.msgServer.cfg.LocalIP, strconv.FormatUint(session.Id(), 10))
		
	glog.Info(sessionStoreData)
	args := make([]string, 0)
	args = append(args, cmd.GetArgs()[0])
	CCmd := protocol.NewCmdInternal(protocol.STORE_SESSION_CMD, args, sessionStoreData)
	
	glog.Info(CCmd)
	
	if self.msgServer.channels[protocol.SYSCTRL_CLIENT_STATUS] != nil {
		err = self.msgServer.channels[protocol.SYSCTRL_CLIENT_STATUS].Channel.Broadcast(libnet.JSON {
			CCmd,
		})
		if err != nil {
			glog.Error(err.Error())
			return err
		}
	}

	self.msgServer.sessions[cmd.GetArgs()[0]] = session
	self.msgServer.sessions[cmd.GetArgs()[0]].State = base.NewSessionState(true, cmd.GetArgs()[0])
	
	return nil
}

func (self *ProtoProc)procSendMessageP2P(cmd protocol.Cmd, session *libnet.Session) error {
	glog.Info("procSendMessageP2P")
	var err error
	send2ID := cmd.GetArgs()[0]
	send2Msg := cmd.GetArgs()[1]
	fromID := cmd.GetArgs()[2]
	store_session, err := common.GetSessionFromCID(self.msgServer.sessionStore, send2ID)
	if err != nil {
		glog.Warningf("no ID : %s", send2ID)
		
		return err
	}
	
	if self.msgServer.sessions[send2ID] == nil {
		//offline
		glog.Info(send2ID + " | is offline")
		exist, err := self.msgServer.offlineMsgStore.IsKeyExist(send2ID)
		if exist.(int64) == 0 {
			tmp := storage.NewOfflineMsgStoreData(send2ID)
			tmp.AddMsg(storage.NewOfflineMsgData(send2Msg, fromID))
			
			self.msgServer.offlineMsgStore.Set(tmp)
			if err != nil {
				glog.Error(err.Error())
				return err
			}
		} else {
			omrd, err := common.GetOfflineMsgFromOwnerName(self.msgServer.offlineMsgStore, send2ID)
			if err != nil {
				glog.Error(err.Error())
				return err
			}
			omrd.AddMsg(storage.NewOfflineMsgData(send2Msg, fromID))
			self.msgServer.offlineMsgStore.Set(omrd)
			if err != nil {
				glog.Error(err.Error())
				return err
			}
		}
	}
	
	if store_session.MsgServerAddr == self.msgServer.cfg.LocalIP {
		glog.Info("in the same server")
		resp := protocol.NewCmdSimple()
		resp.CmdName = protocol.RESP_MESSAGE_P2P_CMD
		resp.Args = append(resp.Args, send2Msg)
		resp.Args = append(resp.Args, fromID)
		
		if self.msgServer.sessions[send2ID] != nil {
			self.msgServer.sessions[send2ID].Send(libnet.JSON {
				resp,
			})
			if err != nil {
				glog.Error(err.Error())
			}
		} 
	} else {
		if self.msgServer.channels[protocol.SYSCTRL_SEND] != nil {
			err = self.msgServer.channels[protocol.SYSCTRL_SEND].Channel.Broadcast(libnet.JSON {
				cmd,
			})
			if err != nil {
				glog.Error(err.Error())
				return err
			}
		}
	}
	
	return nil
}

func (self *ProtoProc)procRouteMessageP2P(cmd protocol.Cmd, session *libnet.Session) error {
	glog.Info("procRouteMessageP2P")
	var err error
	send2ID := cmd.GetArgs()[0]
	send2Msg := cmd.GetArgs()[1]
	_, err = common.GetSessionFromCID(self.msgServer.sessionStore, send2ID)
	if err != nil {
		glog.Warningf("no ID : %s", send2ID)
		
		return err
	}

	resp := protocol.NewCmdSimple()
	resp.CmdName = protocol.RESP_MESSAGE_P2P_CMD
	resp.Args = append(resp.Args, send2Msg)
	
	if self.msgServer.sessions[send2ID] != nil {
		self.msgServer.sessions[send2ID].Send(libnet.JSON {
			resp,
		})
		if err != nil {
			glog.Fatalln(err.Error())
		}
	}

	return nil
}


func (self *ProtoProc)procSendMessageTopic(cmd protocol.Cmd, session *libnet.Session) error {
	glog.Info("procSendMessageTopic")
	var err error
	topicName := cmd.GetArgs()[0]
	send2Msg := cmd.GetArgs()[1]
	glog.Info(send2Msg)
	glog.Info(topicName)

	if self.msgServer.channels[protocol.SYSCTRL_TOPIC_SYNC] != nil {
		err = self.msgServer.channels[protocol.SYSCTRL_TOPIC_SYNC].Channel.Broadcast(libnet.JSON {
			cmd,
		})
		if err != nil {
			glog.Error(err.Error())
			return err
		}
	}
	
	return nil
}

func (self *ProtoProc)procSubscribeChannel(cmd protocol.Cmd, session *libnet.Session) {
	glog.Info("procSubscribeChannel")
	channelName := cmd.GetArgs()[0]
	cUUID := cmd.GetArgs()[1]
	glog.Info(channelName)
	if self.msgServer.channels[channelName] != nil {
		self.msgServer.channels[channelName].Channel.Join(session, nil)
		self.msgServer.channels[channelName].ClientIDlist = append(self.msgServer.channels[channelName].ClientIDlist, cUUID)
	} else {
		glog.Warning(channelName + " is not exist")
	}
}

func (self *ProtoProc)procCreateTopic(cmd protocol.Cmd, session *libnet.Session) error {
	glog.Info("procCreateTopic")
	var err error
	topicName := cmd.GetArgs()[0]
	
	topicStoreData := storage.NewTopicStoreData(topicName, session.State.(*base.SessionState).ClientID, 
		self.msgServer.cfg.LocalIP)

	t := protocol.NewTopic(topicName, self.msgServer.cfg.LocalIP, session.State.(*base.SessionState).ClientID, session)
	t.ClientIDList = append(t.ClientIDList, session.State.(*base.SessionState).ClientID)
	t.TSD = topicStoreData
	self.msgServer.topics[topicName] = t
	

	glog.Info(topicStoreData)
	args := make([]string, 0)
	args = append(args, topicName)
	CCmd := protocol.NewCmdInternal(protocol.STORE_TOPIC_CMD, args, topicStoreData)
	m := storage.NewMember(session.State.(*base.SessionState).ClientID)
	CCmd.AnyData.(*storage.TopicStoreData).MemberList = append(CCmd.AnyData.(*storage.TopicStoreData).MemberList, m)
	
	glog.Info(CCmd)
	
	if self.msgServer.channels[protocol.SYSCTRL_TOPIC_STATUS] != nil {
		err = self.msgServer.channels[protocol.SYSCTRL_TOPIC_STATUS].Channel.Broadcast(libnet.JSON {
			CCmd,
		})
		if err != nil {
			glog.Error(err.Error())
			return err
		}
	}
	
	return nil
}

func (self *ProtoProc)findTopicMsgAddr(topicName string) (*storage.TopicStoreData, error) {
	glog.Info("findTopicMsgAddr")
	t, err := common.GetTopicFromTopicName(self.msgServer.topicStore, topicName)
	
	return t, err
}

func (self *ProtoProc)procJoinTopic(cmd protocol.Cmd, session *libnet.Session) error {
	glog.Info("procJoinTopic")
	var err error
	topicName := cmd.GetArgs()[0]
	
	if self.msgServer.topics[topicName] == nil {
		glog.Warning("no topic :" + topicName)
		t, err := self.findTopicMsgAddr(topicName)
		if err != nil {
			glog.Warningf("no topicName : %s", topicName)
			return err
		}
		
		resp := protocol.NewCmdSimple()
		resp.CmdName = protocol.LOCATE_TOPIC_MSG_ADDR_CMD
		resp.Args = append(resp.Args, t.MsgServerAddr)
		
		err = session.Send(libnet.JSON {
			resp,
		})
		
		if err != nil {
			glog.Error(err.Error())
			return err
		}
		
		return err
	}
	
	m := storage.NewMember(session.State.(*base.SessionState).ClientID)

	self.msgServer.topics[topicName].ClientIDList = append(self.msgServer.topics[topicName].ClientIDList, 
		session.State.(*base.SessionState).ClientID)
	
	self.msgServer.topics[topicName].AddMember(m)
	
	args := make([]string, 0)
	args = append(args, topicName)
	CCmd := protocol.NewCmdInternal(protocol.STORE_TOPIC_CMD, args, self.msgServer.topics[topicName].TSD)
	
	glog.Info(CCmd)
	
	if self.msgServer.channels[protocol.SYSCTRL_TOPIC_STATUS] != nil {
		err = self.msgServer.channels[protocol.SYSCTRL_TOPIC_STATUS].Channel.Broadcast(libnet.JSON {
			CCmd,
		})
		if err != nil {
			glog.Error(err.Error())
			return err
		}
	}
	
	return nil
}
