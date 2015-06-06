//
// Copyright 2015-2099 Hong Miao. All Rights Reserved.
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

package controllers

import (
)


type LoginPostData struct {
	Username               string `json:"username"`
	Password               string `json:"password"`
}

type LoginStatus struct {
	Status               string `json:"status"`
}

func NewLoginStatus() LoginStatus {
	return LoginStatus{}
}

type RebootStatus struct {
	Status               string `json:"status"`
}

func NewRebootStatus() RebootStatus {
	return RebootStatus{}
}

type TotalStatus struct {
	Ip               string `json:"ip"`
	Mac              string `json:"mac"`
	AllStorage       uint64 `json:"allStorage"`
	UsedStorage      uint64 `json:"usedStorage"`
	Modify           int64  `json:"modify"`
	Type             string `json:"type"`
}

func NewTotalStatus() TotalStatus {
	return TotalStatus{}
}

type MsgServerData struct {
	Status           string `json:"status"`
	Num              uint32 `json:"num"`            
}


func NewMsgServerData() MsgServerData {
	return MsgServerData{}
}

