package main

import (
	"context"
	"fmt"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/sirupsen/logrus"
)

type IPFS struct {
	addr string
	sh *shell.Shell
	l *logrus.Entry
}

func NewIpfs(addr,name string,l *logrus.Entry) *IPFS {
	n := IPFS{}
	n.addr = addr
	sh := shell.NewShell(addr)
	n.sh = sh
	n.l = l.WithField("source","ipfs").WithField("name",name)
	return &n
}

func (i *IPFS) Init(){
	// some basic settings
	e := i.sh.Request("/api/v0/config","API.HTTPHeaders.Access-Control-Allow-Origin","[*]").Exec(context.Background(),nil)
	if e != nil {
		fmt.Println(e)
	}
	e = i.sh.Request("/api/v0/config","API.HTTPHeaders.Access-Control-Allow-Methods","[*]").Exec(context.Background(),nil)
	if e != nil {
		fmt.Println(e)
	}
	go i.monitoring()
}

func (i *IPFS) monitoring(){
	// some basic settings
	e := i.sh.Request("/api/v0/config","API.HTTPHeaders.Access-Control-Allow-Origin","[*]").Exec(context.Background(),nil)
	if e != nil {
		fmt.Println(e)
	}
	e = i.sh.Request("/api/v0/config","API.HTTPHeaders.Access-Control-Allow-Methods","[*]").Exec(context.Background(),nil)
	if e != nil {
		fmt.Println(e)
	}
}

func (i *IPFS) GetShell() *shell.Shell {
	return i.sh
}