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
	"flag"
	"strconv"
	"github.com/oikomi/FishChatServer/log"
	"github.com/oikomi/FishChatServer/libnet"
	"github.com/oikomi/FishChatServer/base"
	"github.com/oikomi/FishChatServer/protocol"
	"github.com/oikomi/FishChatServer/common"
	"github.com/oikomi/FishChatServer/storage/redis_store"
	"github.com/oikomi/FishChatServer/storage/mongo_store"
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

func (self *ProtoProc)procSubscribeChannel(cmd protocol.Cmd, session *libnet.Session) {
	log.Info("procSubscribeChannel")
	channelName := cmd.GetArgs()[0]
	cUUID := cmd.GetArgs()[1]
	log.Info(channelName)
	if self.msgServer.channels[channelName] != nil {
		self.msgServer.channels[channelName].Channel.Join(session, nil)
		self.msgServer.channels[channelName].ClientIDlist = append(self.msgServer.channels[channelName].ClientIDlist, cUUID)
	} else {
		log.Warning(channelName + " is not exist")
	}

	log.Info(self.msgServer.channels)
}

func (self *ProtoProc)procPing(cmd protocol.Cmd, session *libnet.Session) error {
	//log.Info("procPing")
	cid := session.State.(*base.SessionState).ClientID
	self.msgServer.scanSessionMutex.Lock()
	defer self.msgServer.scanSessionMutex.Unlock()
	self.msgServer.sessions[cid].State.(*base.SessionState).Alive = true
	
	return nil
}

func (self *ProtoProc)procOfflineMsg(session *libnet.Session, ID string) error {
	var err error
	exist, err := self.msgServer.offlineMsgCache.IsKeyExist(ID)
	if exist.(int64) == 0 {
		return err
	} else {
		omrd, err := common.GetOfflineMsgFromOwnerName(self.msgServer.offlineMsgCache, ID)
		if err != nil {
			log.Error(err.Error())
			return err
		}
		for  _, v := range omrd.MsgList {
			resp := protocol.NewCmdSimple(protocol.RESP_MESSAGE_P2P_CMD)
			resp.AddArg(v.Msg)
			resp.AddArg(v.FromID)
			resp.AddArg(v.Uuid)
			
			if self.msgServer.sessions[ID] != nil {
				self.msgServer.sessions[ID].Send(libnet.Json(resp))
				if err != nil {
					log.Error(err.Error())
					return err
				}
			} 
		}
		
		omrd.ClearMsg()
		self.msgServer.offlineMsgCache.Set(omrd)
	}
	
	return err
}

