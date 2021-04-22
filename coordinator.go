package main

import (
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/sirupsen/logrus"
	"math/rand"
	"sync"
)

type Bench struct {
	aws *AWSWrapper
	c *Config
	l *logrus.Entry
	id string
	orignalLog *logrus.Entry

	// step1: make keyfiles
	names     map[string]string
	keys      map[string]*ec2.CreateKeyPairOutput
	instances map[string]string //ips
	sgs       map[string]*string // security group id
	nodes     map[string]*IPFS
	webuiport int
	counts map[string]int
	countsLock *sync.Mutex
}


func NewBench(l *logrus.Entry, c *Config, aws *AWSWrapper) *Bench {
	b := Bench{}
	b.aws = aws
	b.c = c
	b.id = RandStringRunes(12)
	b.l = l.WithField("source","bench")
	b.keys = map[string]*ec2.CreateKeyPairOutput{}
	b.instances = map[string]string{}
	b.sgs = map[string]*string{}
	b.webuiport = 55001
	b.nodes = map[string]*IPFS{}
	b.orignalLog = l
	b.counts = map[string]int{}
	b.countsLock = &sync.Mutex{}
	return &b
}




func RandStringRunes(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}