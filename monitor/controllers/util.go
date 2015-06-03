package controllers

import (
	"net"
	"fmt"
	"bytes"
	"errors"
	"syscall"
	"strings"
	"os/exec"
	"github.com/astaxie/beego"
	"github.com/oikomi/FishChatServer/monitor/conf"
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
            //fmt.Printf("Interface %q, address %v\n", ifi.Name, a)
            //ipaddr , err := net.ResolveIPAddr("ip", a.String())
            //if err != nil {
            //    beego.Error("GetLocalMac Failed")
            //    return "", err
            //}
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


func GetDiskUsage() (*DiskUsageData, error) {
	dud := NewDiskUsageData()
	fs := syscall.Statfs_t{}
	var disk DiskStatus
	err := syscall.Statfs(conf.BaseDir, &fs)
	if err != nil {
		dud.Status = 1
		return nil, err
	}
	disk.All = fs.Blocks * uint64(fs.Bsize)
	disk.Free = fs.Bfree * uint64(fs.Bsize)
	disk.Used = disk.All - disk.Free
	dud.Status = 0
	dud.All = disk.All / (1024*1024)
	dud.Used = disk.Used / (1024*1024)
	return &dud, nil
}