func (self *ProtoProc)procClientID(cmd protocol.Cmd, session *libnet.Session) error {
	log.Info("procClientID")
	var err error
	ID := cmd.GetArgs()[0]
	// for cache data
	sessionCacheData := redis_store.NewSessionCacheData(cmd.GetArgs()[0], session.Conn().RemoteAddr().String(), 
		self.msgServer.cfg.LocalIP, strconv.FormatUint(session.Id(), 10))
		
	log.Info(sessionCacheData)
	args := make([]string, 0)
	args = append(args, cmd.GetArgs()[0])
	CCmd := protocol.NewCmdInternal(protocol.CACHE_SESSION_CMD, args, sessionCacheData)
	
	log.Info(CCmd)
	
	if self.msgServer.channels[protocol.SYSCTRL_CLIENT_STATUS] != nil {
		_, err = self.msgServer.channels[protocol.SYSCTRL_CLIENT_STATUS].Channel.Broadcast(libnet.Json(CCmd))
		if err != nil {
			log.Error(err.Error())
			return err
		}
	}
	
	// for store data
	sessionStoreData := mongo_store.SessionStoreData{ID, session.Conn().RemoteAddr().String(), 
		self.msgServer.cfg.LocalIP, false}
		
	log.Info(sessionStoreData)
	args = make([]string, 0)
	args = append(args, cmd.GetArgs()[0])
	CCmd = protocol.NewCmdInternal(protocol.STORE_SESSION_CMD, args, sessionStoreData)
	
	log.Info(CCmd)
	
	if self.msgServer.channels[protocol.STORE_CLIENT_INFO] != nil {
		_, err = self.msgServer.channels[protocol.STORE_CLIENT_INFO].Channel.Broadcast(libnet.Json(CCmd))
		if err != nil {
			log.Error(err.Error())
			return err
		}
	}	

	self.msgServer.sessions[cmd.GetArgs()[0]] = session
	self.msgServer.sessions[cmd.GetArgs()[0]].State = base.NewSessionState(true, cmd.GetArgs()[0])
	
	err = self.procOfflineMsg(session, ID)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func (self *ProtoProc)procSendMessageP2P(cmd protocol.Cmd, session *libnet.Session) error {
	log.Info("procSendMessageP2P")
	var err error
	send2ID := cmd.GetArgs()[0]
	send2Msg := cmd.GetArgs()[1]
	fromID := cmd.GetArgs()[2]
	store_session, err := common.GetSessionFromCID(self.msgServer.sessionCache, send2ID)
	if err != nil {
		log.Warningf("no ID : %s", send2ID)
		
		return err
	}
	
	uuid := common.NewV4().String()
	
	log.Info("uuid : ", uuid)
	
	uuidTmpMap := make(map[string]bool)
	uuidTmpMap[uuid] = false
	
	self.msgServer.p2pAckStatus[fromID] = uuidTmpMap
	
	if self.msgServer.sessions[send2ID] == nil {
		//offline
		log.Info(send2ID + " | is offline")
		exist, err := self.msgServer.offlineMsgCache.IsKeyExist(send2ID)
		if exist.(int64) == 0 {
			tmp := redis_store.NewOfflineMsgCacheData(send2ID)
			tmp.AddMsg(redis_store.NewOfflineMsgData(send2Msg, fromID, uuid))
			
			self.msgServer.offlineMsgCache.Set(tmp)
			if err != nil {
				log.Error(err.Error())
				return err
			}
		} else {
			omrd, err := common.GetOfflineMsgFromOwnerName(self.msgServer.offlineMsgCache, send2ID)
			if err != nil {
				log.Error(err.Error())
				return err
			}
			omrd.AddMsg(redis_store.NewOfflineMsgData(send2Msg, fromID, uuid))
			self.msgServer.offlineMsgCache.Set(omrd)
			if err != nil {
				log.Error(err.Error())
				return err
			}
		}
	}
		
	if store_session.MsgServerAddr == self.msgServer.cfg.LocalIP {
		log.Info("in the same server")
		resp := protocol.NewCmdSimple(protocol.RESP_MESSAGE_P2P_CMD)
		resp.AddArg(send2Msg)
		resp.AddArg(fromID)
		// add uuid
		resp.AddArg(uuid)
		
		if self.msgServer.sessions[send2ID] != nil {
			self.msgServer.sessions[send2ID].Send(libnet.Json(resp))
			if err != nil {
				log.Error(err.Error())
			}
		} 
	} else {
		if self.msgServer.channels[protocol.SYSCTRL_SEND] != nil {
			//add uuid
			cmd.AddArg(uuid)
			_, err = self.msgServer.channels[protocol.SYSCTRL_SEND].Channel.Broadcast(libnet.Json(cmd))
			if err != nil {
				log.Error(err.Error())
				return err
			}
		}
	}
	
	return nil
}

func (self *ProtoProc)procRouteMessageP2P(cmd protocol.Cmd, session *libnet.Session) error {
	log.Info("procRouteMessageP2P")
	var err error
	send2ID := cmd.GetArgs()[0]
	send2Msg := cmd.GetArgs()[1]
	fromID := cmd.GetArgs()[2]
	uuid := cmd.GetArgs()[3]
	_, err = common.GetSessionFromCID(self.msgServer.sessionCache, send2ID)
	if err != nil {
		log.Warningf("no ID : %s", send2ID)
		
		return err
	}

	resp := protocol.NewCmdSimple(protocol.RESP_MESSAGE_P2P_CMD)
	resp.AddArg(send2Msg)
	resp.AddArg(fromID)
	// add uuid
	resp.AddArg(uuid)
	
	if self.msgServer.sessions[send2ID] != nil {
		self.msgServer.sessions[send2ID].Send(libnet.Json(resp))
		if err != nil {
			log.Fatalln(err.Error())
		}
	}

	return nil
}


func (self *ProtoProc)procCreateTopic(cmd protocol.Cmd, session *libnet.Session) error {
	log.Info("procCreateTopic")
	var err error
	
	if len(cmd.GetArgs()) != 1 {
		return CMD_NOT_CORRECT
	}
	topicName := cmd.GetArgs()[0]
	
	topicCacheData := redis_store.NewTopicCacheData(topicName, session.State.(*base.SessionState).ClientID, 
		self.msgServer.cfg.LocalIP)

	t := protocol.NewTopic(topicName, self.msgServer.cfg.LocalIP, session.State.(*base.SessionState).ClientID, session)
	t.ClientIDList = append(t.ClientIDList, session.State.(*base.SessionState).ClientID)
	t.TSD = topicCacheData
	self.msgServer.topics[topicName] = t
	self.msgServer.topics[topicName].Channel = libnet.NewChannel(self.msgServer.server.Protocol())
	
	self.msgServer.topics[topicName].Channel.Join(session, nil)
	

	log.Info(topicCacheData)
	args := make([]string, 0)
	args = append(args, topicName)
	CCmd := protocol.NewCmdInternal(protocol.CACHE_TOPIC_CMD, args, topicCacheData)
	m := redis_store.NewMember(session.State.(*base.SessionState).ClientID)
	CCmd.AnyData.(*redis_store.TopicCacheData).MemberList = append(CCmd.AnyData.(*redis_store.TopicCacheData).MemberList, m)
	
	log.Info(CCmd)
	
	if self.msgServer.channels[protocol.SYSCTRL_TOPIC_STATUS] != nil {
		_, err = self.msgServer.channels[protocol.SYSCTRL_TOPIC_STATUS].Channel.Broadcast(libnet.Json(CCmd))
		if err != nil {
			log.Error(err.Error())
			return err
		}
	}
	
	return nil
}

func (self *ProtoProc)findTopicMsgAddr(topicName string) (*redis_store.TopicCacheData, error) {
	log.Info("findTopicMsgAddr")
	t, err := common.GetTopicFromTopicName(self.msgServer.topicCache, topicName)
	
	return t, err
}

func (self *ProtoProc)procJoinTopic(cmd protocol.Cmd, session *libnet.Session) error {
	log.Info("procJoinTopic")
	var err error
	if len(cmd.GetArgs()) != 2 {
		return CMD_NOT_CORRECT
	}
	topicName := cmd.GetArgs()[0]
	clientID := cmd.GetArgs()[1]
	
	if self.msgServer.topics[topicName] == nil {
		log.Warning("no topic in this server :" + topicName)
		t, err := self.findTopicMsgAddr(topicName)
		if err != nil {
			log.Warningf("no topicName : %s", topicName)
			return err
		}
		
		resp := protocol.NewCmdSimple(protocol.LOCATE_TOPIC_MSG_ADDR_CMD)
		resp.AddArg(t.MsgServerAddr)
		resp.AddArg(topicName)
		
		err = session.Send(libnet.Json(resp))
		
		if err != nil {
			log.Error(err.Error())
			return err
		}
		
		return err
	} else {
		t, err := self.findTopicMsgAddr(topicName)
		if err != nil {
			log.Warningf("no topicName : %s", topicName)
			return err
		}
		resp := protocol.NewCmdSimple(protocol.LOCATE_TOPIC_MSG_ADDR_CMD)
		resp.AddArg(t.MsgServerAddr)
		resp.AddArg(topicName)
		
		err = session.Send(libnet.Json(resp))
		
		if err != nil {
			log.Error(err.Error())
			return err
		}
	}
	
	m := redis_store.NewMember(clientID)

	self.msgServer.topics[topicName].ClientIDList = append(self.msgServer.topics[topicName].ClientIDList, 
		clientID)
	
	self.msgServer.topics[topicName].AddMember(m)
	
	self.msgServer.topics[topicName].Channel.Join(session, nil)
	
	args := make([]string, 0)
	args = append(args, topicName)
	CCmd := protocol.NewCmdInternal(protocol.CACHE_TOPIC_CMD, args, self.msgServer.topics[topicName].TSD)
	
	log.Info(CCmd)
	
	if self.msgServer.channels[protocol.SYSCTRL_TOPIC_STATUS] != nil {
		_, err = self.msgServer.channels[protocol.SYSCTRL_TOPIC_STATUS].Channel.Broadcast(libnet.Json(CCmd))
		if err != nil {
			log.Error(err.Error())
			return err
		}
	}
	
	return nil
}

func (self *ProtoProc)procSendMessageTopic(cmd protocol.Cmd, session *libnet.Session) error {
	log.Info("procSendMessageTopic")
	var err error
	topicName := cmd.GetArgs()[0]
	send2Msg := cmd.GetArgs()[1]
	log.Info(send2Msg)
	log.Info(topicName)

	if self.msgServer.topics[topicName] == nil {
		log.Warning(topicName + " is not exist")
	} else {
		resp := protocol.NewCmdSimple(protocol.RESP_MESSAGE_TOPIC_CMD)
		resp.AddArg(topicName)
		resp.AddArg(send2Msg)
		resp.AddArg(session.State.(*base.SessionState).ClientID)

		_, err = self.msgServer.topics[topicName].Channel.Broadcast(libnet.Json(resp))
		if err != nil {
			log.Error(err.Error())
			return err
		}
	}
	
	return err
}


// not a good idea
func (self *ProtoProc)procP2pAck(cmd protocol.Cmd, session *libnet.Session) error {
	log.Info("procP2pAck")
	var err error
	clientID := cmd.GetArgs()[0]
	uuid := cmd.GetArgs()[1]
	self.msgServer.p2pAckMutex.Lock()
	defer self.msgServer.p2pAckMutex.Unlock()
	
	//self.msgServer.p2pAckStatus[clientID][uuid] = true
	
	m, ok := self.msgServer.p2pAckStatus[clientID]
	if ok {
		m[uuid] = true
	}
	
	return err
}