![logo](./docs/pics/logo.png)

FishChatServer
======

FishChat(鱼传 - `鱼传尺素`) - FishChat是一款优秀的即时通讯软件(IM), 它集合了市面上已有产品的优点, 并具备**智能硬件网关管理**.
> **NOTE:** 鱼传取自古乐府《饮马长城窟行》`客从远方来, 遗我双鲤鱼. 呼儿烹鲤鱼, 中有尺素书.` 

FishChat Android客户端见 https://github.com/oikomi/FishChatAndroid

Status
======

It is not usable yet and still under early development.

Todo
======
- [x] support p2p communication
- [x] support group communication
- [ ] support message encryption
- [ ] support file transfer
- [ ] support audio
- [ ] support video

系统架构
======

![](./docs/pics/fishchat.png)


技术细节
======
FishChatServer完全采用Golang开发(https://golang.org/)

音视频方案
---------------------
音视频采用nginx-rtmp架构,借助ffmpeg技术,客户端输出rtmp流,服务器输出hls流

协议
---------------------
### 1. client to MsgServer

**syntax:** *SEND_PING_CMD*

**format:** *SEND_CLIENT_ID*

--

**syntax:** *SEND_CLIENT_ID_CMD*

**format:** *SEND_CLIENT_ID ID*

--


Copyright & License
===================

Copyright 2014 Hong Miao. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

