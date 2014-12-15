FishChatServer
======

FishChat(鱼传 - `鱼传尺素`) - FishChat是一款优秀的即时通讯软件(IM), 它集合了市面上已有产品的有点, 并具备智能网关管理.
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

协议
---------------------
### 1. client to MsgServer

**syntax:** *SEND_PING_CMD*

**format:** *SEND_CLIENT_ID*

--

**syntax:** *SEND_CLIENT_ID_CMD*

**format:** *SEND_CLIENT_ID ID*

--

