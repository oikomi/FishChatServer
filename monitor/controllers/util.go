//
// Copyright 2014-2015 Hong Miao (miaohong@miaohong.org). All Rights Reserved.
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
	"net"
	"fmt"
	"bytes"
	"errors"
	"strings"
	"os/exec"
	"github.com/astaxie/beego"
)


func RunShellCmd(s string) error {
	cmd := exec.Command("/bin/sh", "-c", s)
	var out bytes.Buffer

	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		beego.Error(err)
		return err
	}
	fmt.Printf("%s", out.String())
	
	return nil
}


func GetLocalIP(inter string) (string,error){
    ifi, err := net.InterfaceByName(inter)
    if err != nil {
        beego.Error("GetLocalIP Failed")
        return "", err
    }
    addrs, err := ifi.Addrs()
    if err != nil {
        beego.Error("GetLocalIP Failed")
        return "", err
    }
    for _, a := range addrs {
            fmt.Printf("Interface %q, address %v\n", ifi.Name, a)
            return a.String(), err
    }

    return "", err
}

func GetLocalMac(inter string) (string, error) {
    ifi, err := net.InterfaceByName(inter)
    if err != nil {
        beego.Error("GetLocalMac Failed")
        return "", err
    }
    return ifi.HardwareAddr.String(), nil
}

func GetLocalMask(inter string) (string, error) {
    ifi, err := net.InterfaceByName(inter)
    if err != nil {
        beego.Error("GetLocalMac Failed")
        return "", err
    }
    addrs, err := ifi.Addrs()
    if err != nil {
        beego.Error("GetLocalMac Failed")
        return "", err
    }
    for _, a := range addrs {
            fullIp := strings.Split(a.String(), "/")[0]
            ip := net.ParseIP(fullIp)

            if ip == nil {
                beego.Error("ParseIP Failed")
                return "", errors.New("ParseIP Failed")
            }

            return ip.DefaultMask().String(), err
    }

    return "", err
}
