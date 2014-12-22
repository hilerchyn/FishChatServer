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
	"fmt"
	"flag"
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer/libnet"
	"github.com/oikomi/FishChatServer/common"
	"github.com/oikomi/FishChatServer/protocol"
)

/*
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
const char* build_time(void) {
	static const char* psz_build_time = "["__DATE__ " " __TIME__ "]";
	return psz_build_time;
}
*/
import "C"

var (
	buildTime = C.GoString(C.build_time())
)

func BuildTime() string {
	return buildTime
}

const VERSION string = "0.10"

func version() {
	fmt.Printf("gateway version %s Copyright (c) 2014 Harold Miao (miaohonghit@gmail.com)  \n", VERSION)
}

func init() {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", "false")
}

var InputConfFile = flag.String("conf_file", "gateway.json", "input conf file name")   

func main() {
	version()
	fmt.Printf("built on %s\n", BuildTime())
	flag.Parse()
	cfg := NewGatewayConfig(*InputConfFile)
	err := cfg.LoadConfig()
	if err != nil {
		glog.Error(err.Error())
		return
	}
	
	p := libnet.PacketN(2, libnet.BigEndian)
	
	server, err := libnet.Listen(cfg.TransportProtocols, cfg.Listen, p)
	if err != nil {
		glog.Error(err.Error())
		return
	}
	glog.Info("gateway server start at ", server.Listener().Addr().String())

	server.Handle(func(session *libnet.Session) {
		glog.Info("client ", session.Conn().RemoteAddr().String(), " | in")
		msgServer := common.SelectServer(cfg.MsgServerList, cfg.MsgServerNum)

		resp := protocol.NewCmdSimple(protocol.SELECT_MSG_SERVER_FOR_CLIENT_CMD)
		resp.AddArg(msgServer)
		
		if session != nil {
			session.Send(libnet.JSON {
				resp,
			})
			if err != nil {
				glog.Error(err.Error())
			}
		}
		session.Close(nil)
		glog.Info("client ", session.Conn().RemoteAddr().String(), " | close")
	})
}
