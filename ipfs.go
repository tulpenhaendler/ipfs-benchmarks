package main

import (
	"bytes"
	"context"
	"fmt"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

type IPFS struct {
	addr string
	ma string
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
	i.l.Info("IPFS Basic Setup done - can open webui if u want")
	go i.monitoring()
}

func (i *IPFS) monitoring(){
}

func (i *IPFS) GetShell() *shell.Shell {
	return i.sh
}

func (i *IPFS) Addr() string {
	if i.ma != "" {
		return i.ma
	}
	id := ID{}
	e := i.sh.Request("id",).Exec(context.Background(),&id)
	if e != nil {
		fmt.Println(e)
	}
	i.ma = id.Addresses[len(id.Addresses)-1]
	return i.ma
}



func (i *IPFS) CanGetCid(hash string) bool {
	_,e := i.sh.Request("get",hash).Send(context.Background())
	if e != nil {
		return false
	}
	return true
}

// size in MB
func (i *IPFS) MakeRandomObject(size int) string {
	rand.Seed(time.Now().UnixNano())
	data := make([]byte, 1024*size)
	rand.Read(data)
	cid,e := i.sh.Add(bytes.NewBuffer(data))
	if e != nil {
		fmt.Println(e)
	}
	i.sh.Pin(cid)
	return cid
}




// models

type ID struct {
	ID              string   `json:"ID"`
	PublicKey       string   `json:"PublicKey"`
	Addresses       []string `json:"Addresses"`
	AgentVersion    string   `json:"AgentVersion"`
	ProtocolVersion string   `json:"ProtocolVersion"`
	Protocols       []string `json:"Protocols"`
}
